import { useParams, Link } from 'react-router-dom';

export default function TvShowSeason() {
  const { id, season } = useParams();
  return (
    <div>
      <h3>TV Show ID: {id}</h3>
      <h4>Season: {season}</h4>
      <Link to={`/library/tvshows/${id}`}>Back to TV Show</Link>
    </div>
  );
}
