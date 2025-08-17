import { Outlet } from 'react-router-dom';
import { Container } from 'react-bootstrap';
import NavMenu from './NavMenu';

export default function Layout() {
  return (
    <>
      <NavMenu />
      <Container className="mt-3">
        <Outlet />
      </Container>
    </>
  );
}
