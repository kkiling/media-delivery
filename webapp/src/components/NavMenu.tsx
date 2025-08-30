import { Navbar, Nav, Container } from 'react-bootstrap';
import { Link, useLocation } from 'react-router-dom';
import { Film, Tv } from 'react-bootstrap-icons';
import { ROUTES } from '@/constants/routes';

export default function NavMenu() {
  const location = useLocation();

  const isMoviesActive = location.pathname.startsWith(ROUTES.LIBRARY.MOVIES.ROOT);
  const isTvShowsActive = location.pathname.startsWith(ROUTES.LIBRARY.TV_SHOWS.ROOT);

  return (
    <Navbar bg="dark" variant="dark" expand="lg" className="mb-3">
      <Container>
        <Navbar.Brand as={Link} to={ROUTES.HOME} className="d-flex align-items-center">
          <Film className="me-2" />
          <span>Media Delivery</span>
        </Navbar.Brand>
        <Navbar.Toggle aria-controls="basic-navbar-nav" />
        <Navbar.Collapse id="basic-navbar-nav">
          <Nav className="ms-auto">
            <Nav.Link
              as={Link}
              to={ROUTES.LIBRARY.MOVIES.ROOT}
              className="d-flex align-items-center mx-2"
              active={isMoviesActive}
            >
              <Film className="me-2" />
              My Movies
            </Nav.Link>
            <Nav.Link
              as={Link}
              to={ROUTES.LIBRARY.TV_SHOWS.ROOT}
              className="d-flex align-items-center mx-2"
              active={isTvShowsActive}
            >
              <Tv className="me-2" />
              My TV Shows
            </Nav.Link>
          </Nav>
        </Navbar.Collapse>
      </Container>
    </Navbar>
  );
}
