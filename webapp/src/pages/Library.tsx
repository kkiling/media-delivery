import { Outlet } from 'react-router-dom';

export default function Library() {
  return (
    <div>
      <h2>My Library</h2>
      <Outlet />
    </div>
  );
}
