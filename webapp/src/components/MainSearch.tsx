import { Form, Row, Col, ButtonGroup, ToggleButton } from 'react-bootstrap';
import { useNavigate } from 'react-router-dom';
import { useState } from 'react';
import { Film, Tv } from 'react-bootstrap-icons';
import { ROUTES } from '@/constants/routes';
import { SearchInput } from './SearchInput';

const searchOptions = [
  { value: 'tvshows', label: 'TV Shows', icon: Tv },
  { value: 'movies', label: 'Movies', icon: Film },
] as const;

type SearchType = (typeof searchOptions)[number]['value'];
type MainSearchMode = 'tvshows' | 'movies' | 'all';

interface MainSearchProps {
  mode?: MainSearchMode;
  query?: string;
  onSubmit?: (query: string) => void;
}

export default function MainSearch({ mode = 'all', query = '', onSubmit }: MainSearchProps) {
  const navigate = useNavigate();
  const isToggleVisible = mode === 'all';
  const defaultType: SearchType =
    mode === 'tvshows' ? 'tvshows' : mode === 'movies' ? 'movies' : 'tvshows';
  const [searchType, setSearchType] = useState<SearchType>(defaultType);

  const currentOptions =
    mode === 'all' ? searchOptions : searchOptions.filter((opt) => opt.value === defaultType);

  const handleSearch = (searchQuery: string) => {
    if (onSubmit) {
      onSubmit(searchQuery);
    } else {
      const url =
        searchType === 'tvshows'
          ? ROUTES.SEARCH.createTvShowsUrl(searchQuery)
          : ROUTES.SEARCH.createMoviesUrl(searchQuery);
      navigate(url);
    }
  };

  const placeholder = `Search ${searchType === 'tvshows' ? 'TV shows' : 'movies'}...`;

  return (
    <Form.Group controlId="searchForm">
      <Row>
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

        <Col xs={12}>
          <SearchInput placeholder={placeholder} initialQuery={query} onSubmit={handleSearch} />
        </Col>
      </Row>
    </Form.Group>
  );
}
