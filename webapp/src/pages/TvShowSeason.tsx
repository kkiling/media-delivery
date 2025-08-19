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

// Обновляем конфигурацию для эпизодов
const EPISODE_CONFIG = {
  STILL_HEIGHT: 152,
  OVERVIEW_LINES: 3,
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

// Обновляем компонент EpisodeRow на EpisodeCard
interface EpisodeCardProps {
  episode: Episode;
}

const EpisodeCard = ({ episode }: EpisodeCardProps) => (
  <>
    <div className="px-2">
      <div className="d-flex flex-column flex-sm-row gap-3">
        {episode.still?.w342 && (
          <div className="d-flex justify-content-center">
            <div
              style={{
                width: '200px',
                height: EPISODE_CONFIG.STILL_HEIGHT,
              }}
            >
              <img
                src={episode.still.w342}
                alt={episode.name}
                className="w-100 h-100 object-fit-cover rounded"
              />
            </div>
          </div>
        )}
        <div className="flex-grow-1">
          <div className="d-flex justify-content-between align-items-start mb-2">
            <div>
              <h5 className="mb-1">
                <span className="text-muted me-2">#{episode.episode_number}</span>
                {episode.name}
              </h5>
              <div className="text-muted small mb-2">
                {formatDate(episode.air_date)} • {episode.runtime} min
              </div>
            </div>
            <div>
              <Rating
                voteAverage={episode.vote_average ?? 0}
                voteCount={episode.vote_count ?? 0}
                showVoteCount={true}
              />
            </div>
          </div>
          {episode.overview && (
            <p
              className="text-secondary mb-0 small"
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
          )}
        </div>
      </div>
    </div>
    <hr className="my-3" />
  </>
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
                    <Rating
                      voteAverage={seasonData.vote_average}
                      voteCount={0}
                      showVoteCount={false}
                    />
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

      {/* Обновляем секцию с эпизодами в основном компоненте */}
      <Card>
        <Card.Header as="h5">Episodes</Card.Header>
        <Card.Body>
          <div className="d-flex flex-column">
            {episodes.map((episode, index) => (
              <div key={episode.id}>
                <EpisodeCard episode={episode} />
                {/* Remove last divider */}
                {index === episodes.length - 1 && <hr className="d-none" />}
              </div>
            ))}
          </div>
        </Card.Body>
      </Card>
    </Container>
  );
});

export default TvShowSeason;
