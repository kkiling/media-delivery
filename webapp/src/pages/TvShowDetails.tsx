import { useParams, Link } from 'react-router-dom';
import { ROUTES } from '@/constants/routes';

export default function TvShowDetails() {
  const { id } = useParams();

  return (
    <div>
      <h3>TV Show ID: {id}</h3>
      <Link to={ROUTES.LIBRARY.TV_SHOWS.getSeason(id || '', 1)}>Season 1</Link>
    </div>
  );
}
