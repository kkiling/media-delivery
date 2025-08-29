import { MergeVideoStatus } from '@/api/api';
import ProgressBar from 'react-bootstrap/ProgressBar';
import { FileVideo, CheckCircle2 } from 'lucide-react';

interface MergeVideoProgressProps {
  status?: MergeVideoStatus;
}

export const MergeVideoProgress = ({ status }: MergeVideoProgressProps) => {
  if (!status) {
    return <p className="text-muted">No merge video status available</p>;
  }

  // Round to 2 decimal places
  const progress = Number((status.progress || 0) * 100).toFixed(2);

  return (
    <div className="mb-3">
      <div className="d-flex align-items-center gap-2 mb-2">
        {status.is_complete ? (
          <CheckCircle2 size={20} className="text-success" />
        ) : (
          <FileVideo size={20} className="text-primary" />
        )}
        <span>{status.is_complete ? 'Complete' : 'Processing files'}</span>
        <span className="ms-auto">{progress}%</span>
      </div>
      <ProgressBar
        now={Number(progress)}
        variant={status.is_complete ? 'success' : 'primary'}
        animated={!status.is_complete}
      />
    </div>
  );
};
