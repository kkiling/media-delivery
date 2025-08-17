import { useParams, Link } from 'react-router-dom';

export default function TvShowDetails() {
  const { id } = useParams();
  return (
    <div>
      <h3>TV Show ID: {id}</h3>
      <Link to={`/library/tvshows/${id}/1`}>Season 1</Link>
    </div>
  );
}
