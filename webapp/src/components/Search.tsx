import { Form, Row, Col, ButtonGroup, ToggleButton, Button, InputGroup } from 'react-bootstrap';
import { useNavigate } from 'react-router-dom';
import { useState, FormEvent } from 'react';
import { Search as SearchIcon } from 'react-bootstrap-icons';
import { ROUTES } from '@/constants/routes';

export default function Search() {
  const navigate = useNavigate();
  const [searchQuery, setSearchQuery] = useState('');
  const [searchType, setSearchType] = useState('tvshows');

  const handleSearch = (e: FormEvent) => {
    e.preventDefault();
    if (searchQuery.trim()) {
      const url =
        searchType === 'tvshows'
          ? ROUTES.SEARCH.createTvShowsUrl(searchQuery)
          : ROUTES.SEARCH.createMoviesUrl(searchQuery);
      navigate(url);
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
          <Row>
            <Col xs={12}>
              <InputGroup size="lg" className="h-100">
                <InputGroup.Text
                  className="bg-white p-0"
                  style={{
                    minWidth: '160px', // Increased width for better text fit
                    height: '100%',
                    border: 0,
                  }}
                >
                  <ButtonGroup className="w-100 h-100">
                    <ToggleButton
                      id="toggle-tvshows"
                      type="radio"
                      variant="outline-secondary"
                      name="searchType"
                      value="tvshows"
                      checked={searchType === 'tvshows'}
                      onChange={(e) => setSearchType(e.currentTarget.value)}
                      className="w-50 d-flex align-items-center justify-content-center"
                      style={{
                        borderRadius: '0.375rem 0 0 0.375rem',
                        border: '1px solid #dee2e6',
                        borderRight: 0,
                      }}
                    >
                      TV Shows
                    </ToggleButton>
                    <ToggleButton
                      id="toggle-movies"
                      type="radio"
                      variant="outline-secondary"
                      name="searchType"
                      value="movies"
                      checked={searchType === 'movies'}
                      onChange={(e) => setSearchType(e.currentTarget.value)}
                      className="w-50 d-flex align-items-center justify-content-center"
                      style={{
                        borderRadius: 0,
                        border: '1px solid #dee2e6',
                        borderLeft: 0,
                      }}
                    >
                      Movies
                    </ToggleButton>
                  </ButtonGroup>
                </InputGroup.Text>
                <Form.Control
                  type="search"
                  placeholder={searchType === 'tvshows' ? 'Search TV shows...' : 'Search movies...'}
                  value={searchQuery}
                  onChange={(e) => setSearchQuery(e.target.value)}
                  className="rounded-0"
                  style={{ borderLeft: 0 }}
                />
                <Button type="submit" variant="primary" className="rounded-0 rounded-end px-3">
                  <SearchIcon size={20} />
                </Button>
              </InputGroup>
            </Col>
          </Row>
        </Form.Group>
      </Form>
    </>
  );
}
