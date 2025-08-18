import { Card } from 'react-bootstrap';
import { useNavigate } from 'react-router-dom';
import { ROUTES } from '@/constants/routes';
import { TVShowShort } from '@/api/api';
import { useState } from 'react';

const CARD_CONFIG = {
  IMAGE_HEIGHT: '300px',
} as const;

type TVShowCardProps = {
  show: TVShowShort;
};

export function TVShowCard({ show }: TVShowCardProps) {
  const navigate = useNavigate();
  const [imageError, setImageError] = useState(false);

  const formatDate = (dateString?: string) => {
    if (!dateString) return '';
    return new Date(dateString).toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'long',
      day: 'numeric',
    });
  };

  const getRatingColor = (rating?: number) => {
    if (!rating) return 'secondary'; // Серый, если рейтинга нет
    if (rating >= 7.8) return 'success';
    if (rating >= 5) return 'warning';
    return 'danger';
  };

  const renderRating = () => {
    if (show.vote_average === undefined) return null; // Не показываем, если рейтинга нет

    return (
      <div
        className={`position-absolute top-0 start-0 m-2 bg-${getRatingColor(
          show.vote_average
        )} text-white rounded-circle d-flex align-items-center justify-content-center`}
        style={{
          width: '40px',
          height: '40px',
          border: '2px solid white',
        }}
      >
        <span className="fw-bold">{show.vote_average.toFixed(1)}</span>
      </div>
    );
  };

  return (
    <Card
      className="h-100 cursor-pointer position-relative"
      onClick={() => navigate(ROUTES.LIBRARY.TV_SHOWS.getDetails(show.id))}
      style={{ cursor: 'pointer' }}
    >
      <div style={{ height: CARD_CONFIG.IMAGE_HEIGHT, position: 'relative' }}>
        {show.poster?.w342 && !imageError ? (
          <Card.Img
            variant="top"
            src={show.poster.w342}
            alt={show.name}
            onError={() => setImageError(true)}
            style={{ width: '100%', height: '100%', objectFit: 'cover' }}
          />
        ) : (
          <Card.Body className="d-flex align-items-center justify-content-center bg-secondary w-100 h-100">
            <div className="text-center text-white">
              <i className="bi bi-image" style={{ fontSize: '3rem' }}></i>
              <p className="mt-2 mb-0">No Image Available</p>
            </div>
          </Card.Body>
        )}
        {renderRating()}
      </div>

      <Card.Body className="d-flex flex-column" style={{ flexGrow: 1 }}>
        <Card.Title>{show.name || 'No title'}</Card.Title>
        <Card.Subtitle className="mb-2 text-muted">
          {formatDate(show.first_air_date) || 'Release date unknown'}
        </Card.Subtitle>
        <div style={{ flexGrow: 1 }}>
          <Card.Text
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
        </div>
      </Card.Body>
    </Card>
  );
}
