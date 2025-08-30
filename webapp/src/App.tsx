import { Routes, Route } from 'react-router-dom';
import { Layout, NotFound } from './components';
import {
  Home,
  Library,
  TvShowSeason,
  TvShowDetails,
  SearchTVShows,
  SearchMovies,
  TvShowsLibrary,
} from './pages';
import { ROUTES } from './constants/routes';

function App() {
  return (
    <Routes>
      <Route path={ROUTES.HOME} element={<Layout />}>
        <Route index element={<Home />} />
        <Route path={ROUTES.SEARCH.TV_SHOWS} element={<SearchTVShows />} />
        <Route path={ROUTES.SEARCH.MOVIES} element={<SearchMovies />} />
        <Route path={ROUTES.LIBRARY.ROOT} element={<Library />}>
          {/* <Route path={ROUTES.LIBRARY.MOVIES.ROOT} element={<Movies />} />
          <Route path={ROUTES.LIBRARY.MOVIES.DETAILS} element={<MovieDetails />} /> */}
          <Route path={ROUTES.LIBRARY.TV_SHOWS.ROOT} element={<TvShowsLibrary />} />
          <Route path={ROUTES.LIBRARY.TV_SHOWS.DETAILS} element={<TvShowDetails />} />
          <Route path={ROUTES.LIBRARY.TV_SHOWS.SEASON} element={<TvShowSeason />} />
        </Route>
        <Route path={ROUTES.NOT_FOUND} element={<NotFound />} />
      </Route>
    </Routes>
  );
}

export default App;
