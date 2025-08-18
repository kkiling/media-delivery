import { Form, Row, Col, ButtonGroup, ToggleButton, Button, InputGroup } from 'react-bootstrap';
import { useNavigate } from 'react-router-dom';
import { useState, FormEvent } from 'react';
import { Search as SearchIcon, Film, Tv } from 'react-bootstrap-icons';
import { ROUTES } from '@/constants/routes';

const SEARCH_VALIDATION = {
  MIN_LENGTH: 3,
  MAX_LENGTH: 256,
} as const;

const searchOptions = [
  { value: 'tvshows', label: 'TV Shows', icon: Tv },
  { value: 'movies', label: 'Movies', icon: Film },
] as const;

type SearchType = (typeof searchOptions)[number]['value'];
type MainSearchMode = 'tvshows' | 'movies' | 'all';

interface MainSearchProps {
  mode?: MainSearchMode; // 'tvshows', 'movies', 'all'
  query?: string; // добавляем параметр query
}

export default function MainSearch({ mode = 'all', query = '' }: MainSearchProps) {
  const navigate = useNavigate();
  const [searchQuery, setSearchQuery] = useState(query);
  const [isInvalid, setIsInvalid] = useState(false);

  // Определяем searchType и searchOptions в зависимости от режима
  const isToggleVisible = mode === 'all';
  const defaultType: SearchType =
    mode === 'tvshows' ? 'tvshows' : mode === 'movies' ? 'movies' : 'tvshows';
  const [searchType, setSearchType] = useState<SearchType>(defaultType);

  const currentOptions =
    mode === 'all' ? searchOptions : searchOptions.filter((opt) => opt.value === defaultType);

  const handleSearch = (e: FormEvent) => {
    e.preventDefault();
    const trimmed = searchQuery.trim();
    if (
      trimmed.length < SEARCH_VALIDATION.MIN_LENGTH ||
      trimmed.length > SEARCH_VALIDATION.MAX_LENGTH
    ) {
      setIsInvalid(true);
      return;
    }
    setIsInvalid(false);

    const url =
      searchType === 'tvshows'
        ? ROUTES.SEARCH.createTvShowsUrl(trimmed)
        : ROUTES.SEARCH.createMoviesUrl(trimmed);

    navigate(url);
  };

  const placeholder = `Search ${searchType === 'tvshows' ? 'TV shows' : 'movies'}...`;

  return (
    <>
      <Form onSubmit={handleSearch}>
        <Form.Group controlId="searchForm">
          <Row>
            {/* Переключатель */}
            {isToggleVisible && (
              <Col xs={12} className="order-first mb-3">
                <div className="d-flex justify-content-center">
                  <ButtonGroup className="w-100" style={{ maxWidth: '320px' }}>
                    {currentOptions.map(({ value, label, icon: Icon }) => (
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
            )}

            {/* Строка поиска */}
            <Col xs={12}>
              <InputGroup size="lg">
                <Form.Control
                  type="search"
                  placeholder={placeholder}
                  value={searchQuery}
                  onChange={(e) => {
                    setSearchQuery(e.target.value);
                    if (isInvalid) setIsInvalid(false);
                  }}
                  className="focus-ring-0 shadow-none"
                  isInvalid={isInvalid}
                />
                <Button type="submit" variant="primary" className="px-4">
                  <SearchIcon size={20} className="me-2" />
                  <span className="d-none d-sm-inline">Search</span>
                </Button>
                <Form.Control.Feedback type="invalid">
                  Please enter between {SEARCH_VALIDATION.MIN_LENGTH} and{' '}
                  {SEARCH_VALIDATION.MAX_LENGTH} characters
                </Form.Control.Feedback>
              </InputGroup>
            </Col>
          </Row>
        </Form.Group>
      </Form>
    </>
  );
}
