import { useSearchParams } from 'react-router-dom';
import { Alert, Container, Row, Col, Spinner } from 'react-bootstrap';
import { useEffect } from 'react';
import { observer } from 'mobx-react-lite';
import { MainSearch } from '@/components';
import { TVShowCard } from '@/components/TVShowCard';
import { searchTVShowsStore } from '@/stores/searchTVShowsStore';
import { Search as SearchIcon } from 'react-bootstrap-icons';

const SearchTVShows = observer(() => {
  const [searchParams] = useSearchParams();
  const query = searchParams.get('query')?.trim() || '';

  useEffect(() => {
    if (query) {
      searchTVShowsStore.searchShows(query);
    }
  }, [query]);

  const { loading, error, shows } = searchTVShowsStore;

  return (
    <Container className="mt-4">
      <div className="mb-4">
        <MainSearch mode="tvshows" query={query} />
      </div>

      {loading && (
        <div className="text-center">
          <Spinner animation="border" role="status">
            <span className="visually-hidden">Loading...</span>
          </Spinner>
        </div>
      )}

      {!loading && error && <Alert variant="danger">{error}</Alert>}

      {!loading && !error && (
        <>
          {shows.length > 0 ? (
            <Row xs={1} md={2} lg={4} className="g-4">
              {shows.map((show) => (
                <Col key={show.id}>
                  <TVShowCard show={show} />
                </Col>
              ))}
            </Row>
          ) : (
            query && (
              <div className="text-center py-5">
                <SearchIcon size={48} className="text-muted mb-3" />
                <h4 className="text-muted">No TV shows found</h4>
                <p className="text-muted mb-0">
                  Try adjusting your search to find what you're looking for.
                </p>
              </div>
            )
          )}
        </>
      )}
    </Container>
  );
});

export default SearchTVShows;
