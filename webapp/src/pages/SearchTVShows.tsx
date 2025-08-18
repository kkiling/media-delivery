import { useSearchParams } from 'react-router-dom';
import { Alert, Container, Row, Col, Spinner } from 'react-bootstrap';
import { useEffect } from 'react';
import { observer } from 'mobx-react-lite';
import { TVShowCard } from '@/components/TVShowCard';
import { MainSearch } from '@/components';
import { searchTVShowsStore } from '@/stores/searchTVShowsStore';
import { Search as SearchIcon } from 'react-bootstrap-icons';

const SearchTVShows = observer(() => {
  const [searchParams] = useSearchParams();
  const query = searchParams.get('query') || '';

  useEffect(() => {
    searchTVShowsStore.searchShows(query);
  }, [query]);

  const renderNoResults = () => (
    <div className="text-center py-5">
      <SearchIcon size={48} className="text-muted mb-3" />
      <h4 className="text-muted">No TV shows found</h4>
      <p className="text-muted mb-0">Try adjusting your search to find what you're looking for.</p>
    </div>
  );

  return (
    <Container className="mt-4">
      <div className="mb-4">
        <MainSearch mode="tvshows" query={query} />
      </div>

      {searchTVShowsStore.loading ? (
        <div className="text-center">
          <Spinner animation="border" role="status">
            <span className="visually-hidden">Loading...</span>
          </Spinner>
        </div>
      ) : (
        <>
          {searchTVShowsStore.error ? (
            <Alert variant="danger">{searchTVShowsStore.error}</Alert>
          ) : (
            <>
              {searchTVShowsStore.shows.length > 0 ? (
                <Row xs={1} md={2} lg={4} className="g-4">
                  {searchTVShowsStore.shows.map((show) => (
                    <Col key={show.id}>
                      <TVShowCard show={show} />
                    </Col>
                  ))}
                </Row>
              ) : (
                query && renderNoResults()
              )}
            </>
          )}
        </>
      )}
    </Container>
  );
});

export default SearchTVShows;
