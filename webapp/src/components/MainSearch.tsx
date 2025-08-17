import { Form, Row, Col, ButtonGroup, ToggleButton, Button, InputGroup } from 'react-bootstrap';
import { useNavigate } from 'react-router-dom';
import { useState, FormEvent } from 'react';
import { Search as SearchIcon, Film, Tv } from 'react-bootstrap-icons';
import { ROUTES } from '@/constants/routes';

const searchOptions = [
  { value: 'tvshows', label: 'TV Shows', icon: Tv },
  { value: 'movies', label: 'Movies', icon: Film },
] as const;

type SearchType = (typeof searchOptions)[number]['value'];

export default function MainSearch() {
  const navigate = useNavigate();
  const [searchQuery, setSearchQuery] = useState('');
  const [searchType, setSearchType] = useState<SearchType>('tvshows');

  const handleSearch = (e: FormEvent) => {
    e.preventDefault();
    if (!searchQuery.trim()) return;

    const url =
      searchType === 'tvshows'
        ? ROUTES.SEARCH.createTvShowsUrl(searchQuery)
        : ROUTES.SEARCH.createMoviesUrl(searchQuery);

    navigate(url);
  };

  const placeholder = `Search ${searchType === 'tvshows' ? 'TV shows' : 'movies'}...`;

  return (
    <>
      <h1 className="text-center mb-3">Welcome to Media Delivery</h1>
      <p className="text-center text-muted mb-4">
        Your gateway to thousands of movies and TV shows
      </p>

      <Form onSubmit={handleSearch}>
        <Form.Group controlId="searchForm">
          <Row>
            {/* Переключатель */}
            <Col xs={12} className="order-first mb-3">
              <div className="d-flex justify-content-center">
                <ButtonGroup className="w-100" style={{ maxWidth: '320px' }}>
                  {searchOptions.map(({ value, label, icon: Icon }) => (
                    <ToggleButton
                      key={value}
                      id={`toggle-${value}`}
                      type="radio"
                      variant="outline-secondary"
                      size="lg"
                      name="searchType"
                      value={value}
                      checked={searchType === value}
                      onChange={(e) => setSearchType(e.currentTarget.value as SearchType)}
                      className="d-flex align-items-center justify-content-center flex-grow-1"
                    >
                      <Icon className="me-2" /> {label}
                    </ToggleButton>
                  ))}
                </ButtonGroup>
              </div>
            </Col>

            {/* Строка поиска */}
            <Col xs={12}>
              <InputGroup size="lg">
                <Form.Control
                  type="search"
                  placeholder={placeholder}
                  value={searchQuery}
                  onChange={(e) => setSearchQuery(e.target.value)}
                  className="focus-ring-0 shadow-none"
                />
                <Button type="submit" variant="primary" className="px-4">
                  <SearchIcon size={20} className="me-2" />
                  <span className="d-none d-sm-inline">Search</span>
                </Button>
              </InputGroup>
            </Col>
          </Row>
        </Form.Group>
      </Form>
    </>
  );
}
