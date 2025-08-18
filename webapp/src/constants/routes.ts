export const ROUTES = {
  HOME: '/',
  SEARCH: {
    ROOT: '/search',
    TV_SHOWS: '/search/tvshows',
    MOVIES: '/search/movies',
    createTvShowsUrl: (query: string) => {
      const params = new URLSearchParams({ query });
      return `/search/tvshows?${params.toString()}`;
    },
    createMoviesUrl: (query: string) => {
      const params = new URLSearchParams({ query });
      return `/search/movies?${params.toString()}`;
    },
  },
  LIBRARY: {
    ROOT: '/library',
    MOVIES: {
      ROOT: '/library/movies',
      DETAILS: '/library/movies/:id',
      // Helper function to generate movie details URL
      getDetails: (id: number | undefined) => {
        if (!id) return '/not-found';
        return `/library/movies/${id}`;
      },
    },
    TV_SHOWS: {
      ROOT: '/library/tvshows',
      DETAILS: '/library/tvshows/:id',
      SEASON: '/library/tvshows/:id/:season',
      // Helper functions to generate TV show URLs
      getDetails: (id: number | undefined) => {
        if (!id) return '/not-found';
        return `/library/tvshows/${id}`;
      },
      getSeason: (id: number | undefined, season: number | undefined) => {
        if (!id || !season) return '/not-found';
        return `/library/tvshows/${id}/${season}`;
      },
    },
  },
  NOT_FOUND: '*',
} as const;
