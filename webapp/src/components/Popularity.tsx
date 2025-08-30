const POPULARITY_CONFIG = {
  MIN: 0,
  MAX: 500,
  BAR_HEIGHT: 12,
  THRESHOLDS: [
    { value: 10, color: 'secondary' },
    { value: 50, color: 'info' },
    { value: 200, color: 'primary' },
    { value: 400, color: 'warning' },
    { value: 500, color: 'danger' },
  ],
} as const;

interface PopularityProps {
  popularity: number;
}

export const Popularity = ({ popularity }: PopularityProps) => {
  const getPopularityPercentage = (value: number) => {
    // Using natural logarithm for smooth scaling
    const logMin = Math.log(POPULARITY_CONFIG.MIN + 1); // +1 to avoid log(0)
    const logMax = Math.log(POPULARITY_CONFIG.MAX + 1);
    const logValue = Math.log(value + 1);

    const percentage = ((logValue - logMin) / (logMax - logMin)) * 100;
    return Math.min(Math.max(percentage, 0), 100);
  };

  const getPopularityColor = (value: number) => {
    const threshold =
      POPULARITY_CONFIG.THRESHOLDS.find((t) => value <= t.value) ||
      POPULARITY_CONFIG.THRESHOLDS[POPULARITY_CONFIG.THRESHOLDS.length - 1];
    return threshold.color;
  };

  const percentage = getPopularityPercentage(popularity);
  const color = getPopularityColor(popularity);

  return (
    <div>
      <div className="text-center mb-1">
        <small className="text-muted" style={{ fontSize: '0.8rem' }}>
          Popularity
        </small>
      </div>
      <div className="progress" style={{ height: `${POPULARITY_CONFIG.BAR_HEIGHT}px` }}>
        <div
          className={`progress-bar bg-${color}`}
          role="progressbar"
          style={{ width: `${percentage}%` }}
          aria-valuenow={percentage}
          aria-valuemin={0}
          aria-valuemax={100}
        />
      </div>
    </div>
  );
};

export default Popularity;
