import { useState } from 'react';
import { Form, Button, Row, Col } from 'react-bootstrap';
import { Search } from 'react-bootstrap-icons';
import { useNavigate } from 'react-router-dom';

export default function TrendingMovies() {
  const [searchQuery, setSearchQuery] = useState('');
  const navigate = useNavigate();

  const handleSearch = (e: React.FormEvent) => {
    e.preventDefault();
    if (searchQuery.trim()) {
      navigate(`/search?query=${encodeURIComponent(searchQuery)}`);
    }
  };

  return (
    <>
      <h1 className="text-center mb-3">Welcome to Media Delivery</h1>
      <p className="text-center text-muted mb-4">
        Your gateway to thousands of movies and TV shows
      </p>
      <Form onSubmit={handleSearch}>
        <Form.Group controlId="searchForm">
          <Row className="g-2">
            <Col md={10}>
              <Form.Control
                type="search"
                placeholder="Search movies and TV shows..."
                value={searchQuery}
                onChange={(e) => setSearchQuery(e.target.value)}
              />
            </Col>
            <Col md={2}>
              <Button
                variant="primary"
                type="submit"
                className="w-100 d-flex align-items-center justify-content-center"
              >
                <Search className="me-2" />
                <span>Search</span>
              </Button>
            </Col>
          </Row>
        </Form.Group>
      </Form>
    </>
  );
}
