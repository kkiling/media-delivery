import { useParams, Link } from 'react-router-dom';
import { ROUTES } from '@/constants/routes';

export default function TvShowSeason() {
  const { id, season } = useParams();

  return (
    <div>
      <h3>TV Show ID: {id}</h3>
      <h4>Season: {season}</h4>
      <Link to={ROUTES.LIBRARY.TV_SHOWS.getDetails(id || '')}>Back to TV Show</Link>
    </div>
  );
}
