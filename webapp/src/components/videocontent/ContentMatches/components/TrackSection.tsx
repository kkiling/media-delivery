import React from 'react';
import { Button } from 'react-bootstrap';
import { Trash2 } from 'lucide-react';
import { Track, TrackType } from '@/api/api';

interface TrackSectionProps {
  tracks: Track[];
  onRemove: (index: number, track: Track) => void;
  onReplace: (index: number, newTrackPath: string) => void;
  unallocatedTracks: Track[];
  type: TrackType;
  showName?: boolean;
}

export const TrackSection: React.FC<TrackSectionProps> = ({
  tracks,
  onRemove,
  onReplace,
  unallocatedTracks,
  type,
  showName = false,
}) => {
  if (tracks.length === 0) return null;

  return (
    <div className="mb-3">
      {tracks.map((track, index) => (
        <div key={index} className="mb-2 ps-3">
          <div className="d-flex align-items-center gap-2">
            {showName && (
              <div className="track-name text-truncate d-none d-md-block">
                {track.name || 'Unnamed track'}
              </div>
            )}
            <div className="flex-grow-1">
              <select
                className="form-select"
                value={track.relative_path}
                onChange={(e) => onReplace(index, e.target.value)}
              >
                <option value={track.relative_path}>{track.relative_path}</option>
                {unallocatedTracks
                  .filter((t) => t.type === type && t.name !== track.name)
                  .map((t, i) => (
                    <option key={i} value={t.relative_path}>
                      {t.relative_path}
                    </option>
                  ))}
              </select>
            </div>
            <Button variant="outline-danger" size="sm" onClick={() => onRemove(index, track)}>
              <Trash2 size={16} />
            </Button>
          </div>
        </div>
      ))}
    </div>
  );
};
