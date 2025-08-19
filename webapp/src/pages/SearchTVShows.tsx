import { useSearchParams } from 'react-router-dom';
import { Alert, Container, Row, Col, Spinner } from 'react-bootstrap';
import { useEffect } from 'react';
import { observer } from 'mobx-react-lite';
import { MainSearch, Rating } from '@/components';
import { searchTVShowsStore } from '@/stores/searchTVShowsStore';
import { Search as SearchIcon } from 'react-bootstrap-icons';
import { Card } from 'react-bootstrap';
import { useNavigate } from 'react-router-dom';
import { ROUTES } from '@/constants/routes';
import { TVShowShort } from '@/api/api';
import { useState } from 'react';
import { Image as ImageIcon } from 'react-bootstrap-icons';
import { formatDate } from '@/utils/formatting';

const CARD_CONFIG = {
  IMAGE_HEIGHT: 400,
} as const;

type TVShowCardProps = {
  show: TVShowShort;
};

function TVShowCard({ show }: TVShowCardProps) {
  const navigate = useNavigate();
  const [imageError, setImageError] = useState(false);
  return (
    <Card
      className="h-100 position-relative cursor-pointer"
      onClick={() => navigate(ROUTES.LIBRARY.TV_SHOWS.getDetails(show.id))}
    >
      <div className="position-relative" style={{ height: `${CARD_CONFIG.IMAGE_HEIGHT}px` }}>
        {show.poster?.w342 && !imageError ? (
          <Card.Img
            variant="top"
            src={show.poster.w342}
            alt={show.name}
            onError={() => setImageError(true)}
            className="w-100 h-100 object-fit-cover"
          />
        ) : (
          <Card.Body className="d-flex flex-column align-items-center justify-content-center bg-secondary text-white w-100 h-100">
            <ImageIcon size={48} className="mb-2" />
            <p className="mb-0">No Image Available</p>
          </Card.Body>
        )}

        {show.vote_average !== undefined && (
          <div className="position-absolute top-0 start-0 m-2">
            <Rating voteAverage={show.vote_average} voteCount={0} showVoteCount={false} />
          </div>
        )}
      </div>

      <Card.Body className="d-flex flex-column">
        <Card.Title>{show.name || 'No title'}</Card.Title>
        <Card.Subtitle className="mb-2 text-muted">
          {formatDate(show.first_air_date) || 'Release date unknown'}
        </Card.Subtitle>
        <Card.Text
          className="flex-grow-1"
          style={{
            overflow: 'hidden',
            textOverflow: 'ellipsis',
            display: '-webkit-box',
            WebkitLineClamp: 5,
            WebkitBoxOrient: 'vertical',
          }}
        >
          {show.overview || 'No overview available'}
        </Card.Text>
      </Card.Body>
    </Card>
  );
}

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
