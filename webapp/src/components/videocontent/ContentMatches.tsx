import { Button, Accordion, Spinner } from 'react-bootstrap';
import { Video, Mic2, Subtitles, CheckCircle } from 'lucide-react';

export interface Episode {
  episode_number: number;
  season_number: number;
  episode_file: string;
}

export interface Track {
  name: string;
  file: string;
}

export interface MediaFiles {
  video_file: string;
  audio_files: Track[];
  subtitles_files: Track[];
}

export interface ContentMatches {
  episode: Episode;
  media_files: MediaFiles;
}

export interface ContentMatchesProps {
  loading: boolean;
  contentMatches?: ContentMatches[];
  onConfirm: () => void;
}

export const ContentMatche = ({ contentMatches, onConfirm, loading }: ContentMatchesProps) => {
  return (
    <div className="container-fluid py-3">
      <h4 className="mb-3">Confirm file matching</h4>

      <div className="border-top border-bottom">
        {contentMatches?.map((content, index) => (
          <div key={index} className={`py-3 ${index !== 0 ? 'border-top' : ''}`}>
            <div className="d-flex justify-content-between align-items-center mb-3">
              <h6 className="mb-0">{content.episode.episode_file}</h6>
            </div>

            <div className="d-flex align-items-center gap-2 mb-3 ps-2">
              <Video size={18} className="text-primary" />
              <div className="text-break">{content.media_files.video_file}</div>
            </div>

            <div className="accordion-custom">
              {content.media_files.audio_files.length > 0 && (
                <Accordion>
                  <Accordion.Item eventKey="0" className="border-0 mb-2">
                    <Accordion.Header>
                      <div className="d-flex align-items-center gap-2">
                        <Mic2 size={18} className="text-primary" />
                        <span>Audio files ({content.media_files.audio_files.length})</span>
                      </div>
                    </Accordion.Header>
                    <Accordion.Body className="pt-2 pb-1">
                      {content.media_files.audio_files.map((audio, idx) => (
                        <div key={idx} className="mb-3 ps-4">
                          <div className="fw-semibold">{audio.name}</div>
                          <div className="text-muted small">{audio.file}</div>
                        </div>
                      ))}
                    </Accordion.Body>
                  </Accordion.Item>
                </Accordion>
              )}

              {content.media_files.subtitles_files.length > 0 && (
                <Accordion>
                  <Accordion.Item eventKey="0" className="border-0">
                    <Accordion.Header>
                      <div className="d-flex align-items-center gap-2">
                        <Subtitles size={18} className="text-primary" />
                        <span>Subtitles ({content.media_files.subtitles_files.length})</span>
                      </div>
                    </Accordion.Header>
                    <Accordion.Body className="pt-2 pb-1">
                      {content.media_files.subtitles_files.map((sub, idx) => (
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
      </div>
      <div className="d-flex justify-content-end mt-3">
        <Button
          variant="primary"
          onClick={onConfirm}
          disabled={!contentMatches?.length || loading}
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
        `}
      </style>
    </div>
  );
};
