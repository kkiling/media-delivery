import { makeAutoObservable } from 'mobx';
import { Api, TVShow } from '@/api/api';

export class TVShowDetailsStore {
  show: TVShow | null = null;
  loading = false;
  error: string | null = null;

  constructor() {
    makeAutoObservable(this);
  }

  setShow(show: TVShow | null) {
    this.show = show;
  }

  setLoading(loading: boolean) {
    this.loading = loading;
  }

  setError(error: string | null) {
    this.error = error;
  }

  reset() {
    this.show = null;
    this.loading = false;
    this.error = null;
  }

  async fetchTVShowDetails(tvShowId: string) {
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

      const response = await api.v1.tvShowLibraryServiceGetTvShowInfo(tvShowId);
      this.setShow(response.data.result || null);
    } catch (error) {
      console.error('Error:', error);
      this.setError('Failed to fetch TV show details. Please try again.');
    } finally {
      this.setLoading(false);
    }
  }
}

export const tvShowDetailsStore = new TVShowDetailsStore();
