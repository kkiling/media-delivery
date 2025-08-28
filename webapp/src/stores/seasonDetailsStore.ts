import { makeAutoObservable } from 'mobx';
import { Api, Season, Episode } from '@/api/api';

export class SeasonDetailsStore {
  season: Season | null = null;
  episodes: Episode[] = [];
  loading = false;
  error: string | null = null;

  constructor() {
    makeAutoObservable(this);
  }

  setSeason(season: Season | null) {
    this.season = season;
  }

  setEpisodes(episodes: Episode[]) {
    this.episodes = episodes;
  }

  setLoading(loading: boolean) {
    this.loading = loading;
  }

  setError(error: string | null) {
    this.error = error;
  }

  reset() {
    this.season = null;
    this.episodes = [];
    this.loading = false;
    this.error = null;
  }

  async fetchSeasonDetails(tvShowId: number, seasonNumber: number) {
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

      const response = await api.v1.tvShowLibraryServiceGetSeasonInfo(
        tvShowId.toString(),
        seasonNumber
      );

      this.setSeason(response.data.season || null);
      this.setEpisodes(response.data.episodes || []);
    } catch (error) {
      console.error('Error:', error);
      this.setError('Failed to fetch season details. Please try again.');
    } finally {
      this.setLoading(false);
    }
  }
}

export const seasonDetailsStore = new SeasonDetailsStore();
