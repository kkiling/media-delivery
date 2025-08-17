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

function App() {
  return (
    <Routes>
      <Route path="/" element={<Layout />}>
        <Route index element={<Home />} />
        <Route path="search" element={<Search />} />
        <Route path="library" element={<Library />}>
          <Route path="movies" element={<Movies />} />
          <Route path="movies/:id" element={<MovieDetails />} />
          <Route path="tvshows" element={<TvShows />} />
          <Route path="tvshows/:id" element={<TvShowDetails />} />
          <Route path="tvshows/:id/:season" element={<TvShowSeason />} />
        </Route>
        <Route path="*" element={<NotFound />} />
      </Route>
    </Routes>
  );
}

export default App;
