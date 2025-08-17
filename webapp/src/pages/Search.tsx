import { useSearchParams } from 'react-router-dom';
import { Alert, Container } from 'react-bootstrap';
import { useEffect } from 'react';
import { Api } from '@/api/api';

export default function Search() {
  const [searchParams] = useSearchParams();
  const query = searchParams.get('query') || '';

  useEffect(() => {
    const api = new Api({
      baseApiParams: {
        headers: {
          'Content-Type': 'application/json',
          Accept: 'application/json',
        },
      },
    });

    api.v1
      .tvShowLibraryServiceSearchTvShow({
        query: query,
      })
      .then((response) => {
        console.log('Success:', response.data);
      })
      .catch((error) => {
        console.error('Error:', error);
      });
  }, [query]);

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
