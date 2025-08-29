import { Form, Button, InputGroup } from 'react-bootstrap';
import { Search as SearchIcon } from 'react-bootstrap-icons';
import { useState, FormEvent, ChangeEvent } from 'react';

export enum SearchInputSize {
  Small = 'sm',
  Medium = 'md',
  Large = 'lg',
}

const SEARCH_VALIDATION = {
  MIN_LENGTH: 3,
  MAX_LENGTH: 256,
} as const;

interface SearchInputProps {
  size: SearchInputSize;
  placeholder?: string;
  initialQuery?: string;
  onSubmit: (query: string) => void;
}

export function SearchInput({
  size,
  placeholder = 'Search...',
  initialQuery = '',
  onSubmit,
}: SearchInputProps) {
  const [searchQuery, setSearchQuery] = useState(initialQuery);
  const [isInvalid, setIsInvalid] = useState(false);

  const handleSubmit = (e: FormEvent) => {
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
    onSubmit(trimmed);
  };

  const handleChange = (e: ChangeEvent<HTMLInputElement>) => {
    setSearchQuery(e.target.value);
    if (isInvalid) setIsInvalid(false);
  };

  return (
    <Form onSubmit={handleSubmit}>
      <InputGroup size={size === SearchInputSize.Medium ? undefined : size}>
        <Form.Control
          type="search"
          placeholder={placeholder}
          value={searchQuery}
          onChange={handleChange}
          className="focus-ring-0 shadow-none"
          isInvalid={isInvalid}
        />
        <Button type="submit" variant="primary" className="px-4">
          <SearchIcon size={size === SearchInputSize.Large ? 20 : 16} className="me-2" />
          <span className="d-none d-sm-inline">Search</span>
        </Button>
        <Form.Control.Feedback type="invalid">
          Please enter between {SEARCH_VALIDATION.MIN_LENGTH} and {SEARCH_VALIDATION.MAX_LENGTH}{' '}
          characters
        </Form.Control.Feedback>
      </InputGroup>
    </Form>
  );
}
