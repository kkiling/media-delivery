import React from 'react';
import { Accordion, Form, Button } from 'react-bootstrap';
import { Settings, Trash2 } from 'lucide-react';
import { Options } from '@/api/api';

interface OptionsPanelProps {
  options: Options;
  setOptions: React.Dispatch<React.SetStateAction<Options>>;
  uniqueAudioNames: string[];
  uniqueSubtitleNames: string[];
  removeTracksByName: (type: 'audio' | 'subtitle', name: string) => void;
}

export const OptionsPanel: React.FC<OptionsPanelProps> = ({
  options,
  setOptions,
  uniqueAudioNames,
  uniqueSubtitleNames,
  removeTracksByName,
}) => (
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
          {uniqueAudioNames.length > 0 && (
            <div>
              <label className="form-label">Remove Audio Tracks</label>
              <div className="d-flex flex-wrap gap-2">
                {uniqueAudioNames.map((name) => (
                  <Button
                    key={name}
                    variant="outline-danger"
                    size="sm"
                    onClick={() => removeTracksByName('audio', name)}
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
                    onClick={() => removeTracksByName('subtitle', name)}
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
);
