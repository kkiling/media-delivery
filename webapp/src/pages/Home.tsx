import { Container, Row, Col } from 'react-bootstrap';
import { TrendingMovies, TrendingTvShows, MainSearch } from '@/components';

export default function Home() {
  return (
    <Container className="mt-4">
      {/* Hero Section with Search */}
      <Row className="mb-4" style={{ minHeight: '200px' }}>
        <Col>
          <h1 className="text-center mb-3">Welcome to Media Delivery</h1>
          <p className="text-center text-muted mb-4">
            Your gateway to thousands of movies and TV shows
          </p>
          <MainSearch />
        </Col>
      </Row>

      {/* Trending Sections - теперь всегда в разных строках */}
      <Row>
        <Col xs={12} className="mb-4">
          <TrendingMovies />
        </Col>
        <Col xs={12}>
          <TrendingTvShows />
        </Col>
      </Row>
    </Container>
  );
}
