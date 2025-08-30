import { TorrentSearch } from '@/api/api';
import { Download, HardDrive, Upload, Users } from 'lucide-react';

interface TorrentCardProps {
  torrent: TorrentSearch;
  clickable?: boolean;
  onSelect?: (href: string) => void;
  isLast?: boolean;
}

export const TorrentCard = ({
  torrent,
  clickable = true,
  onSelect,
  isLast = false,
}: TorrentCardProps) => {
  const handleClick = (e: React.MouseEvent) => {
    e.preventDefault();
    if (clickable && onSelect && torrent.href) {
      onSelect(torrent.href);
    }
  };

  return (
    <div className="px-2 py-3">
      <div>
        <h5 className="mb-1">
          {clickable && torrent.href ? (
            <a
              href="#"
              onClick={handleClick}
              className="text-primary text-decoration-none hover-underline"
              style={{ cursor: 'pointer' }}
            >
              {torrent.title}
            </a>
          ) : (
            <span>{torrent.title}</span>
          )}
        </h5>
        <div className="text-muted small">
          {torrent.category} • {torrent.added_date}
        </div>
      </div>
      <div className="d-flex flex-wrap gap-2 mt-2">
        <span className="badge bg-light text-dark d-flex align-items-center gap-1">
          <HardDrive size={18} /> {torrent.size}
        </span>
        <span className="badge bg-light text-dark d-flex align-items-center gap-1">
          <Download size={18} /> {torrent.downloads}
        </span>
        <span className="badge bg-success text-white d-flex align-items-center gap-1">
          <Upload size={18} /> {torrent.seeds}
        </span>
        <span className="badge bg-danger text-white d-flex align-items-center gap-1">
          <Users size={18} /> {torrent.leeches}
        </span>
      </div>
      {!isLast && <hr className="my-3" />}
    </div>
  );
};
