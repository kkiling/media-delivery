import React, { useState, useEffect, useCallback } from 'react';
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
  const handleRemoveVideo = useCallback(
    (index: number, video: Track) => {
      console.log('Removing video at index:', index);
      updateMatch(index, (match) => ({
        ...match,
        video: undefined,
      }));
      addToUnallocated(video);
      console.log('Video removed');
    },
    [updateMatch, addToUnallocated]
  );

  const handleReplaceVideo = (index: number, newPath: string) => {
    const newContent = unallocated.find((t) => t.relative_path === newPath);
    if (!newContent) return;

    updateMatch(index, (match) => {
      const old = matches[index].video;
      if (old) {
        addToUnallocated(old);
      }
      return {
        ...match,
        video: newContent,
      };
    });
    removeFromUnallocated(newPath);
  };

  // Обработчики для аудио
  const handleRemoveAudio = (contentIndex: number, audioIndex: number, audio: Track) => {
    updateMatch(contentIndex, (match) => ({
      ...match,
      audio_files: (match.audio_files || []).filter((_, idx) => idx !== audioIndex),
    }));
    addToUnallocated(audio);
  };

  const handleReplaceAudio = (contentIndex: number, audioIndex: number, newAudioPath: string) => {
    const newAudio = unallocated.find((t) => t.relative_path === newAudioPath);
    if (!newAudio) return;

    updateMatch(contentIndex, (match) => {
      const audioFiles = [...(match.audio_files || [])];
      const oldAudio = audioFiles[audioIndex];
      if (oldAudio) {
        addToUnallocated(oldAudio);
      }
      audioFiles[audioIndex] = newAudio;
      return {
        ...match,
        audio_files: audioFiles,
      };
    });
    removeFromUnallocated(newAudioPath);
  };

  const handleAddAudio = (contentIndex: number) => {
    const newAudio = unallocated.find(
      (t) =>
        t.type === TrackType.TRACK_TYPE_AUDIO &&
        !matches[contentIndex].audio_files?.some((a) => a.name === t.name)
    );
    if (!newAudio) return;

    updateMatch(contentIndex, (match) => ({
      ...match,
      audio_files: [...(match.audio_files || []), newAudio],
    }));
    removeFromUnallocated(newAudio.relative_path);
  };

  // Обработчики для субтитров
  const handleRemoveSubtitle = (contentIndex: number, subtitleIndex: number, subtitle: Track) => {
    updateMatch(contentIndex, (match) => ({
      ...match,
      subtitles: (match.subtitles || []).filter((_, idx) => idx !== subtitleIndex),
    }));
    addToUnallocated(subtitle);
  };

  const handleReplaceSubtitle = (
    contentIndex: number,
    subtitleIndex: number,
    newSubtitlePath: string
  ) => {
    const newSubtitle = unallocated.find((t) => t.relative_path === newSubtitlePath);
    if (!newSubtitle) return;

    updateMatch(contentIndex, (match) => {
      const subtitles = [...(match.subtitles || [])];
      const oldSubtitle = subtitles[subtitleIndex];
      if (oldSubtitle) {
        addToUnallocated(oldSubtitle);
      }
      subtitles[subtitleIndex] = newSubtitle;
      return {
        ...match,
        subtitles: subtitles,
      };
    });
    removeFromUnallocated(newSubtitlePath);
  };

  const handleAddSubtitle = (contentIndex: number) => {
    const newSubtitle = unallocated.find(
      (t) =>
        t.type === TrackType.TRACK_TYPE_SUBTITLE &&
        !matches[contentIndex].subtitles?.some((s) => s.name === t.name)
    );
    if (!newSubtitle) return;

    updateMatch(contentIndex, (match) => ({
      ...match,
      subtitles: [...(match.subtitles || []), newSubtitle],
    }));
    removeFromUnallocated(newSubtitle.relative_path);
  };

  // Удаление треков по имени
  const removeTracksByName = useCallback(
    (type: 'audio' | 'subtitle', name: string) => {
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
    },
    [matches, setMatches, addToUnallocated]
  );

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
          {matches.map((content, index) => (
            <ContentItem
              key={index}
              content={content}
              index={index}
              unallocated={unallocated}
              onRemoveVideo={handleRemoveVideo}
              onRemoveAudio={handleRemoveAudio}
              onRemoveSubtitle={handleRemoveSubtitle}
              onReplaceVideo={handleReplaceVideo}
              onReplaceAudio={handleReplaceAudio}
              onReplaceSubtitle={handleReplaceSubtitle}
              onAddAudio={handleAddAudio}
              onAddSubtitle={handleAddSubtitle}
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
