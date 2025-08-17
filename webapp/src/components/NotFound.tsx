import { Container, Row, Col, Button } from 'react-bootstrap';
import { useNavigate } from 'react-router-dom';
import { ROUTES } from '@/constants/routes';

export default function NotFound() {
  const navigate = useNavigate();

  return (
    <Container className="mt-5">
      <Row className="justify-content-center text-center">
        <Col md={6}>
          <h1 style={{ fontSize: '8rem', fontWeight: 'bold', color: '#dc3545' }}>404</h1>
          <h2 className="mb-4">Page Not Found</h2>
          <p className="lead mb-4">
            Oops! The page you&apos;re looking for seems to have vanished into thin air like a movie
            plot twist.
          </p>
          <Button
            variant="primary"
            size="lg"
            onClick={() => navigate(ROUTES.HOME)}
            className="px-4"
          >
            Back to Homepage
          </Button>
        </Col>
      </Row>
    </Container>
  );
}
