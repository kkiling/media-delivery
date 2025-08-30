import { Card } from 'react-bootstrap';
import { Film } from 'react-bootstrap-icons';

export default function TrendingMovies() {
  return (
    <Card className="h-100">
      <Card.Body className="d-flex flex-column">
        <Card.Title className="d-flex align-items-center">
          <Film className="me-2" />
          Trending Movies
        </Card.Title>
        <Card.Text className="flex-grow-1">Placeholder for trending movies content</Card.Text>
      </Card.Body>
    </Card>
  );
}
