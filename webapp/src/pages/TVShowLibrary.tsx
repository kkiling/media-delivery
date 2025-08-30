import { Container, Row, Col, Spinner, Alert } from 'react-bootstrap';
import { useEffect } from 'react';
import { observer } from 'mobx-react-lite';
import { TVShowCard } from '@/components/TVShowCard';
import { libraryTVShowsStore } from '@/stores/libraryTVShowsStore';
import { Collection as CollectionIcon } from 'react-bootstrap-icons';

const TVShowLibrary = observer(() => {
  useEffect(() => {
    libraryTVShowsStore.loadLibraryShows();
  }, []);

  const { loading, error, shows } = libraryTVShowsStore;

  return (
    <Container className="mt-4">
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
            <div className="text-center py-5">
              <CollectionIcon size={48} className="text-muted mb-3" />
              <h4 className="text-muted">Your library is empty</h4>
              <p className="text-muted mb-0">Search for TV shows and add them to your library.</p>
            </div>
          )}
        </>
      )}
    </Container>
  );
});

export default TVShowLibrary;
