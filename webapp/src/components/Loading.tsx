import { Container, Spinner } from 'react-bootstrap';

interface LoadingProps {
  text?: string;
}

export default function Loading({ text }: LoadingProps) {
  return (
    <Container className="mt-4 mb-4 text-center">
      <div className="d-flex align-items-center justify-content-center">
        <Spinner animation="border" role="status" />
        {text && <span className="ms-2">{text}</span>}
      </div>
    </Container>
  );
}
