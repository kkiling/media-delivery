import { useState, useEffect } from 'react';
import { Button, Accordion, Spinner, Modal, Alert } from 'react-bootstrap';
import { Video, Mic2, Subtitles, CheckCircle, Pencil } from 'lucide-react';
import AceEditor from 'react-ace';
import yaml from 'js-yaml';

import 'ace-builds/src-noconflict/mode-yaml';
import 'ace-builds/src-noconflict/theme-github';

export interface Episode {
  episode_number: number;
  season_number: number;
  episode_file: string;
}

export interface Track {
  name: string;
  file: string;
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

  const [showModal, setShowModal] = useState(false);
  const [editorValue, setEditorValue] = useState('');
  const [error, setError] = useState<string | null>(null);

  // Инициализация из props
  useEffect(() => {
    setMatches(contentMatches);
    setUnallocated(unallocatedTracks);
  }, [contentMatches, unallocatedTracks]);

  const handleEdit = () => {
    const yamlString = yaml.dump(
      {
        episodes: matches.map((c) => ({
          episode: c.episode.episode_file,
          video: c.video.file,
          audio: c.audio_tracks.map((a) => a.file),
          subtitles: c.subtitle_tracks.map((s) => s.file),
        })),
      },
      { lineWidth: -1 } // всегда в одну строку
    );

    setEditorValue(yamlString);
    setError(null);
    setShowModal(true);
  };

  const handleSave = () => {
    try {
      // eslint-disable-next-line @typescript-eslint/no-explicit-any
      const parsed = yaml.load(editorValue) as any;

      if (!parsed || !Array.isArray(parsed.episodes)) {
        throw new Error('The YAML must have an "episodes" section');
      }

      // 1. Собираем все допустимые файлы
      const allowedFiles = new Set<string>([
        ...matches.map((m) => m.video.file),
        ...matches.flatMap((m) => m.audio_tracks.map((a) => a.file)),
        ...matches.flatMap((m) => m.subtitle_tracks.map((s) => s.file)),
        ...unallocated.map((t) => t.file),
      ]);

      const usedFiles = new Set<string>();
      // eslint-disable-next-line @typescript-eslint/no-explicit-any
      const updatedEpisodes: ContentMatches[] = parsed.episodes.map((ep: any, idx: number) => {
        if (!ep.episode || !ep.video) {
          throw new Error(`Эпизод #${idx + 1}: поля "episode" и "video" обязательны.`);
        }

        const collectFile = (file: string, context: string) => {
          if (!allowedFiles.has(file)) {
            throw new Error(`Файл "${file}" (${context}) не входит в исходный список файлов.`);
          }
          if (usedFiles.has(file)) {
            throw new Error(`Файл "${file}" используется более чем в одном месте.`);
          }
          usedFiles.add(file);
          return { name: file, file };
        };

        const video = collectFile(ep.video, `video эпизода ${ep.episode}`);
        const audio_tracks = (ep.audio || []).map((f: string) =>
          collectFile(f, `audio эпизода ${ep.episode}`)
        );
        const subtitle_tracks = (ep.subtitles || []).map((f: string) =>
          collectFile(f, `subtitles эпизода ${ep.episode}`)
        );

        return {
          episode: {
            episode_number: idx + 1,
            season_number: 1,
            episode_file: ep.episode,
          },
          video,
          audio_tracks,
          subtitle_tracks,
        };
      });

      // 2. Всё, что осталось → в unallocated
      const updatedUnallocated: Track[] = Array.from(allowedFiles)
        .filter((f) => !usedFiles.has(f))
        .map((f) => ({ name: f, file: f }));

      // 3. Обновляем state
      setMatches(updatedEpisodes);
      setUnallocated(updatedUnallocated);

      setError(null);
      setShowModal(false);
    } catch (e: any) {
      setError(e.message || 'Ошибка парсинга YAML');
    }
  };

  return (
    <div className="container-fluid py-3">
      <div className="d-flex justify-content-between align-items-center mb-3">
        <h4 className="mb-0">Confirm file matching</h4>
        <Button variant="outline-secondary" onClick={handleEdit}>
          <Pencil size={18} className="me-1" /> Edit
        </Button>
      </div>

      <div className="border-top border-bottom">
        {matches.map((content, index) => (
          <div key={index} className={`py-3 ${index !== 0 ? 'border-top' : ''}`}>
            <div className="d-flex justify-content-between align-items-center mb-3">
              <h6 className="mb-0">{content.episode.episode_file}</h6>
            </div>

            <div className="d-flex align-items-center gap-2 mb-3 ps-2">
              <Video size={18} className="text-primary" />
              <div className="text-break">{content.video.file}</div>
            </div>

            <div className="accordion-custom">
              {content.audio_tracks.length > 0 && (
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
                          <div className="fw-semibold">{audio.name}</div>
                          <div className="text-muted small">{audio.file}</div>
                        </div>
                      ))}
                    </Accordion.Body>
                  </Accordion.Item>
                </Accordion>
              )}

              {content.subtitle_tracks.length > 0 && (
                <Accordion>
                  <Accordion.Item eventKey="0" className="border-0">
                    <Accordion.Header>
                      <div className="d-flex align-items-center gap-2">
                        <Subtitles size={18} className="text-primary" />
                        <span>Subtitles ({content.subtitle_tracks.length})</span>
                      </div>
                    </Accordion.Header>
                    <Accordion.Body className="pt-2 pb-1">
                      {content.subtitle_tracks.map((sub, idx) => (
                        <div key={idx} className="mb-3 ps-4">
                          <div className="fw-semibold">{sub.name}</div>
                          <div className="text-muted small">{sub.file}</div>
                        </div>
                      ))}
                    </Accordion.Body>
                  </Accordion.Item>
                </Accordion>
              )}
            </div>
          </div>
        ))}

        {unallocated.length > 0 && (
          <div className="border-top py-3">
            <h6>Unallocated files</h6>
            {unallocated.map((t, idx) => (
              <div key={idx} className="text-muted small ps-2">
                {t.file}
              </div>
            ))}
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

      {/* Модалка редактора */}
      <Modal show={showModal} onHide={() => setShowModal(false)} fullscreen>
        <Modal.Header closeButton>
          <Modal.Title>Edit episodes</Modal.Title>
        </Modal.Header>
        <Modal.Body className="d-flex flex-row gap-3" style={{ height: '100%' }}>
          <div className="flex-fill">
            {error && <Alert variant="danger">{error}</Alert>}
            <AceEditor
              mode="yaml"
              theme="github"
              name="yamlEditor"
              value={editorValue}
              onChange={setEditorValue}
              fontSize={18}
              width="100%"
              height="calc(100vh - 200px)"
              setOptions={{ useWorker: false }}
            />
          </div>
          <div className="border-start ps-3" style={{ width: '30%', overflowY: 'auto' }}>
            <h6>Unallocated files</h6>
            {unallocated.length === 0 ? (
              <div className="text-muted small">No unallocated files</div>
            ) : (
              unallocated.map((t, idx) => (
                <div key={idx} className="text-muted small mb-1">
                  {t.file}
                </div>
              ))
            )}
          </div>
        </Modal.Body>
        <Modal.Footer>
          <Button variant="secondary" onClick={() => setShowModal(false)}>
            Cancel
          </Button>
          <Button variant="primary" onClick={handleSave}>
            Save
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
        `}
      </style>
    </div>
  );
};
