import { useSearchParams } from 'react-router-dom';
import { Alert, Container } from 'react-bootstrap';

export default function Search() {
  const [searchParams] = useSearchParams();
  const query = searchParams.get('query') || '';

  return (
    <Container className="mt-4">
      <h2>Search Results</h2>
      {query && (
        <Alert variant="info">
          Showing results for: <strong>{query}</strong>
        </Alert>
      )}
      {/* Здесь позже можно добавить результаты поиска */}
    </Container>
  );
}
