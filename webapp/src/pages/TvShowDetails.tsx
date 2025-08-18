import { useParams, Navigate } from 'react-router-dom';
import { useEffect } from 'react';
import { observer } from 'mobx-react-lite';
import { Card, Container, Row, Col, Badge, Spinner, Alert } from 'react-bootstrap';
import { tvShowDetailsStore } from '@/stores/tvShowDetailsStore';
import { ROUTES } from '@/constants/routes';
import { formatDate, getRatingColor } from '@/utils/formatting';
import getCountryFlag from 'country-flag-icons/unicode';

const POPULARITY_CONFIG = {
  MIN: 0,
  MAX: 500,
  BAR_HEIGHT: 6,
} as const;

const RatingSection = ({ voteAverage, voteCount }: { voteAverage: number; voteCount: number }) => {
  return (
    <div className="mb-3">
      <div
        className={`bg-${getRatingColor(voteAverage)} 
          text-white rounded-circle 
          d-flex align-items-center justify-content-center 
          mx-auto mb-1`}
        style={{
          width: 42,
          height: 42,
          border: '2px solid white',
          boxShadow: '0 2px 4px rgba(0,0,0,0.2)',
        }}
      >
        <div className="fw-bold" style={{ fontSize: '1.1rem' }}>
          {voteAverage.toFixed(1)}
        </div>
      </div>
      <div>
        <small className="text-muted" style={{ fontSize: '0.8rem' }}>
          {voteCount.toLocaleString()} votes
        </small>
      </div>
    </div>
  );
};

const PopularitySection = ({ popularity }: { popularity: number }) => {
  const getPopularityPercentage = (value: number) => {
    const percentage =
      ((value - POPULARITY_CONFIG.MIN) / (POPULARITY_CONFIG.MAX - POPULARITY_CONFIG.MIN)) * 100;
    return Math.min(Math.max(percentage, 0), 100);
  };

  return (
    <div>
      <small className="text-muted d-block mb-1" style={{ fontSize: '0.8rem' }}>
        Popularity
      </small>
      <div className="progress" style={{ height: `${POPULARITY_CONFIG.BAR_HEIGHT}px` }}>
        <div
          className="progress-bar bg-info"
          role="progressbar"
          style={{ width: `${getPopularityPercentage(popularity)}%` }}
          aria-valuenow={getPopularityPercentage(popularity)}
          aria-valuemin={0}
          aria-valuemax={100}
        />
      </div>
    </div>
  );
};

const RatingAndPopularitySection = ({
  voteAverage,
  voteCount,
  popularity,
}: {
  voteAverage: number;
  voteCount: number;
  popularity: number;
}) => {
  return (
    <div className="text-center" style={{ width: '90px' }}>
      <RatingSection voteAverage={voteAverage} voteCount={voteCount} />
      <PopularitySection popularity={popularity} />
    </div>
  );
};

const TvShowDetails = observer(() => {
  const { id } = useParams<{ id: string }>();
  const numberId = id ? parseInt(id, 10) : null;
  const { show, loading, error } = tvShowDetailsStore;

  useEffect(() => {
    if (numberId) {
      tvShowDetailsStore.fetchTVShowDetails(numberId.toString());
    }
  }, [numberId]);

  if (!numberId || isNaN(numberId)) {
    return <Navigate to={ROUTES.NOT_FOUND} />;
  }

  if (loading) {
    return (
      <Container className="mt-4 text-center">
        <Spinner animation="border" role="status">
          <span className="visually-hidden">Loading...</span>
        </Spinner>
      </Container>
    );
  }

  if (error) {
    return (
      <Container className="mt-4">
        <Alert variant="danger">{error}</Alert>
      </Container>
    );
  }

  if (!show) {
    return null;
  }

  return (
    <Container className="mt-4">
      <Card className="mb-4">
        <Row className="g-0">
          <Col md={3}>
            {show.poster?.w342 && (
              <img
                src={show.poster.w342}
                alt={show.name}
                className="img-fluid h-100 object-fit-cover"
                style={{ maxHeight: '500px' }}
              />
            )}
          </Col>
          <Col md={9}>
            <Card.Body>
              <div className="d-flex justify-content-between align-items-start mb-4">
                <div>
                  <div className="d-flex align-items-center mb-2">
                    <h2 className="mb-0 me-3">{show.name}</h2>
                    {show.origin_country && show.origin_country.length > 0 && (
                      <div className="d-flex align-items-center">
                        {show.origin_country.map((country) => (
                          <span
                            key={country}
                            className="me-2 p-1"
                            style={{
                              fontSize: '1.5rem',
                              filter: 'drop-shadow(1px 1px 1px rgba(0,0,0,0.3))',
                            }}
                          >
                            {getCountryFlag(country)}
                          </span>
                        ))}
                      </div>
                    )}
                  </div>
                  {show.original_name && show.original_name !== show.name && (
                    <h5 className="text-muted mb-2">{show.original_name}</h5>
                  )}
                  {show.first_air_date && (
                    <div className="text-muted">{formatDate(show.first_air_date)}</div>
                  )}
                </div>
                <div>
                  {show.vote_average !== undefined && show.popularity !== undefined && (
                    <RatingAndPopularitySection
                      voteAverage={show.vote_average}
                      voteCount={show.vote_count || 0}
                      popularity={show.popularity}
                    />
                  )}
                </div>
              </div>

              {/* Genres */}
              {show.genres && show.genres.length > 0 && (
                <div className="mb-4">
                  {show.genres.map((genre) => (
                    <Badge bg="secondary" className="me-2 mb-2" key={genre}>
                      {genre}
                    </Badge>
                  ))}
                </div>
              )}

              {/* Overview */}
              {show.overview && (
                <div className="mb-3">
                  <p className="text-secondary mb-0">{show.overview}</p>
                </div>
              )}
            </Card.Body>
          </Col>
        </Row>
      </Card>

      {/* Seasons Grid */}
      {show.seasons && (
        <Row xs={1} md={2} lg={4} className="g-4">
          {show.seasons.map((season) => (
            <Col key={season.id}>
              <Card className="h-100">
                {season.poster?.w342 && (
                  <Card.Img
                    variant="top"
                    src={season.poster.w342}
                    alt={season.name}
                    style={{ height: '300px', objectFit: 'cover' }}
                  />
                )}
                <Card.Body>
                  <Card.Title>{season.name}</Card.Title>
                  <Card.Text className="text-muted">{season.episode_count} episodes</Card.Text>
                  {season.air_date && (
                    <Card.Text className="text-muted">
                      <small>{formatDate(season.air_date)}</small>
                    </Card.Text>
                  )}
                </Card.Body>
              </Card>
            </Col>
          ))}
        </Row>
      )}
    </Container>
  );
});

export default TvShowDetails;
