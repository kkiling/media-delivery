import { useParams, Link, Navigate } from 'react-router-dom';
import { ROUTES } from '@/constants/routes';

export default function TvShowSeason() {
  const { id, season } = useParams<{ id: string; season: string }>();

  const numberId = id ? parseInt(id, 10) : null;
  const numberSeason = season ? parseInt(season, 10) : null;

  if (!numberId || isNaN(numberId) || !numberSeason || isNaN(numberSeason)) {
    return <Navigate to={ROUTES.NOT_FOUND} />;
  }

  return (
    <div>
      <h3>TV Show ID: {numberId}</h3>
      <h4>Season: {numberSeason}</h4>
      <Link to={ROUTES.LIBRARY.TV_SHOWS.getDetails(numberId)}>Back to TV Show</Link>
    </div>
  );
}
