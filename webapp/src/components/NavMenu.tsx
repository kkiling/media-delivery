import { Navbar, Nav, Container } from 'react-bootstrap';
import { Link } from 'react-router-dom';
import { Film, Tv } from 'react-bootstrap-icons';

export default function NavMenu() {
  return (
    <Navbar bg="dark" variant="dark" expand="lg" className="mb-3">
      <Container>
        <Navbar.Brand as={Link} to="/" className="d-flex align-items-center">
          <Film className="me-2" />
          <span>Media Delivery</span>
        </Navbar.Brand>
        <Navbar.Toggle aria-controls="basic-navbar-nav" />
        <Navbar.Collapse id="basic-navbar-nav">
          <Nav className="ms-auto">
            <Nav.Link as={Link} to="/library/movies" className="d-flex align-items-center mx-2">
              <Film className="me-2" />
              My Movies
            </Nav.Link>
            <Nav.Link as={Link} to="/library/tvshows" className="d-flex align-items-center mx-2">
              <Tv className="me-2" />
              My TV Shows
            </Nav.Link>
          </Nav>
        </Navbar.Collapse>
      </Container>
    </Navbar>
  );
}
