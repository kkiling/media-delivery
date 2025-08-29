import { TorrentState, TorrentDownloadStatus } from '@/api/api';
import ProgressBar from 'react-bootstrap/ProgressBar';
import { Download, Upload, Pause, AlertCircle, HelpCircle, Loader } from 'lucide-react';
import Loading from '../Loading';

interface TorrentDownloadProgressProps {
  status?: TorrentDownloadStatus;
}

const getStatusIcon = (state?: TorrentState) => {
  switch (state) {
    case TorrentState.TORRENT_STATE_DOWNLOADING:
      return <Download size={20} className="text-primary" />;
    case TorrentState.TORRENT_STATE_UPLOADING:
      return <Upload size={20} className="text-success" />;
    case TorrentState.TORRENT_STATE_STOPPED:
      return <Pause size={20} className="text-warning" />;
    case TorrentState.TORRENT_STATE_ERROR:
      return <AlertCircle size={20} className="text-danger" />;
    default:
      return <HelpCircle size={20} className="text-muted" />;
  }
};

const getStatusText = (state?: TorrentState) => {
  switch (state) {
    case TorrentState.TORRENT_STATE_DOWNLOADING:
      return 'Downloading';
    case TorrentState.TORRENT_STATE_UPLOADING:
      return 'Seeding';
    case TorrentState.TORRENT_STATE_STOPPED:
      return 'Paused';
    case TorrentState.TORRENT_STATE_ERROR:
      return 'Error';
    case TorrentState.TORRENT_STATE_QUEUED:
      return 'Queued';
    default:
      return 'Unknown';
  }
};

export const TorrentDownloadProgress = ({ status }: TorrentDownloadProgressProps) => {
  if (!status) {
    return <p className="text-muted">No torrent download status available</p>;
  }

  // Округляем до 2 знаков после запятой
  const progress = Number((status.progress || 0) * 100).toFixed(2);

  return (
    <div className="mb-3">
      <div className="d-flex align-items-center gap-2 mb-2">
        {getStatusIcon(status.state)}
        <span>{getStatusText(status.state)}</span>
        <span className="ms-auto">{progress}%</span>
      </div>
      <ProgressBar
        now={Number(progress)}
        variant={status.state === TorrentState.TORRENT_STATE_ERROR ? 'danger' : 'primary'}
        animated={status.state === TorrentState.TORRENT_STATE_DOWNLOADING}
      />
    </div>
  );
};

interface TorrentWaitingFilesProps {
  status?: TorrentDownloadStatus;
}

export const TorrentWaitingFiles = ({ status }: TorrentWaitingFilesProps) => {
  return (
    <div className="mb-3">
      <div className="text-center mb-3">
        <Loading text="Waiting torrent files" />
      </div>
      {status && (
        <div className="d-flex align-items-center gap-2 justify-content-center">
          {getStatusIcon(status.state)}
          <span>{getStatusText(status.state)}</span>
        </div>
      )}
    </div>
  );
};
