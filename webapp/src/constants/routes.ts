export const ROUTES = {
  HOME: '/',
  SEARCH: {
    ROOT: '/search',
    // Helper function to generate search URL with query
    createUrl: (query: string) => {
      const params = new URLSearchParams({ query });
      return `/search?${params.toString()}`;
    },
  },
  LIBRARY: {
    ROOT: '/library',
    MOVIES: {
      ROOT: '/library/movies',
      DETAILS: '/library/movies/:id',
      // Helper function to generate movie details URL
      getDetails: (id: string | number) => `/library/movies/${id}`,
    },
    TV_SHOWS: {
      ROOT: '/library/tvshows',
      DETAILS: '/library/tvshows/:id',
      SEASON: '/library/tvshows/:id/:season',
      // Helper functions to generate TV show URLs
      getDetails: (id: string | number) => `/library/tvshows/${id}`,
      getSeason: (id: string | number, season: string | number) =>
        `/library/tvshows/${id}/${season}`,
    },
  },
  NOT_FOUND: '*',
} as const;
