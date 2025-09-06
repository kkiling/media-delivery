import React from 'react';
import { Accordion, Button } from 'react-bootstrap';
import { Video, Mic2, Subtitles, Trash2 } from 'lucide-react';
import { ContentMatch, Track, TrackType } from '@/api/api';

interface ContentItemProps {
  content: ContentMatch;
  index: number;
  unallocated: Track[];
  onRemoveVideo: (index: number, video: Track) => void;
  onRemoveAudio: (contentIndex: number, audioIndex: number, audio: Track) => void;
  onRemoveSubtitle: (contentIndex: number, subtitleIndex: number, subtitle: Track) => void;
  onReplaceVideo: (index: number, newVideoPath: string) => void;
  onReplaceAudio: (contentIndex: number, audioIndex: number, newAudioPath: string) => void;
  onReplaceSubtitle: (contentIndex: number, subtitleIndex: number, newSubtitlePath: string) => void;
  onAddAudio: (contentIndex: number) => void;
  onAddSubtitle: (contentIndex: number) => void;
}

export const ContentItem: React.FC<ContentItemProps> = ({
  content,
  index,
  unallocated,
  onRemoveVideo,
  onRemoveAudio,
  onRemoveSubtitle,
  onReplaceVideo,
  onReplaceAudio,
  onReplaceSubtitle,
  onAddAudio,
  onAddSubtitle,
}) => {
  const hasAvailableAudio = unallocated.some(
    (t) =>
      t.type === TrackType.TRACK_TYPE_AUDIO &&
      !(content.audio_files || []).some((a) => a.name === t.name)
  );

  const hasAvailableSubtitles = unallocated.some(
    (t) =>
      t.type === TrackType.TRACK_TYPE_SUBTITLE &&
      !(content.subtitles || []).some((s) => s.name === t.name)
  );

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
            onChange={(e) => onReplaceVideo(index, e.target.value)}
          >
            <option value={content.video?.relative_path || ''}>
              {content.video?.relative_path || ''}
            </option>
            {/* Фильтруем только видео файлы из unallocated */}
            {unallocated
              .filter((t) => t.type === TrackType.TRACK_TYPE_VIDEO)
              .map((track, idx) => (
                <option key={idx} value={track.relative_path}>
                  {track.relative_path}
                </option>
              ))}
          </select>
          {/* Кнопка удаления видео */}
          <Button
            variant="outline-danger"
            size="sm"
            disabled={!content.video?.relative_path}
            onClick={() => content.video && onRemoveVideo(index, content.video)}
          >
            <Trash2 size={18} />
          </Button>
        </div>
      </div>

      <div className="accordion-custom">
        {/* Аудио секция */}
        {((content.audio_files?.length || 0) > 0 || hasAvailableAudio) && (
          <Accordion>
            <Accordion.Item eventKey="0" className="border-0 mb-2">
              <Accordion.Header>
                <div className="d-flex align-items-center gap-2">
                  <Mic2 size={18} className="text-primary" />
                  <span>Audio files ({content.audio_files?.length || 0})</span>
                </div>
              </Accordion.Header>
              <Accordion.Body className="pt-2 pb-1">
                {content.audio_files?.map((audio, idx) => (
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
                          onChange={(e) => onReplaceAudio(index, idx, e.target.value)}
                        >
                          <option value={audio.relative_path}>{audio.relative_path}</option>
                          {unallocated
                            .filter(
                              (t) => t.type === TrackType.TRACK_TYPE_AUDIO && t.name !== audio.name
                            )
                            .map((track, i) => (
                              <option key={i} value={track.relative_path}>
                                {track.relative_path}
                              </option>
                            ))}
                        </select>
                      </div>
                      <Button
                        variant="outline-danger"
                        size="sm"
                        onClick={() => onRemoveAudio(index, idx, audio)}
                      >
                        <Trash2 size={16} />
                      </Button>
                    </div>
                  </div>
                ))}
                {hasAvailableAudio && (
                  <Button
                    variant="outline-primary"
                    size="sm"
                    className={`ms-4 ${!content.audio_files?.length ? 'mt-2' : ''}`}
                    onClick={() => onAddAudio(index)}
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
                          onChange={(e) => onReplaceSubtitle(index, idx, e.target.value)}
                        >
                          <option value={subtitle.relative_path}>{subtitle.relative_path}</option>
                          {unallocated
                            .filter(
                              (t) =>
                                t.type === TrackType.TRACK_TYPE_SUBTITLE && t.name !== subtitle.name
                            )
                            .map((track, i) => (
                              <option key={i} value={track.relative_path}>
                                {track.relative_path}
                              </option>
                            ))}
                        </select>
                      </div>
                      <Button
                        variant="outline-danger"
                        size="sm"
                        onClick={() => onRemoveSubtitle(index, idx, subtitle)}
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
                    onClick={() => onAddSubtitle(index)}
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
