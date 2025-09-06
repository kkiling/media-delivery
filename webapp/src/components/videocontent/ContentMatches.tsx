import { useState, useEffect, useCallback } from 'react';
import { Button, Accordion, Spinner, Form, Modal } from 'react-bootstrap';
import { Video, Mic2, Subtitles, CheckCircle, Trash2, Settings } from 'lucide-react';
import {
  ContentMatch,
  ContentMatches as MatchContents,
  Options,
  Track,
  TrackType,
} from '@/api/api';

export interface ContentMatchesProps {
  loading: boolean;
  contentMatches: MatchContents;
  onConfirm: (contentMatches?: MatchContents) => void;
}

const sortTracks = (tracks: Track[]) => {
  return [...tracks].sort((a, b) => {
    // Сортировка по типу (видео -> аудио -> субтитры)
    if (a.type !== b.type) {
      const typeOrder = {
        [TrackType.TRACK_TYPE_VIDEO]: 0,
        [TrackType.TRACK_TYPE_AUDIO]: 1,
        [TrackType.TRACK_TYPE_SUBTITLE]: 2,
      };
      const aType = a.type ?? '';
      const bType = b.type ?? '';
      return (
        (typeOrder[aType as keyof typeof typeOrder] ?? 3) -
        (typeOrder[bType as keyof typeof typeOrder] ?? 3)
      );
    }

    // Если имена треков одинаковые - сортируем по relative_path
    if (a.name === b.name) {
      return (a.relative_path || '').localeCompare(b.relative_path || '');
    }

    // Сортировка по имени трека
    return (a.name || '').localeCompare(b.name || '');
  });
};

export const ContentMatches = ({ contentMatches, onConfirm, loading }: ContentMatchesProps) => {
  const [matches, setMatches] = useState<ContentMatch[]>([]);
  const [unallocated, setUnallocated] = useState<Track[]>([]);
  const [options, setOptions] = useState<Options>({});
  const [showConfirmModal, setShowConfirmModal] = useState(false);

  // Инициализация из props
  useEffect(() => {
    setMatches(contentMatches.matches || []);
    setUnallocated(sortTracks(contentMatches.unallocated || []));
    setOptions(contentMatches.options || {});
  }, [contentMatches]);

  // Устанавливаем первый аудио трек если не выбран
  useEffect(() => {
    if (
      matches.length > 0 &&
      !options.default_audio_track_name &&
      matches.some((m) => m.audio_files && m.audio_files.length > 0)
    ) {
      const firstAudioTrack = matches.flatMap((m) => m.audio_files || []).find((t) => t.name)?.name;
      if (firstAudioTrack) {
        setOptions((s) => ({
          ...s,
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
        .map((t) => t.name)
        .filter(Boolean)
    )
  );
  const uniqueSubtitleNames = Array.from(
    new Set(
      matches
        .flatMap((m) => m.subtitles || [])
        .map((t) => t.name)
        .filter(Boolean)
    )
  );

  // Удаление треков по имени (без состояния)
  const removeTracksByName = useCallback(
    (type: 'audio' | 'subtitle', name: string) => {
      const removedTracks: Track[] = [];
      const newMatches = matches.map((match) => {
        if (type === 'audio') {
          const tracksToRemove = (match.audio_files || []).filter((t) => t.name === name);
          removedTracks.push(...tracksToRemove);
          return {
            ...match,
            audio_files: (match.audio_files || []).filter((t) => t.name !== name),
          };
        } else {
          const tracksToRemove = (match.subtitles || []).filter((t) => t.name === name);
          removedTracks.push(...tracksToRemove);
          return {
            ...match,
            subtitles: (match.subtitles || []).filter((t) => t.name !== name),
          };
        }
      });
      setUnallocated(sortTracks([...unallocated, ...removedTracks]));
      setMatches(newMatches);
    },
    [matches, unallocated]
  );

  // Удаление видео/аудио/субтитров по индексу
  const handleRemoveVideo = (index: number, video: Track) => {
    const newMatches = matches.map((m, idx) =>
      idx === index ? { ...m, video: { file: '', type: TrackType.TRACK_TYPE_VIDEO } } : m
    );
    setUnallocated(sortTracks([...unallocated, video]));
    setMatches(newMatches);
  };

  const handleRemoveAudio = (contentIndex: number, audioIndex: number, audio: Track) => {
    const newMatches = matches.map((m, idx) => {
      if (idx === contentIndex) {
        const newAudioTracks = [...(m.audio_files || [])];
        newAudioTracks.splice(audioIndex, 1);
        return { ...m, audio_files: newAudioTracks };
      }
      return m;
    });
    setUnallocated(sortTracks([...unallocated, audio]));
    setMatches(newMatches);
  };

  const handleRemoveSubtitle = (contentIndex: number, subtitleIndex: number, subtitle: Track) => {
    const newMatches = matches.map((m, idx) => {
      if (idx === contentIndex) {
        const newSubtitles = [...(m.subtitles || [])];
        newSubtitles.splice(subtitleIndex, 1);
        return { ...m, subtitles: newSubtitles };
      }
      return m;
    });
    setUnallocated(sortTracks([...unallocated, subtitle]));
    setMatches(newMatches);
  };

  const handleConfirmClick = () => setShowConfirmModal(true);
  const handleConfirmClose = () => setShowConfirmModal(false);

  const handleConfirmSubmit = () => {
    onConfirm({
      matches,
      unallocated,
      options: {
        keep_original_audio: options.keep_original_audio,
        keep_original_subtitles: options.keep_original_subtitles,
        default_audio_track_name: options.default_audio_track_name,
        default_subtitle_track: options.default_subtitle_track,
      },
    });
    setShowConfirmModal(false);
  };

  return (
    <>
      <div className="container-fluid py-3">
        {/* Track Settings Panel */}
        <div className="mb-3">
          <Accordion>
            <Accordion.Item eventKey="0" className="border-0">
              <Accordion.Header>
                <div className="d-flex align-items-center gap-2">
                  <Settings size={18} className="text-primary" />
                  <h6 className="mb-0">Options</h6>
                </div>
              </Accordion.Header>
              <Accordion.Body>
                <div className="d-flex flex-column gap-3">
                  {/* Original Tracks */}
                  <div>
                    <Form.Check
                      type="checkbox"
                      id="keepOriginalAudio"
                      label="Source audio tracks"
                      checked={!!options.keep_original_audio}
                      onChange={(e) =>
                        setOptions((s) => ({
                          ...s,
                          keep_original_audio: e.target.checked,
                        }))
                      }
                    />
                    <Form.Check
                      type="checkbox"
                      id="keepOriginalSubtitles"
                      label="Source subtitles"
                      checked={!!options.keep_original_subtitles}
                      onChange={(e) =>
                        setOptions((s) => ({
                          ...s,
                          keep_original_subtitles: e.target.checked,
                        }))
                      }
                    />
                  </div>
                  {/* Audio Track Selection */}
                  <div>
                    <label className="form-label mb-1">Default audio track</label>
                    <Form.Select
                      size="sm"
                      value={options.default_audio_track_name || ''}
                      onChange={(e) =>
                        setOptions((s) => ({
                          ...s,
                          default_audio_track_name: e.target.value,
                        }))
                      }
                    >
                      <option value=""></option>
                      {uniqueAudioNames.map((name) => (
                        <option key={name} value={name}>
                          {name}
                        </option>
                      ))}
                    </Form.Select>
                  </div>
                  {/* Subtitle Track Selection */}
                  <div>
                    <label className="form-label mb-1">Default subtitle</label>
                    <Form.Select
                      size="sm"
                      value={options.default_subtitle_track || ''}
                      onChange={(e) =>
                        setOptions((s) => ({
                          ...s,
                          default_subtitle_track: e.target.value,
                        }))
                      }
                    >
                      <option value=""></option>
                      {uniqueSubtitleNames.map((name) => (
                        <option key={name} value={name}>
                          {name}
                        </option>
                      ))}
                    </Form.Select>
                  </div>
                  {/* Batch Remove */}
                  {uniqueAudioNames.length > 0 && (
                    <div>
                      <label className="form-label">Remove Audio Tracks</label>
                      <div className="d-flex flex-wrap gap-2">
                        {uniqueAudioNames.map((name) => (
                          <Button
                            key={name}
                            variant="outline-danger"
                            size="sm"
                            onClick={() => removeTracksByName('audio', name || '')}
                            className="d-inline-flex align-items-center gap-2"
                          >
                            <Trash2 size={14} />
                            {name}
                          </Button>
                        ))}
                      </div>
                    </div>
                  )}
                  {uniqueSubtitleNames.length > 0 && (
                    <div>
                      <label className="form-label">Remove Subtitles</label>
                      <div className="d-flex flex-wrap gap-2">
                        {uniqueSubtitleNames.map((name) => (
                          <Button
                            key={name}
                            variant="outline-danger"
                            size="sm"
                            onClick={() => removeTracksByName('subtitle', name || '')}
                            className="d-inline-flex align-items-center gap-2"
                          >
                            <Trash2 size={14} />
                            {name}
                          </Button>
                        ))}
                      </div>
                    </div>
                  )}
                </div>
              </Accordion.Body>
            </Accordion.Item>
          </Accordion>
        </div>
        <div className="border-top border-bottom">
          {matches.map((content, index) => (
            <div key={index} className={`py-3 ${index !== 0 ? 'border-top' : ''}`}>
              <div className="d-flex justify-content-between align-items-center mb-3">
                <h6 className="mb-0">{content.episode?.relative_path}</h6>
              </div>
              <div className="d-flex align-items-center gap-2 mb-3 ps-2">
                <Video size={18} className="text-primary" />
                <div className="d-flex align-items-center gap-2 flex-grow-1">
                  <select
                    className="form-select"
                    value={content.video?.relative_path || ''}
                    onChange={(e) => {
                      const newVideo = unallocated.find((t) => t.relative_path === e.target.value);
                      if (!newVideo) return;
                      const newMatches = matches.map((m, idx) =>
                        idx === index ? { ...m, video: newVideo } : m
                      );
                      const newUnallocated = unallocated
                        .filter((t) => t.relative_path !== newVideo.relative_path)
                        .concat(content.video ? content.video : []);
                      setUnallocated(newUnallocated);
                      setMatches(newMatches);
                    }}
                  >
                    <option value={content.video?.relative_path || ''}>
                      {content.video?.relative_path}
                    </option>
                    {unallocated
                      .filter((t) => t.type === TrackType.TRACK_TYPE_VIDEO)
                      .map((track, idx) => (
                        <option key={idx} value={track.relative_path}>
                          {track.relative_path}
                        </option>
                      ))}
                  </select>
                  <Button
                    variant="outline-danger"
                    size="sm"
                    disabled={!content.video?.relative_path}
                    onClick={() => handleRemoveVideo(index, content.video!)}
                  >
                    <Trash2 size={16} />
                  </Button>
                </div>
              </div>
              <div className="accordion-custom">
                {/* Аудио */}
                {(content.audio_files?.length || 0) > 0 ||
                unallocated.some(
                  (t) =>
                    t.type === TrackType.TRACK_TYPE_AUDIO &&
                    !(content.audio_files || []).some((a) => a.name === t.name)
                ) ? (
                  <Accordion>
                    <Accordion.Item eventKey="0" className="border-0 mb-2">
                      <Accordion.Header>
                        <div className="d-flex align-items-center gap-2">
                          <Mic2 size={18} className="text-primary" />
                          <span>Audio files ({content.audio_files?.length || 0})</span>
                        </div>
                      </Accordion.Header>
                      <Accordion.Body className="pt-2 pb-1">
                        {(content.audio_files || []).map((audio, idx) => (
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
                                  onChange={(e) => {
                                    const newAudio = unallocated.find(
                                      (t) => t.relative_path === e.target.value
                                    );
                                    if (!newAudio) return;
                                    const newMatches = matches.map((m, mIdx) => {
                                      if (mIdx === index) {
                                        const newAudioTracks = [...(m.audio_files || [])];
                                        newAudioTracks[idx] = newAudio;
                                        return { ...m, audio_files: newAudioTracks };
                                      }
                                      return m;
                                    });
                                    const newUnallocated = unallocated
                                      .filter((t) => t.relative_path !== newAudio.relative_path)
                                      .concat(audio);
                                    setUnallocated(newUnallocated);
                                    setMatches(newMatches);
                                  }}
                                >
                                  <option value={audio.relative_path}>{audio.relative_path}</option>
                                  {unallocated
                                    .filter((t) => t.type === TrackType.TRACK_TYPE_AUDIO)
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
                                onClick={() => handleRemoveAudio(index, idx, audio)}
                              >
                                <Trash2 size={16} />
                              </Button>
                            </div>
                          </div>
                        ))}
                        {unallocated.some(
                          (t) =>
                            t.type === TrackType.TRACK_TYPE_AUDIO &&
                            !(content.audio_files || []).some((a) => a.name === t.name)
                        ) && (
                          <Button
                            variant="outline-primary"
                            size="sm"
                            className={`ms-4 ${!content.audio_files?.length ? 'mt-2' : ''}`}
                            onClick={() => {
                              const newAudio = unallocated.find(
                                (t) =>
                                  t.type === TrackType.TRACK_TYPE_AUDIO &&
                                  !(content.audio_files || []).some((a) => a.name === t.name)
                              );
                              if (!newAudio) return;
                              const newMatches = matches.map((m, idx) =>
                                idx === index
                                  ? {
                                      ...m,
                                      audio_files: [...(m.audio_files || []), newAudio],
                                    }
                                  : m
                              );
                              setUnallocated(
                                unallocated.filter(
                                  (t) => t.relative_path !== newAudio.relative_path
                                )
                              );
                              setMatches(newMatches);
                            }}
                          >
                            Add Audio Track
                          </Button>
                        )}
                      </Accordion.Body>
                    </Accordion.Item>
                  </Accordion>
                ) : null}
                {/* Субтитры */}
                {(content.subtitles?.length || 0) > 0 ||
                unallocated.some(
                  (t) =>
                    t.type === TrackType.TRACK_TYPE_SUBTITLE &&
                    !(content.subtitles || []).some((s) => s.name === t.name)
                ) ? (
                  <Accordion>
                    <Accordion.Item eventKey="0" className="border-0">
                      <Accordion.Header>
                        <div className="d-flex align-items-center gap-2">
                          <Subtitles size={18} className="text-primary" />
                          <span>Subtitles ({content.subtitles?.length || 0})</span>
                        </div>
                      </Accordion.Header>
                      <Accordion.Body className="pt-2 pb-1">
                        {(content.subtitles || []).map((subtitle, idx) => (
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
                                  onChange={(e) => {
                                    const newSubtitle = unallocated.find(
                                      (t) => t.relative_path === e.target.value
                                    );
                                    if (!newSubtitle) return;
                                    const newMatches = matches.map((m, mIdx) => {
                                      if (mIdx === index) {
                                        const newSubtitles = [...(m.subtitles || [])];
                                        newSubtitles[idx] = newSubtitle;
                                        return { ...m, subtitles: newSubtitles };
                                      }
                                      return m;
                                    });
                                    const newUnallocated = unallocated
                                      .filter((t) => t.relative_path !== newSubtitle.relative_path)
                                      .concat(subtitle);
                                    setUnallocated(newUnallocated);
                                    setMatches(newMatches);
                                  }}
                                >
                                  <option value={subtitle.relative_path}>
                                    {subtitle.relative_path}
                                  </option>
                                  {unallocated
                                    .filter((t) => t.type === TrackType.TRACK_TYPE_SUBTITLE)
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
                                onClick={() => handleRemoveSubtitle(index, idx, subtitle)}
                              >
                                <Trash2 size={16} />
                              </Button>
                            </div>
                          </div>
                        ))}
                        {unallocated.some(
                          (t) =>
                            t.type === TrackType.TRACK_TYPE_SUBTITLE &&
                            !(content.subtitles || []).some((s) => s.name === t.name)
                        ) && (
                          <Button
                            variant="outline-primary"
                            size="sm"
                            className={`ms-4 ${!content.subtitles?.length ? 'mt-2' : ''}`}
                            onClick={() => {
                              const newSubtitle = unallocated.find(
                                (t) =>
                                  t.type === TrackType.TRACK_TYPE_SUBTITLE &&
                                  !(content.subtitles || []).some((s) => s.name === t.name)
                              );
                              if (!newSubtitle) return;
                              const newMatches = matches.map((m, idx) =>
                                idx === index
                                  ? {
                                      ...m,
                                      subtitles: [...(m.subtitles || []), newSubtitle],
                                    }
                                  : m
                              );
                              setUnallocated(
                                unallocated.filter(
                                  (t) => t.relative_path !== newSubtitle.relative_path
                                )
                              );
                              setMatches(newMatches);
                            }}
                          >
                            Add Subtitle Track
                          </Button>
                        )}
                      </Accordion.Body>
                    </Accordion.Item>
                  </Accordion>
                ) : null}
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
                    {unallocated.filter((t) => t.type === TrackType.TRACK_TYPE_VIDEO).length >
                      0 && (
                      <div className="mb-3">
                        <div className="d-flex align-items-center gap-2 mb-2">
                          <Video size={18} className="text-primary" />
                          <span>Video files</span>
                        </div>
                        {unallocated
                          .filter((t) => t.type === TrackType.TRACK_TYPE_VIDEO)
                          .map((track, idx) => (
                            <div key={idx} className="ps-4 mb-2 d-flex align-items-center">
                              <div className="text-break">{track.relative_path}</div>
                            </div>
                          ))}
                      </div>
                    )}
                    {/* Аудио файлы */}
                    {unallocated.filter((t) => t.type === TrackType.TRACK_TYPE_AUDIO).length >
                      0 && (
                      <div className="mb-3">
                        <div className="d-flex align-items-center gap-2 mb-2">
                          <Mic2 size={18} className="text-primary" />
                          <span>Audio files</span>
                        </div>
                        {unallocated
                          .filter((t) => t.type === TrackType.TRACK_TYPE_AUDIO)
                          .map((track, idx) => (
                            <div key={idx} className="ps-4 mb-2 d-flex align-items-center">
                              <div className="text-break">{track.name || 'Unnamed track'}</div>
                              <div className="text-secondary ms-2">{track.relative_path}</div>
                            </div>
                          ))}
                      </div>
                    )}
                    {/* Субтитры */}
                    {unallocated.filter((t) => t.type === TrackType.TRACK_TYPE_SUBTITLE).length >
                      0 && (
                      <div className="mb-3">
                        <div className="d-flex align-items-center gap-2 mb-2">
                          <Subtitles size={18} className="text-primary" />
                          <span>Subtitle files</span>
                        </div>
                        {unallocated
                          .filter((t) => t.type === TrackType.TRACK_TYPE_SUBTITLE)
                          .map((track, idx) => (
                            <div key={idx} className="ps-4 mb-2 d-flex align-items-center">
                              <div className="text-break">{track.name || 'Unnamed track'}</div>
                              <div className="text-secondary ms-2">{track.relative_path}</div>
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
        {/* Confirm button */}
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
      {/* Modal */}
      <Modal show={showConfirmModal} onHide={handleConfirmClose} centered>
        <Modal.Header closeButton>
          <Modal.Title>Confirm File Matching</Modal.Title>
        </Modal.Header>
        <Modal.Body>Are you sure you want to confirm the file matching?</Modal.Body>
        <Modal.Footer>
          <Button variant="secondary" onClick={handleConfirmClose}>
            Cancel
          </Button>
          <Button variant="primary" onClick={handleConfirmSubmit}>
            Confirm
          </Button>
        </Modal.Footer>
      </Modal>
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
    </>
  );
};
