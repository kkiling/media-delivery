import { makeAutoObservable } from 'mobx';
import { Api, TVShowShort } from '@/api/api';

export class LibraryTVShowsStore {
  shows: TVShowShort[] = [];
  loading = false;
  error: string | null = null;

  constructor() {
    makeAutoObservable(this);
  }

  setShows(shows: TVShowShort[]) {
    this.shows = shows;
  }

  setLoading(loading: boolean) {
    this.loading = loading;
  }

  setError(error: string | null) {
    this.error = error;
  }

  reset() {
    this.shows = [];
    this.loading = false;
    this.error = null;
  }

  async loadLibraryShows() {
    this.reset();
    this.setLoading(true);

    try {
      const api = new Api({
        baseApiParams: {
          headers: {
            'Content-Type': 'application/json',
            Accept: 'application/json',
          },
        },
      });

      const response = await api.v1.tvShowLibraryServiceGetTvShowsFromLibrary();
      this.setShows(response.data.items || []);
    } catch (error) {
      console.error('Error:', error);
      this.setError('Failed to fetch library TV shows. Please try again.');
    } finally {
      this.setLoading(false);
    }
  }
}

export const libraryTVShowsStore = new LibraryTVShowsStore();
