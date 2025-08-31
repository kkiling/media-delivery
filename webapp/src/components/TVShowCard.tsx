import { Card } from 'react-bootstrap';
import { useNavigate } from 'react-router-dom';
import { TVShowShort } from '@/api/api';
import { PosterImage, Rating } from '@/components';
import { ROUTES } from '@/constants/routes';
import { formatDate } from '@/utils/formatting';

const CARD_CONFIG = {
  IMAGE_HEIGHT: 400,
} as const;

type TVShowCardProps = {
  show: TVShowShort;
};

export function TVShowCard({ show }: TVShowCardProps) {
  const navigate = useNavigate();

  return (
    <Card
      className="h-100 position-relative cursor-pointer"
      onClick={() => navigate(ROUTES.LIBRARY.TV_SHOWS.getDetails(show.id))}
    >
      <div className="position-relative" style={{ height: `${CARD_CONFIG.IMAGE_HEIGHT}px` }}>
        <PosterImage
          src={show.poster?.w185 || show.poster?.w342}
          alt={show.name}
          minHeight={CARD_CONFIG.IMAGE_HEIGHT}
        />

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
