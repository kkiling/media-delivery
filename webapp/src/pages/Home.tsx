import { Container, Row, Col } from 'react-bootstrap';
import { TrendingMovies, TrendingTvShows, Search } from '@/components';

export default function Home() {
  return (
    <Container className="mt-4">
      {/* Hero Section with Search */}
      <Row className="mb-4" style={{ minHeight: '200px' }}>
        <Col>
          <Search />
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
