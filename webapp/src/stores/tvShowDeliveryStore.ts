import { makeAutoObservable } from 'mobx';
import { Api, TVShowDeliveryState, ContentID } from '@/api/api';

export class TVShowDeliveryStore {
  deliveryState: TVShowDeliveryState | null = null;
  loading = false;
  error: string | null = null;

  constructor() {
    makeAutoObservable(this);
  }

  setDeliveryState = (state: TVShowDeliveryState | null) => {
    this.deliveryState = state;
  };

  setLoading = (loading: boolean) => {
    this.loading = loading;
  };

  setError = (error: string | null) => {
    this.error = error;
  };

  reset = () => {
    this.deliveryState = null;
    this.loading = false;
    this.error = null;
  };

  fetchDeliveryData = async (contentId: ContentID) => {
    this.setLoading(true);
    this.setError(null);

    try {
      const api = new Api({
        baseApiParams: {
          headers: {
            'Content-Type': 'application/json',
            Accept: 'application/json',
          },
        },
      });

      const query: Record<string, any> = {};
      if (contentId?.movie_id) query['content_id.movie_id'] = contentId.movie_id;
      if (contentId?.tv_show?.id) query['content_id.tv_show.id'] = contentId.tv_show.id;
      if (contentId?.tv_show?.season_number)
        query['content_id.tv_show.season_number'] = contentId.tv_show.season_number;

      const response = await api.v1.videoContentServiceGetTvShowDeliveryData(query);
      this.setDeliveryState(response.data.result || null);
    } catch (error) {
      console.error('Error:', error);
      this.setError('Failed to fetch delivery data. Please try again.');
    } finally {
      this.setLoading(false);
    }
  };

  confirmFileMatches = async (contentId: ContentID) => {
    this.setLoading(true);
    this.setError(null);

    try {
      const api = new Api({
        baseApiParams: {
          headers: {
            'Content-Type': 'application/json',
            Accept: 'application/json',
          },
        },
      });

      const response = await api.v1.videoContentServiceChoseFileMatchesOptions({
        content_id: contentId,
        approve: true,
      });

      this.setDeliveryState(response.data.result || null);
    } catch (error) {
      console.error('Error:', error);
      this.setError('Failed to confirm file matches. Please try again.');
      throw error;
    } finally {
      this.setLoading(false);
    }
  };

  selectTorrent = async (contentId: ContentID, href?: string, newSearchQuery?: string) => {
    this.setLoading(true);
    this.setError(null);

    try {
      const api = new Api({
        baseApiParams: {
          headers: {
            'Content-Type': 'application/json',
            Accept: 'application/json',
          },
        },
      });

      const response = await api.v1.videoContentServiceChoseTorrentOptions({
        content_id: contentId,
        href,
        new_search_query: newSearchQuery,
      });

      this.setDeliveryState(response.data.result || null);
    } catch (error) {
      console.error('Error:', error);
      this.setError('Failed to select torrent. Please try again.');
      throw error;
    } finally {
      this.setLoading(false);
    }
  };
}

export const tvShowDeliveryStore = new TVShowDeliveryStore();
