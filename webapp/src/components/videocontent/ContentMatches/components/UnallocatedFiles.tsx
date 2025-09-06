import React from 'react';
import { Accordion } from 'react-bootstrap';
import { Video, Mic2, Subtitles } from 'lucide-react';
import { Track, TrackType } from '@/api/api';
import { FileSection } from './FileSection';

interface UnallocatedFilesProps {
  unallocated: Track[];
}

export const UnallocatedFiles: React.FC<UnallocatedFilesProps> = ({ unallocated }) => {
  if (unallocated.length === 0) return null;

  const videoFiles = unallocated.filter((t) => t.type === TrackType.TRACK_TYPE_VIDEO);
  const audioFiles = unallocated.filter((t) => t.type === TrackType.TRACK_TYPE_AUDIO);
  const subtitleFiles = unallocated.filter((t) => t.type === TrackType.TRACK_TYPE_SUBTITLE);

  return (
    <div className="border-top py-3">
      <Accordion>
        <Accordion.Item eventKey="0" className="border-0">
          <Accordion.Header>
            <h6 className="mb-0">Unallocated files ({unallocated.length})</h6>
          </Accordion.Header>
          <Accordion.Body>
            {videoFiles.length > 0 && (
              <FileSection
                title="Video files"
                icon={<Video size={18} className="text-primary" />}
                files={videoFiles}
              />
            )}
            {audioFiles.length > 0 && (
              <FileSection
                title="Audio files"
                icon={<Mic2 size={18} className="text-primary" />}
                files={audioFiles}
                showName
              />
            )}
            {subtitleFiles.length > 0 && (
              <FileSection
                title="Subtitle files"
                icon={<Subtitles size={18} className="text-primary" />}
                files={subtitleFiles}
                showName
              />
            )}
          </Accordion.Body>
        </Accordion.Item>
      </Accordion>
    </div>
  );
};
