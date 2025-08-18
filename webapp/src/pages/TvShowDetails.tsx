import { useParams, Navigate } from 'react-router-dom';
import { useEffect } from 'react';
import { observer } from 'mobx-react-lite';
import { Card, Container, Row, Col, Badge, Spinner, Alert } from 'react-bootstrap';
import { tvShowDetailsStore } from '@/stores/tvShowDetailsStore';
import { ROUTES } from '@/constants/routes';
import { formatDate } from '@/utils/formatting';
import getCountryFlag from 'country-flag-icons/unicode';
import { hasFlag } from 'country-flag-icons';
import { RatingSection } from '@/components/RatingSection';
import { PopularitySection } from '@/components/PopularitySection';
import { Season, TVShow } from '@/api/api';

const CountrySection = ({ countries }: { countries: string[] }) => {
  return (
    <div className="mb-4">
      <div className="d-flex flex-wrap gap-3">
        {countries.map((country) => {
          return (
            <div key={country} className="d-flex align-items-center bg-light rounded px-3 py-2">
              {hasFlag(country) && (
                <span className="me-2" style={{ fontSize: '1.2rem' }}>
                  {getCountryFlag(country)}
                </span>
              )}
              <span>
                {new Intl.DisplayNames(['en'], { type: 'region' }).of(country) || country}
              </span>
            </div>
          );
        })}
      </div>
    </div>
  );
};

const TVShowInfo = ({ show }: { show: TVShow }) => {
  return (
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
                <h2 className="mb-0 me-1">{show.name}</h2>
                {show.original_name && show.original_name !== show.name && (
                  <h5 className="text-muted mb-2">{show.original_name}</h5>
                )}
                {show.first_air_date && (
                  <div className="text-muted">{formatDate(show.first_air_date)}</div>
                )}
              </div>
              <div className="text-center" style={{ width: '90px' }}>
                {show.vote_average !== undefined && (
                  <RatingSection
                    voteAverage={show.vote_average}
                    voteCount={show.vote_count || 0}
                    showVoteCount={true}
                  />
                )}
                {show.popularity !== undefined && (
                  <PopularitySection popularity={show.popularity} />
                )}
              </div>
            </div>

            {show.origin_country && show.origin_country.length > 0 && (
              <CountrySection countries={show.origin_country} />
            )}

            {show.genres && show.genres.length > 0 && (
              <div className="mb-4">
                <h6 className="text-muted mb-2">Genres:</h6>
                {show.genres.map((genre) => (
                  <Badge bg="secondary" className="me-2 mb-2" key={genre}>
                    {genre}
                  </Badge>
                ))}
              </div>
            )}

            {show.overview && (
              <div className="mb-3">
                <p className="text-secondary mb-0">{show.overview}</p>
              </div>
            )}
          </Card.Body>
        </Col>
      </Row>
    </Card>
  );
};

const SeasonCard = ({ season }: { season: Season }) => {
  return (
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
      <TVShowInfo show={show} />

      {show.seasons && (
        <Row xs={1} md={2} lg={4} className="g-4">
          {show.seasons.map((season) => (
            <Col key={season.id}>
              <SeasonCard season={season} />
            </Col>
          ))}
        </Row>
      )}
    </Container>
  );
});

export default TvShowDetails;
