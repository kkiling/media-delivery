import { useSearchParams } from 'react-router-dom';
import { Alert, Container, Row, Col, Spinner } from 'react-bootstrap';
import { useEffect, useState } from 'react';
import { Api, TVShowShort } from '@/api/api';
import { TVShowCard } from '@/components/TVShowCard';

export default function Search() {
  const [searchParams] = useSearchParams();
  const [shows, setShows] = useState<TVShowShort[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  
  const query = searchParams.get('query') || '';

  useEffect(() => {
    if (!query.trim()) return;

    setLoading(true);
    setError(null);

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
        setShows(response.data.items || []);
      })
      .catch((error) => {
        console.error('Error:', error);
        setError('Failed to fetch TV shows. Please try again.');
      })
      .finally(() => {
        setLoading(false);
      });
  }, [query]);

  return (
    <Container className="mt-4">
      <h2>Search Results</h2>
      {query && (
        <Alert variant="info" className="mb-4">
          Showing results for: <strong>{query}</strong>
        </Alert>
      )}

      {loading && (
        <div className="text-center">
          <Spinner animation="border" role="status">
            <span className="visually-hidden">Loading...</span>
          </Spinner>
        </div>
      )}

      {error && (
        <Alert variant="danger">
          {error}
        </Alert>
      )}

      <Row xs={1} md={2} lg={3} className="g-4">
        {shows.map((show) => (
          <Col key={show.id}>
            <TVShowCard show={show} />
          </Col>
        ))}
      </Row>

      {!loading && !error && shows.length === 0 && query && (
        <Alert variant="info">
          No TV shows found for your search.
        </Alert>
      )}
    </Container>
  );
}
