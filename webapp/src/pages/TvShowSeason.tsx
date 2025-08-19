import { useParams, Link, Navigate } from 'react-router-dom';
import { useEffect, useState } from 'react';
import { observer } from 'mobx-react-lite';
import { Card, Container, Row, Col, Spinner, Alert, Table, Button } from 'react-bootstrap';
import { seasonDetailsStore } from '@/stores/seasonDetailsStore';
import { ROUTES } from '@/constants/routes';
import { formatDate } from '@/utils/formatting';
import { Rating } from '@/components';
import { Episode } from '@/api/api';
import { ArrowLeft, Image as ImageIcon } from 'react-bootstrap-icons';

const SEASON_INFO_CONFIG = {
  MIN_IMAGE_HEIGHT: 400,
  OVERVIEW_LINES: 8,
} as const;

const EPISODE_CONFIG = {
  IMAGE_HEIGHT: 120, // Height for episode still images
  OVERVIEW_LINES: 4,
} as const;

const TABLE_CONFIG = {
  COLUMNS: {
    NUMBER: '30px',
    EPISODE: 'auto',
    AIR_DATE: '160px',
    DURATION: '80px',
    RATING: '80px',
  },
} as const;

interface NoImageFallbackProps {
  text?: string;
}

const NoImageFallback = ({ text = 'No Image Available' }: NoImageFallbackProps) => (
  <div className="d-flex flex-column align-items-center justify-content-center bg-secondary text-white w-100 h-100">
    <ImageIcon size={48} className="mb-2" />
    <p className="mb-0">{text}</p>
  </div>
);

interface PosterImageProps {
  src?: string;
  alt: string;
  minHeight?: number;
}

const PosterImage = ({ src, alt, minHeight }: PosterImageProps) => {
  const [error, setError] = useState(false);

  if (!src || error) {
    return <NoImageFallback />;
  }

  return (
    <img
      src={src}
      alt={alt}
      onError={() => setError(true)}
      className="w-100 h-100 object-fit-cover"
      style={{ minHeight: minHeight ? `${minHeight}px` : 'auto' }}
    />
  );
};

interface EpisodeRowProps {
  episode: Episode;
}

const EpisodeRow = ({ episode }: EpisodeRowProps) => (
  <tr>
    <td>{episode.episode_number}</td>
    <td>
      <div className="d-flex align-items-center">
        {episode.still?.w342 && (
          <img
            src={episode.still.w342}
            alt={episode.name}
            className="me-3"
            style={{
              height: EPISODE_CONFIG.IMAGE_HEIGHT,
              objectFit: 'cover',
            }}
          />
        )}
        <div>
          <strong>{episode.name}</strong>
          <p
            className="text-muted mb-0 small"
            style={{
              overflow: 'hidden',
              textOverflow: 'ellipsis',
              display: '-webkit-box',
              WebkitLineClamp: EPISODE_CONFIG.OVERVIEW_LINES,
              WebkitBoxOrient: 'vertical',
            }}
          >
            {episode.overview}
          </p>
        </div>
      </div>
    </td>
    <td>{formatDate(episode.air_date)}</td>
    <td>{episode.runtime} min</td>
    <td>
      <Rating voteAverage={episode.vote_average ?? 0} voteCount={episode.vote_count ?? 0} />
    </td>
  </tr>
);

const TvShowSeason = observer(() => {
  const { id, season } = useParams<{ id: string; season: string }>();
  const numberId = id ? parseInt(id, 10) : null;
  const numberSeason = season ? parseInt(season, 10) : null;
  const { season: seasonData, episodes, loading, error } = seasonDetailsStore;

  useEffect(() => {
    if (numberId && numberSeason !== null) {
      seasonDetailsStore.fetchSeasonDetails(numberId.toString(), numberSeason);
    }
  }, [numberId, numberSeason]);

  if (!numberId || isNaN(numberId) || numberSeason === null || isNaN(numberSeason)) {
    return <Navigate to={ROUTES.NOT_FOUND} />;
  }

  if (loading) {
    return (
      <Container className="mt-4 text-center">
        <Spinner animation="border" role="status" />
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

  if (!seasonData) {
    return null;
  }

  return (
    <Container className="mt-4">
      <div className="mb-4">
        <Link to={ROUTES.LIBRARY.TV_SHOWS.getDetails(numberId)}>
          <Button variant="outline-primary" className="d-inline-flex align-items-center">
            <ArrowLeft className="me-2" /> Back to TV Show
          </Button>
        </Link>
      </div>

      <Card className="mb-4">
        <Row className="g-0">
          <Col md={3}>
            <div style={{ height: SEASON_INFO_CONFIG.MIN_IMAGE_HEIGHT }}>
              <PosterImage
                src={seasonData.poster?.w342}
                alt={seasonData.name || 'Season Poster'}
                minHeight={SEASON_INFO_CONFIG.MIN_IMAGE_HEIGHT}
              />
            </div>
          </Col>
          <Col md={9}>
            <Card.Body>
              <div className="d-flex justify-content-between align-items-start mb-4">
                <div>
                  <Card.Title as="h2">{seasonData.name}</Card.Title>
                  <Card.Subtitle className="mb-3 text-muted">
                    {formatDate(seasonData.air_date)} • {seasonData.episode_count} episodes
                  </Card.Subtitle>
                </div>
                <div className="text-center" style={{ width: '90px' }}>
                  {seasonData.vote_average && (
                    <Rating voteAverage={seasonData.vote_average} voteCount={0} />
                  )}
                </div>
              </div>
              {seasonData.overview && (
                <Card.Text
                  className="text-secondary"
                  style={{
                    overflow: 'hidden',
                    textOverflow: 'ellipsis',
                    display: '-webkit-box',
                    WebkitLineClamp: SEASON_INFO_CONFIG.OVERVIEW_LINES,
                    WebkitBoxOrient: 'vertical',
                  }}
                >
                  {seasonData.overview}
                </Card.Text>
              )}
            </Card.Body>
          </Col>
        </Row>
      </Card>

      <Card>
        <Card.Header as="h5">Episodes</Card.Header>
        <Card.Body>
          <Table responsive hover>
            <colgroup>
              <col style={{ width: TABLE_CONFIG.COLUMNS.NUMBER }} />
              <col style={{ width: TABLE_CONFIG.COLUMNS.EPISODE }} />
              <col style={{ width: TABLE_CONFIG.COLUMNS.AIR_DATE }} />
              <col style={{ width: TABLE_CONFIG.COLUMNS.DURATION }} />
              <col style={{ width: TABLE_CONFIG.COLUMNS.RATING }} />
            </colgroup>
            <tbody>
              {episodes.map((episode) => (
                <EpisodeRow key={episode.id} episode={episode} />
              ))}
            </tbody>
          </Table>
        </Card.Body>
      </Card>
    </Container>
  );
});

export default TvShowSeason;
