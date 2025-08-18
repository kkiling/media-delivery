import { Card } from 'react-bootstrap';
import { useNavigate } from 'react-router-dom';
import { ROUTES } from '@/constants/routes';
import { TVShowShort } from '@/api/api';
import { useState } from 'react';
import { Image as ImageIcon } from 'react-bootstrap-icons';
import { formatDate } from '@/utils/formatting';
import { RatingSection } from './RatingSection';

const CARD_CONFIG = {
  IMAGE_HEIGHT: 500,
} as const;

type TVShowCardProps = {
  show: TVShowShort;
};

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
          <div className="position-absolute top-0 start-0 m-2">
            <RatingSection voteAverage={show.vote_average} voteCount={0} showVoteCount={false} />
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
