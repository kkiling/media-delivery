import { Card } from 'react-bootstrap';
import { Tv } from 'react-bootstrap-icons';

export default function TrendingTvShows() {
  return (
    <Card className="h-100">
      <Card.Body className="d-flex flex-column">
        <Card.Title className="d-flex align-items-center">
          <Tv className="me-2" />
          Trending TV Shows
        </Card.Title>
        <Card.Text className="flex-grow-1">Placeholder for trending TV shows content</Card.Text>
      </Card.Body>
    </Card>
  );
}
