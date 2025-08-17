import { useSearchParams } from 'react-router-dom';
import { Alert, Container } from 'react-bootstrap';

export default function SearchMovies() {
  const [searchParams] = useSearchParams();
  const query = searchParams.get('query') || '';

  return (
    <Container className="mt-4">
      <h2>Movie Search Results</h2>
      <Alert variant="info">
        Searching for movies with query: <strong>{query}</strong>
      </Alert>
    </Container>
  );
}
