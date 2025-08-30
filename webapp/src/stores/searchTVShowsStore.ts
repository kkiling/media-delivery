import { makeAutoObservable } from 'mobx';
import { Api, TVShowShort } from '@/api/api';

export class SearchTVShowsStore {
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

  async searchShows(query: string) {
    this.reset();

    if (!query.trim()) {
      return;
    }

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

      const response = await api.v1.tvShowLibraryServiceSearchTvShow({
        query: query,
      });

      this.setShows(response.data.items || []);
    } catch (error) {
      console.error('Error:', error);
      this.setError('Failed to fetch TV shows. Please try again.');
    } finally {
      this.setLoading(false);
    }
  }
}

export const searchTVShowsStore = new SearchTVShowsStore();
