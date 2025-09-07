import React, { useState, useEffect } from 'react';
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
  const [options, setOptions] = useState<Options>(contentMatches.options || {});
  const { matches, setMatches, updateMatch } = useMatchesState(contentMatches.matches || []);
  const { unallocated, addToUnallocated, removeFromUnallocated } = useTrackManagement(
    sortTracks(contentMatches.unallocated || [])
  );

  // Устанавливаем первый аудио трек если не выбран
  useEffect(() => {
    if (matches.length > 0 && !options.default_audio_track_name) {
      const firstAudioTrack = matches.flatMap((m) => m.audio_files || []).find((t) => t.name)?.name;
      if (firstAudioTrack) {
        setOptions((prev) => ({
          ...prev,
          default_audio_track_name: firstAudioTrack,
        }));
      }
    }
  }, [matches, options.default_audio_track_name]);

  // Получить уникальные имена аудио и субтитров
  const uniqueAudioNames = Array.from(
    new Set(
      matches
        .flatMap((m) => m.audio_files || [])
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
        audio_files: (match.audio_files || []).filter(
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
          audio_files: oldTrack
            ? (match.audio_files || []).map((track) =>
                track.relative_path === oldTrack.relative_path ? newTrack : track
              )
            : [...(match.audio_files || []), newTrack],
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
    const newTrack = unallocated.find(
      (t) => t.type === type && !matches[contentIndex].audio_files?.some((a) => a.name === t.name)
    );
    if (!newTrack) return;

    if (type === TrackType.TRACK_TYPE_AUDIO) {
      updateMatch(contentIndex, (match) => ({
        ...match,
        audio_files: [...(match.audio_files || []), newTrack],
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
  const removeTracksByName = (type: 'audio' | 'subtitle', name: string) => {
    const newMatches = matches.map((match) => {
      if (type === 'audio') {
        const removedTracks = (match.audio_files || []).filter((t) => t.name === name);
        removedTracks.forEach(addToUnallocated);
        return {
          ...match,
          audio_files: (match.audio_files || []).filter((t) => t.name !== name),
        };
      } else {
        const removedTracks = (match.subtitles || []).filter((t) => t.name === name);
        removedTracks.forEach(addToUnallocated);
        return {
          ...match,
          subtitles: (match.subtitles || []).filter((t) => t.name !== name),
        };
      }
    });
    setMatches(newMatches);
  };

  const handleConfirmClick = () => setShowConfirmModal(true);

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
          removeTracksByName={removeTracksByName}
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
