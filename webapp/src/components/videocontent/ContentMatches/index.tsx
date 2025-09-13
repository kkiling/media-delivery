import React, { useState } from 'react';
import { Button, Spinner } from 'react-bootstrap';
import { CheckCircle } from 'lucide-react';
import { ContentMatches as MatchContents, Options, Track, TrackType } from '@/api/api';
import { ConfirmModal } from './components/ConfirmModal';
import { ContentItem } from './components/ContentItem';
import { OptionsPanel } from './components/OptionsPanel';
import { UnallocatedFiles } from './components/UnallocatedFiles';
import { useMatchesState } from './hooks/useMatchesState';
import { useTrackManagement } from './hooks/useTrackManagement';
import { sortTracks } from './utils/sortTracks';

export interface ContentMatchesProps {
  loading: boolean;
  contentMatches: MatchContents;
  onConfirm: (contentMatches?: MatchContents) => void;
}

export const ContentMatches: React.FC<ContentMatchesProps> = ({
  contentMatches,
  onConfirm,
  loading,
}) => {
  const [showConfirmModal, setShowConfirmModal] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [options, setOptions] = useState<Options>(contentMatches.options || {});
  const { matches, setMatches, updateMatch } = useMatchesState(contentMatches.matches || []);
  const { unallocated, addToUnallocated, removeFromUnallocated } = useTrackManagement(
    sortTracks(contentMatches.unallocated || [])
  );

  // Получить уникальные имена аудио и субтитров
  const uniqueAudioNames = Array.from(
    new Set(
      matches
        .flatMap((m) => m.audio_tracks || [])
        .map((t) => t.name || '')
        .filter((t) => t !== '')
    )
  );
  const uniqueSubtitleNames = Array.from(
    new Set(
      matches
        .flatMap((m) => m.subtitles || [])
        .map((t) => t.name || '')
        .filter((t) => t !== '')
    )
  );

  // Обработчики для видео
  const handleRemoveContent = (index: number, track: Track) => {
    if (track.type === TrackType.TRACK_TYPE_VIDEO) {
      updateMatch(index, (match) => ({
        ...match,
        video: undefined,
      }));
    }

    if (track.type === TrackType.TRACK_TYPE_AUDIO) {
      updateMatch(index, (match) => ({
        ...match,
        audio_tracks: (match.audio_tracks || []).filter(
          (t) => t.relative_path !== track.relative_path
        ),
      }));
    }

    if (track.type === TrackType.TRACK_TYPE_SUBTITLE) {
      updateMatch(index, (match) => ({
        ...match,
        subtitles: (match.subtitles || []).filter((t) => t.relative_path !== track.relative_path),
      }));
    }

    addToUnallocated(track);
  };

  const handleReplaceContent = (contentIndex: number, newTrack: Track, oldTrack?: Track) => {
    updateMatch(contentIndex, (match) => {
      if (newTrack.type === TrackType.TRACK_TYPE_VIDEO) {
        const old = matches[contentIndex].video;
        if (old) {
          addToUnallocated(old);
        }
        return {
          ...match,
          video: newTrack,
        };
      } else if (newTrack.type === TrackType.TRACK_TYPE_AUDIO) {
        // Если есть oldTrack, заменяем его на newTrack
        // Если нет, добавляем newTrack в конец массива
        return {
          ...match,
          audio_tracks: oldTrack
            ? (match.audio_tracks || []).map((track) =>
                track.relative_path === oldTrack.relative_path ? newTrack : track
              )
            : [...(match.audio_tracks || []), newTrack],
        };
      } else if (newTrack.type === TrackType.TRACK_TYPE_SUBTITLE) {
        // Аналогичная логика для субтитров
        return {
          ...match,
          subtitles: oldTrack
            ? (match.subtitles || []).map((track) =>
                track.relative_path === oldTrack.relative_path ? newTrack : track
              )
            : [...(match.subtitles || []), newTrack],
        };
      }
      return match;
    });

    removeFromUnallocated(newTrack.relative_path);
  };

  const handleAddContent = (contentIndex: number, type: TrackType) => {
    // Определяем, какие треки уже назначены для данного эпизода по нужному типу
    const assignedTracks =
      type === TrackType.TRACK_TYPE_AUDIO
        ? matches[contentIndex].audio_tracks || []
        : matches[contentIndex].subtitles || [];

    // Ищем первый неиспользованный трек нужного типа
    const newTrack = unallocated.find(
      (t) => t.type === type && !assignedTracks.some((a) => a.name === t.name)
    );

    if (!newTrack) return;

    if (type === TrackType.TRACK_TYPE_AUDIO) {
      updateMatch(contentIndex, (match) => ({
        ...match,
        audio_tracks: [...(match.audio_tracks || []), newTrack],
      }));
    } else if (type === TrackType.TRACK_TYPE_SUBTITLE) {
      updateMatch(contentIndex, (match) => ({
        ...match,
        subtitles: [...(match.subtitles || []), newTrack],
      }));
    }

    removeFromUnallocated(newTrack.relative_path);
  };

  // Удаление треков по имени
  const removeTracks = (type: TrackType[], name?: string) => {
    // Собираем треки для перемещения в unallocated
    const check = (track: Track) =>
      track.type !== undefined && type.includes(track.type) && (!name || track.name === name);
    matches
      .flatMap((match) => [
        ...(match.audio_tracks || []).filter(check),
        ...(match.subtitles || []).filter(check),
        ...(match.video ? [match.video].filter(check) : []),
      ])
      .forEach((track) => addToUnallocated(track));

    // Обновляем состояние matches
    const filterTrack = (track: Track) => (name ? track.name !== name : false);
    setMatches(
      matches.map((match) => ({
        ...match,
        audio_tracks: type.includes(TrackType.TRACK_TYPE_AUDIO)
          ? match.audio_tracks?.filter(filterTrack)
          : match.audio_tracks,
        subtitles: type.includes(TrackType.TRACK_TYPE_SUBTITLE)
          ? match.subtitles?.filter(filterTrack)
          : match.subtitles,
        video:
          type.includes(TrackType.TRACK_TYPE_VIDEO) && (!name || match.video?.name === name)
            ? undefined
            : match.video,
      }))
    );
  };

  const handleConfirmClick = () => {
    setError(null);

    const invalidMatch = matches.find(
      (match) =>
        !match.video?.full_path &&
        ((match.audio_tracks && match.audio_tracks.length > 0) ||
          (match.subtitles && match.subtitles.length > 0))
    );

    if (invalidMatch) {
      setError(
        `Episode "${invalidMatch.episode?.relative_path}" has audio or subtitle tracks but no video assigned`
      );
      return;
    }

    setShowConfirmModal(true);
  };

  const handleConfirmClose = () => setShowConfirmModal(false);

  const handleConfirmSubmit = () => {
    onConfirm({
      matches,
      unallocated,
      options,
    });
    setShowConfirmModal(false);
  };

  return (
    <>
      <div className="container-fluid py-3">
        <OptionsPanel
          options={options}
          setOptions={setOptions}
          uniqueAudioNames={uniqueAudioNames}
          uniqueSubtitleNames={uniqueSubtitleNames}
          removeTracks={removeTracks}
        />

        <div className="border-top border-bottom">
          {matches.map((content, contentIndex) => (
            <ContentItem
              key={contentIndex}
              content={content}
              contentIndex={contentIndex}
              unallocated={unallocated}
              onRemoveContent={handleRemoveContent}
              onReplaceContent={handleReplaceContent}
              onAddContent={handleAddContent}
            />
          ))}
          <UnallocatedFiles unallocated={unallocated} />
        </div>

        {error && <div className="alert alert-danger mt-3">{error}</div>}

        <div className="mt-3 d-flex justify-content-end">
          <Button
            variant="primary"
            onClick={handleConfirmClick}
            disabled={!matches.length || loading}
            className="d-inline-flex align-items-center gap-2"
          >
            {loading ? <Spinner animation="border" size="sm" /> : <CheckCircle size={24} />}
            Confirm
          </Button>
        </div>
      </div>

      <ConfirmModal
        show={showConfirmModal}
        onClose={handleConfirmClose}
        onSubmit={handleConfirmSubmit}
      />

      <style>{`
        .accordion-custom .accordion-button {
          padding: 0.5rem 1rem;
          box-shadow: none;
          background-color: transparent;
        }
        .accordion-custom .accordion-button:not(.collapsed) {
          background-color: transparent;
          box-shadow: none;
        }
        .accordion-custom .accordion-button:focus {
          box-shadow: none;
          border-color: transparent;
        }
        .accordion-custom .accordion-body {
          padding-left: 1rem;
          padding-right: 1rem;
        }
        .track-name {
          width: 200px;
          min-width: 200px;
          color: var(--bs-secondary);
        }
        @media (max-width: 768px) {
          .track-name {
            width: 100%;
            min-width: 100%;
          }
        }
      `}</style>
    </>
  );
};
