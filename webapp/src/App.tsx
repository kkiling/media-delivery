import { Routes, Route } from 'react-router-dom';
import { Layout, NotFound } from './components';
import {
  Home,
  Library,
  MovieDetails,
  Movies,
  TvShows,
  TvShowSeason,
  TvShowDetails,
  Search,
} from './pages';
import { ROUTES } from './constants/routes';

function App() {
  return (
    <Routes>
      <Route path={ROUTES.HOME} element={<Layout />}>
        <Route index element={<Home />} />
        <Route path={ROUTES.SEARCH.ROOT} element={<Search />} />
        <Route path={ROUTES.LIBRARY.ROOT} element={<Library />}>
          <Route path={ROUTES.LIBRARY.MOVIES.ROOT} element={<Movies />} />
          <Route path={ROUTES.LIBRARY.MOVIES.DETAILS} element={<MovieDetails />} />
          <Route path={ROUTES.LIBRARY.TV_SHOWS.ROOT} element={<TvShows />} />
          <Route path={ROUTES.LIBRARY.TV_SHOWS.DETAILS} element={<TvShowDetails />} />
          <Route path={ROUTES.LIBRARY.TV_SHOWS.SEASON} element={<TvShowSeason />} />
        </Route>
        <Route path={ROUTES.NOT_FOUND} element={<NotFound />} />
      </Route>
    </Routes>
  );
}

export default App;
