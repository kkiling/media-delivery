import { makeAutoObservable } from 'mobx';
import { Api, VideoContent, ContentID } from '@/api/api';

export class VideoContentStore {
  content: VideoContent | null = null;
  loading = false;
  error: string | null = null;

  constructor() {
    makeAutoObservable(this);
  }

  setContent = (content: VideoContent | null) => {
    this.content = content;
  };

  setLoading = (loading: boolean) => {
    this.loading = loading;
  };

  setError = (error: string | null) => {
    this.error = error;
  };

  reset = () => {
    this.content = null;
    this.loading = false;
    this.error = null;
  };

  fetchVideoContent = async (contentId: ContentID) => {
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

      const query: Record<string, unknown> = {};
      if (contentId?.movie_id) query['content_id.movie_id'] = contentId.movie_id;
      if (contentId?.tv_show?.id) query['content_id.tv_show.id'] = contentId.tv_show.id;
      if (contentId?.tv_show?.season_number)
        query['content_id.tv_show.season_number'] = contentId.tv_show.season_number;

      const response = await api.v1.videoContentServiceGetVideoContent(query);

      this.setContent(response.data.items?.[0] || null);
    } catch (error) {
      console.error('Error:', error);
      this.setError('Failed to fetch video content. Please try again.');
    } finally {
      this.setLoading(false);
    }
  };

  createVideoContent = async (contentId: ContentID) => {
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

      const response = await api.v1.videoContentServiceCreateVideoContent({
        content_id: contentId,
      });

      if (response.data.result) {
        this.setContent(response.data.result);
      }
    } catch (error) {
      console.error('Error:', error);
      this.setError('Failed to create video content. Please try again.');
    } finally {
      this.setLoading(false);
    }
  };
}

export const videoContentStore = new VideoContentStore();
