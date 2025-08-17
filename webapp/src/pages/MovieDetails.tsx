import { useParams } from 'react-router-dom';

export default function MovieDetails() {
  const { id } = useParams();
  return <h3>Movie ID: {id}</h3>;
}
