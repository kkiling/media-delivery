import React from 'react';
import { Accordion, Button } from 'react-bootstrap';
import { Video, Mic2, Subtitles, Trash2 } from 'lucide-react';
import { ContentMatch, Track, TrackType } from '@/api/api';

interface ContentItemProps {
  content: ContentMatch;
  contentIndex: number;
  unallocated: Track[];
  onRemoveContent: (contentIndex: number, track: Track) => void;
  onReplaceContent: (contentIndex: number, newTrack: Track, oldTrack?: Track) => void;
  onAddContent: (contentIndex: number, type: TrackType) => void;
}

export const ContentItem: React.FC<ContentItemProps> = ({
  content,
  contentIndex,
  unallocated,
  onRemoveContent,
  onReplaceContent,
  onAddContent,
}) => {
  const hasAvailableAudio = unallocated.some(
    (t) =>
      t.type === TrackType.TRACK_TYPE_AUDIO &&
      !(content.audio_tracks || []).some((a) => a.name === t.name)
  );

  const hasAvailableSubtitles = unallocated.some(
    (t) =>
      t.type === TrackType.TRACK_TYPE_SUBTITLE &&
      !(content.subtitles || []).some((s) => s.name === t.name)
  );

  const onSelectTrack = (e: React.ChangeEvent<HTMLSelectElement>, oldTrack?: Track) => {
    const selectedTrack = unallocated.find((t) => t.relative_path === e.target.value);
    if (selectedTrack) {
      onReplaceContent(contentIndex, selectedTrack, oldTrack);
    }
  };

  // Фильтруем только видео файлы из нераспределенных треков
  const availableVideos = unallocated.filter((t) => t.type === TrackType.TRACK_TYPE_VIDEO);

  // Фильтруем аудио треки, исключая те, что уже используются в content.audio_files,
  // кроме текущего выбранного трека (currentAudio)
  const availableTracks = (current: Track) => {
    // Определяем массив треков в зависимости от типа
    const tracks =
      current.type === TrackType.TRACK_TYPE_AUDIO ? content.audio_tracks : content.subtitles;

    // Собираем уникальные name всех треков, кроме текущего
    const usedNames = new Set(
      (tracks || []).filter((t) => t.name !== current?.name).map((t) => t.name)
    );

    // Фильтруем unallocated по типу и исключаем использованные имена
    return unallocated.filter((t) => t.type === current.type && !usedNames.has(t.name));
  };

  return (
    <div className="py-3 border-top">
      <div className="d-flex justify-content-between align-items-center mb-3">
        <h6 className="mb-0">{content.episode?.relative_path}</h6>
      </div>
      {/* Видео секция */}
      <div className="d-flex align-items-center gap-2 mb-3 ps-2">
        <Video size={18} className="text-primary" />
        <div className="d-flex align-items-center gap-2 flex-grow-1">
          {/* Выпадающий список для выбора видео */}
          <select
            className="form-select"
            value={content.video?.relative_path || ''}
            onChange={(e) => onSelectTrack(e, content.video)}
          >
            <option value={content.video?.relative_path || ''}>
              {content.video?.relative_path || ''}
            </option>
            {/* Фильтруем только видео файлы из unallocated */}
            {availableVideos.map((track, idx) => (
              <option key={idx} value={track.relative_path}>
                {track.relative_path}
              </option>
            ))}
          </select>
          {/* Кнопка удаления*/}
          <Button
            variant="outline-danger"
            size="sm"
            disabled={!content.video?.relative_path}
            onClick={() => content.video && onRemoveContent(contentIndex, content.video)}
          >
            <Trash2 size={18} />
          </Button>
        </div>
      </div>

      <div className="accordion-custom">
        {/* Аудио секция */}
        {((content.audio_tracks?.length || 0) > 0 || hasAvailableAudio) && (
          <Accordion>
            <Accordion.Item eventKey="0" className="border-0 mb-2">
              <Accordion.Header>
                <div className="d-flex align-items-center gap-2">
                  <Mic2 size={18} className="text-primary" />
                  <span>Audio files ({content.audio_tracks?.length || 0})</span>
                </div>
              </Accordion.Header>
              <Accordion.Body className="pt-2 pb-1">
                {content.audio_tracks?.map((audio, idx) => (
                  <div key={idx} className="mb-3 ps-4">
                    <div className="d-block d-md-none mb-2 text-truncate text-secondary">
                      {audio.name || 'Unnamed track'}
                    </div>
                    <div className="d-flex align-items-center gap-2">
                      <div className="track-name text-truncate d-none d-md-block">
                        {audio.name || 'Unnamed track'}
                      </div>
                      <div className="flex-grow-1">
                        <select
                          className="form-select"
                          value={audio.relative_path}
                          onChange={(e) => onSelectTrack(e, audio)}
                        >
                          <option value={audio.relative_path}>{audio.relative_path}</option>
                          {availableTracks(audio).map((track, idx) => (
                            <option key={idx} value={track.relative_path}>
                              {track.relative_path}
                            </option>
                          ))}
                        </select>
                      </div>
                      <Button
                        variant="outline-danger"
                        size="sm"
                        onClick={() => onRemoveContent(contentIndex, audio)}
                      >
                        <Trash2 size={18} />
                      </Button>
                    </div>
                  </div>
                ))}
                {hasAvailableAudio && (
                  <Button
                    variant="outline-primary"
                    size="sm"
                    className={`ms-4 ${!content.audio_tracks?.length ? 'mt-2' : ''}`}
                    onClick={() => onAddContent(contentIndex, TrackType.TRACK_TYPE_AUDIO)}
                  >
                    Add Audio Track
                  </Button>
                )}
              </Accordion.Body>
            </Accordion.Item>
          </Accordion>
        )}

        {/* Субтитры секция */}
        {((content.subtitles?.length || 0) > 0 || hasAvailableSubtitles) && (
          <Accordion>
            <Accordion.Item eventKey="0" className="border-0">
              <Accordion.Header>
                <div className="d-flex align-items-center gap-2">
                  <Subtitles size={18} className="text-primary" />
                  <span>Subtitles ({content.subtitles?.length || 0})</span>
                </div>
              </Accordion.Header>
              <Accordion.Body className="pt-2 pb-1">
                {content.subtitles?.map((subtitle, idx) => (
                  <div key={idx} className="mb-3 ps-4">
                    <div className="d-block d-md-none mb-2 text-truncate text-secondary">
                      {subtitle.name || 'Unnamed track'}
                    </div>
                    <div className="d-flex align-items-center gap-2">
                      <div className="track-name text-truncate d-none d-md-block">
                        {subtitle.name || 'Unnamed track'}
                      </div>
                      <div className="flex-grow-1">
                        <select
                          className="form-select"
                          value={subtitle.relative_path}
                          onChange={(e) => onSelectTrack(e, subtitle)}
                        >
                          <option value={subtitle.relative_path}>{subtitle.relative_path}</option>
                          {availableTracks(subtitle).map((track, idx) => (
                            <option key={idx} value={track.relative_path}>
                              {track.relative_path}
                            </option>
                          ))}
                        </select>
                      </div>
                      <Button
                        variant="outline-danger"
                        size="sm"
                        onClick={() => onRemoveContent(contentIndex, subtitle)}
                      >
                        <Trash2 size={18} />
                      </Button>
                    </div>
                  </div>
                ))}
                {hasAvailableSubtitles && (
                  <Button
                    variant="outline-primary"
                    size="sm"
                    className={`ms-4 ${!content.subtitles?.length ? 'mt-2' : ''}`}
                    onClick={() => onAddContent(contentIndex, TrackType.TRACK_TYPE_SUBTITLE)}
                  >
                    Add Subtitle Track
                  </Button>
                )}
              </Accordion.Body>
            </Accordion.Item>
          </Accordion>
        )}
      </div>
    </div>
  );
};
