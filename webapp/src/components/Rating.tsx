import { getRatingColor } from '@/utils/formatting';

interface RatingProps {
  voteAverage: number;
  voteCount: number;
  showVoteCount?: boolean;
}

export const Rating = ({ voteAverage, voteCount, showVoteCount = true }: RatingProps) => {
  if (voteAverage === 0) {
    return null;
  }

  return (
    <div className="mb-3">
      <div
        className={`bg-${getRatingColor(voteAverage)} 
          text-white rounded-circle 
          d-flex align-items-center justify-content-center 
          mx-auto mb-1`}
        style={{
          width: 42,
          height: 42,
          border: '2px solid white',
          boxShadow: '0 2px 4px rgba(0,0,0,0.2)',
        }}
      >
        <div className="fw-bold" style={{ fontSize: '1.1rem' }}>
          {voteAverage.toFixed(1)}
        </div>
      </div>
      {showVoteCount && (
        <div className="text-center">
          <small className="text-muted" style={{ fontSize: '0.8rem' }}>
            {voteCount.toLocaleString()} votes
          </small>
        </div>
      )}
    </div>
  );
};

export default Rating;
