import { Card } from 'react-bootstrap';
import { useNavigate } from 'react-router-dom';
import { ROUTES } from '@/constants/routes';
import { TVShowShort } from '@/api/api';
import { useState } from 'react';
import { Image as ImageIcon } from 'react-bootstrap-icons';

const CARD_CONFIG = {
  IMAGE_HEIGHT: 300,
} as const;

type TVShowCardProps = {
  show: TVShowShort;
};

function formatDate(dateString?: string) {
  if (!dateString) return '';
  return new Date(dateString).toLocaleDateString('en-US', {
    year: 'numeric',
    month: 'long',
    day: 'numeric',
  });
}

function getRatingColor(rating?: number) {
  if (!rating) return 'secondary';
  if (rating >= 7.8) return 'success';
  if (rating >= 5) return 'warning';
  return 'danger';
}

export function TVShowCard({ show }: TVShowCardProps) {
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
          <div
            className={`position-absolute top-0 start-0 m-2 bg-${getRatingColor(
              show.vote_average
            )} text-white rounded-circle d-flex align-items-center justify-content-center`}
            style={{
              width: 40,
              height: 40,
              border: '2px solid white',
            }}
          >
            <span className="fw-bold">{show.vote_average.toFixed(1)}</span>
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
