import { useState, useEffect } from 'react';
import { Button, Accordion, Spinner } from 'react-bootstrap';
import { Video, Mic2, Subtitles, CheckCircle, Trash2 } from 'lucide-react';

export interface Episode {
  episode_number: number;
  season_number: number;
  episode_file: string;
}

export type TrackType = 'video' | 'audio' | 'subtitle';

export interface Track {
  name?: string;
  file: string;
  type: TrackType;
}

export interface ContentMatches {
  episode: Episode;
  video: Track;
  audio_tracks: Track[];
  subtitle_tracks: Track[];
}

export interface ContentMatchesProps {
  loading: boolean;
  unallocatedTracks: Track[];
  contentMatches?: ContentMatches[];
  onConfirm: (contentMatches: ContentMatches[], unallocated: Track[]) => void;
}

export const ContentMatche = ({
  contentMatches = [],
  unallocatedTracks,
  onConfirm,
  loading,
}: ContentMatchesProps) => {
  const [matches, setMatches] = useState<ContentMatches[]>([]);
  const [unallocated, setUnallocated] = useState<Track[]>([]);
  // Инициализация из props
  useEffect(() => {
    setMatches(contentMatches);
    setUnallocated(unallocatedTracks);
  }, [contentMatches, unallocatedTracks]);

  const handleRemoveVideo = (index: number, video: Track) => {
    // Удаляем видео из эпизода
    const newMatches = matches.map((m, idx) => {
      if (idx === index) {
        return {
          ...m,
          video: { file: '', type: 'video' as TrackType }, // пустое видео
        };
      }
      return m;
    });

    // Добавляем видео в unallocated
    setUnallocated([...unallocated, video]);
    setMatches(newMatches);
  };

  const handleRemoveAudio = (contentIndex: number, audioIndex: number, audio: Track) => {
    const newMatches = matches.map((m, idx) => {
      if (idx === contentIndex) {
        const newAudioTracks = [...m.audio_tracks];
        newAudioTracks.splice(audioIndex, 1);
        return {
          ...m,
          audio_tracks: newAudioTracks,
        };
      }
      return m;
    });

    setUnallocated([...unallocated, audio]);
    setMatches(newMatches);
  };

  const handleRemoveSubtitle = (contentIndex: number, subtitleIndex: number, subtitle: Track) => {
    const newMatches = matches.map((m, idx) => {
      if (idx === contentIndex) {
        const newSubtitles = [...m.subtitle_tracks];
        newSubtitles.splice(subtitleIndex, 1);
        return {
          ...m,
          subtitle_tracks: newSubtitles,
        };
      }
      return m;
    });

    setUnallocated([...unallocated, subtitle]);
    setMatches(newMatches);
  };

  return (
    <div className="container-fluid py-3">
      <div className="border-top border-bottom">
        {matches.map((content, index) => (
          <div key={index} className={`py-3 ${index !== 0 ? 'border-top' : ''}`}>
            <div className="d-flex justify-content-between align-items-center mb-3">
              <h6 className="mb-0">{content.episode.episode_file}</h6>
            </div>

            <div className="d-flex align-items-center gap-2 mb-3 ps-2">
              <Video size={18} className="text-primary" />
              <div className="d-flex align-items-center gap-2 flex-grow-1">
                <select
                  className="form-select"
                  value={content.video.file}
                  onChange={(e) => {
                    const newVideo = unallocated.find((t) => t.file === e.target.value);
                    if (!newVideo) return;

                    // Обновляем matches и unallocated за одну операцию
                    const newMatches = matches.map((m, idx) => {
                      if (idx === index) {
                        return {
                          ...m,
                          video: newVideo,
                        };
                      }
                      return m;
                    });

                    // Обновляем unallocated за одну операцию
                    const newUnallocated = unallocated
                      .filter((t) => t.file !== newVideo.file)
                      .concat(content.video.file ? content.video : []);

                    setUnallocated(newUnallocated);
                    setMatches(newMatches);
                  }}
                >
                  <option value={content.video.file}>{content.video.file}</option>
                  {unallocated
                    .filter((t) => t.type === 'video')
                    .map((track, idx) => (
                      <option key={idx} value={track.file}>
                        {track.file}
                      </option>
                    ))}
                </select>
                <Button
                  variant="outline-danger"
                  size="sm"
                  disabled={!content.video.file}
                  onClick={() => handleRemoveVideo(index, content.video)}
                >
                  <Trash2 size={16} />
                </Button>
              </div>
            </div>

            <div className="accordion-custom">
              {/* Секция аудио - показываем всегда если есть треки или доступны для добавления */}
              {(content.audio_tracks.length > 0 || unallocated.some((t) => t.type === 'audio')) && (
                <Accordion>
                  <Accordion.Item eventKey="0" className="border-0 mb-2">
                    <Accordion.Header>
                      <div className="d-flex align-items-center gap-2">
                        <Mic2 size={18} className="text-primary" />
                        <span>Audio files ({content.audio_tracks.length})</span>
                      </div>
                    </Accordion.Header>
                    <Accordion.Body className="pt-2 pb-1">
                      {content.audio_tracks.map((audio, idx) => (
                        <div key={idx} className="mb-3 ps-4">
                          {/* На маленьких экранах имя сверху */}
                          <div className="d-block d-md-none mb-2 text-truncate text-secondary">
                            {audio.name || 'Unnamed track'}
                          </div>
                          <div className="d-flex align-items-center gap-2">
                            {/* На средних и больших экранах имя слева */}
                            <div className="track-name text-truncate d-none d-md-block">
                              {audio.name || 'Unnamed track'}
                            </div>
                            <div className="flex-grow-1">
                              <select
                                className="form-select"
                                value={audio.file}
                                onChange={(e) => {
                                  const newAudio = unallocated.find(
                                    (t) => t.file === e.target.value
                                  );
                                  if (!newAudio) return;

                                  const newMatches = matches.map((m, mIdx) => {
                                    if (mIdx === index) {
                                      const newAudioTracks = [...m.audio_tracks];
                                      newAudioTracks[idx] = newAudio;
                                      return {
                                        ...m,
                                        audio_tracks: newAudioTracks,
                                      };
                                    }
                                    return m;
                                  });

                                  // Обновляем unallocated за одну операцию
                                  const newUnallocated = unallocated
                                    .filter((t) => t.file !== newAudio.file)
                                    .concat(audio);

                                  setUnallocated(newUnallocated);
                                  setMatches(newMatches);
                                }}
                              >
                                <option value={audio.file}>{audio.file}</option>
                                {unallocated
                                  .filter((t) => t.type === 'audio')
                                  .map((track, i) => (
                                    <option key={i} value={track.file}>
                                      {track.file}
                                    </option>
                                  ))}
                              </select>
                            </div>
                            <Button
                              variant="outline-danger"
                              size="sm"
                              onClick={() => handleRemoveAudio(index, idx, audio)}
                            >
                              <Trash2 size={16} />
                            </Button>
                          </div>
                        </div>
                      ))}
                      {/* Кнопка добавления новой аудио дорожки */}
                      {unallocated.some((t) => t.type === 'audio') && (
                        <Button
                          variant="outline-primary"
                          size="sm"
                          className={`ms-4 ${content.audio_tracks.length === 0 ? 'mt-2' : ''}`}
                          onClick={() => {
                            const newAudio = unallocated.find((t) => t.type === 'audio');
                            if (!newAudio) return;

                            const newMatches = matches.map((m, idx) => {
                              if (idx === index) {
                                return {
                                  ...m,
                                  audio_tracks: [...m.audio_tracks, newAudio],
                                };
                              }
                              return m;
                            });

                            setUnallocated(unallocated.filter((t) => t.file !== newAudio.file));
                            setMatches(newMatches);
                          }}
                        >
                          Add Audio Track
                        </Button>
                      )}
                    </Accordion.Body>
                  </Accordion.Item>
                </Accordion>
              )}

              {/* Секция субтитров - показываем всегда если есть треки или доступны для добавления */}
              {(content.subtitle_tracks.length > 0 ||
                unallocated.some((t) => t.type === 'subtitle')) && (
                <Accordion>
                  <Accordion.Item eventKey="0" className="border-0">
                    <Accordion.Header>
                      <div className="d-flex align-items-center gap-2">
                        <Subtitles size={18} className="text-primary" />
                        <span>Subtitles ({content.subtitle_tracks.length})</span>
                      </div>
                    </Accordion.Header>
                    <Accordion.Body className="pt-2 pb-1">
                      {content.subtitle_tracks.map((subtitle, idx) => (
                        <div key={idx} className="mb-3 ps-4">
                          {/* На маленьких экранах имя сверху */}
                          <div className="d-block d-md-none mb-2 text-truncate text-secondary">
                            {subtitle.name || 'Unnamed track'}
                          </div>
                          <div className="d-flex align-items-center gap-2">
                            {/* На средних и больших экранах имя слева */}
                            <div className="track-name text-truncate d-none d-md-block">
                              {subtitle.name || 'Unnamed track'}
                            </div>
                            <div className="flex-grow-1">
                              <select
                                className="form-select"
                                value={subtitle.file}
                                onChange={(e) => {
                                  const newSubtitle = unallocated.find(
                                    (t) => t.file === e.target.value
                                  );
                                  if (!newSubtitle) return;

                                  const newMatches = matches.map((m, mIdx) => {
                                    if (mIdx === index) {
                                      const newSubtitles = [...m.subtitle_tracks];
                                      newSubtitles[idx] = newSubtitle;
                                      return {
                                        ...m,
                                        subtitle_tracks: newSubtitles,
                                      };
                                    }
                                    return m;
                                  });

                                  // Обновляем unallocated за одну операцию
                                  const newUnallocated = unallocated
                                    .filter((t) => t.file !== newSubtitle.file)
                                    .concat(subtitle);

                                  setUnallocated(newUnallocated);
                                  setMatches(newMatches);
                                }}
                              >
                                <option value={subtitle.file}>{subtitle.file}</option>
                                {unallocated
                                  .filter((t) => t.type === 'subtitle')
                                  .map((track, i) => (
                                    <option key={i} value={track.file}>
                                      {track.file}
                                    </option>
                                  ))}
                              </select>
                            </div>
                            <Button
                              variant="outline-danger"
                              size="sm"
                              onClick={() => handleRemoveSubtitle(index, idx, subtitle)}
                            >
                              <Trash2 size={16} />
                            </Button>
                          </div>
                        </div>
                      ))}
                      {/* Кнопка добавления новых субтитров */}
                      {unallocated.some((t) => t.type === 'subtitle') && (
                        <Button
                          variant="outline-primary"
                          size="sm"
                          className={`ms-4 ${content.subtitle_tracks.length === 0 ? 'mt-2' : ''}`}
                          onClick={() => {
                            const newSubtitle = unallocated.find((t) => t.type === 'subtitle');
                            if (!newSubtitle) return;

                            const newMatches = matches.map((m, idx) => {
                              if (idx === index) {
                                return {
                                  ...m,
                                  subtitle_tracks: [...m.subtitle_tracks, newSubtitle],
                                };
                              }
                              return m;
                            });

                            setUnallocated(unallocated.filter((t) => t.file !== newSubtitle.file));
                            setMatches(newMatches);
                          }}
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
        ))}

        {unallocated.length > 0 && (
          <div className="border-top py-3">
            <Accordion>
              <Accordion.Item eventKey="0" className="border-0">
                <Accordion.Header>
                  <div className="d-flex align-items-center gap-2">
                    <h6 className="mb-0">Unallocated files ({unallocated.length})</h6>
                  </div>
                </Accordion.Header>
                <Accordion.Body className="pt-2 pb-1">
                  {/* Видео файлы */}
                  {unallocated.filter((t) => t.type === 'video').length > 0 && (
                    <div className="mb-3">
                      <div className="d-flex align-items-center gap-2 mb-2">
                        <Video size={18} className="text-primary" />
                        <span>Video files</span>
                      </div>
                      {unallocated
                        .filter((t) => t.type === 'video')
                        .map((track, idx) => (
                          <div key={idx} className="ps-4 mb-2 d-flex align-items-center">
                            <div className="text-break">{track.file}</div>
                          </div>
                        ))}
                    </div>
                  )}

                  {/* Аудио файлы */}
                  {unallocated.filter((t) => t.type === 'audio').length > 0 && (
                    <div className="mb-3">
                      <div className="d-flex align-items-center gap-2 mb-2">
                        <Mic2 size={18} className="text-primary" />
                        <span>Audio files</span>
                      </div>
                      {unallocated
                        .filter((t) => t.type === 'audio')
                        .map((track, idx) => (
                          <div key={idx} className="ps-4 mb-2 d-flex align-items-center">
                            <div className="text-break">{track.name || 'Unnamed track'}</div>
                            <div className="text-secondary ms-2">{track.file}</div>
                          </div>
                        ))}
                    </div>
                  )}

                  {/* Субтитры */}
                  {unallocated.filter((t) => t.type === 'subtitle').length > 0 && (
                    <div className="mb-3">
                      <div className="d-flex align-items-center gap-2 mb-2">
                        <Subtitles size={18} className="text-primary" />
                        <span>Subtitle files</span>
                      </div>
                      {unallocated
                        .filter((t) => t.type === 'subtitle')
                        .map((track, idx) => (
                          <div key={idx} className="ps-4 mb-2 d-flex align-items-center">
                            <div className="text-break">{track.name || 'Unnamed track'}</div>
                            <div className="text-secondary ms-2">{track.file}</div>
                          </div>
                        ))}
                    </div>
                  )}
                </Accordion.Body>
              </Accordion.Item>
            </Accordion>
          </div>
        )}
      </div>

      <div className="d-flex justify-content-end mt-3">
        <Button
          variant="primary"
          onClick={() => onConfirm(matches, unallocated)}
          disabled={!matches.length || loading}
          className="d-inline-flex align-items-center gap-3"
        >
          {loading ? <Spinner animation="border" size="sm" /> : <CheckCircle size={24} />}
          Confirm
        </Button>
      </div>
      <style>
        {`
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
        `}
      </style>
    </div>
  );
};
