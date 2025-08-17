import { Card } from 'react-bootstrap';
import { useNavigate } from 'react-router-dom';
import { ROUTES } from '@/constants/routes';
import { TVShowShort } from '@/api/api';

type TVShowCardProps = {
  show: TVShowShort;
};

export function TVShowCard({ show }: TVShowCardProps) {
  const navigate = useNavigate();

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'long',
      day: 'numeric',
    });
  };

  return (
    <Card
      className="h-100 cursor-pointer"
      onClick={() => navigate(ROUTES.LIBRARY.TV_SHOWS.getDetails(show.id))}
      style={{ cursor: 'pointer' }}
    >
      <Card.Img
        variant="top"
        src={show.poster?.w342}
        alt={show.name}
        style={{ height: '400px', objectFit: 'cover' }}
      />
      <Card.Body>
        <Card.Title>{show.name}</Card.Title>
        <Card.Subtitle className="mb-2 text-muted">
          {show.first_air_date && formatDate(show.first_air_date)}
        </Card.Subtitle>
        <Card.Text
          style={{
            overflow: 'hidden',
            textOverflow: 'ellipsis',
            display: '-webkit-box',
            WebkitLineClamp: 3,
            WebkitBoxOrient: 'vertical',
          }}
        >
          {show.overview}
        </Card.Text>
      </Card.Body>
    </Card>
  );
}
