import { ContentMatches as ApiContentMatches } from '@/api/api';
import { ContentMatches as UiContentMatches } from './ContentMatches';

// Маппинг одного ContentMatches
function mapContentMatch(apiMatch: ApiContentMatches): UiContentMatches {
  return {
    episode: {
      episode_number: apiMatch.episode?.episode_number ?? 0,
      season_number: apiMatch.episode?.season_number ?? 0,
      episode_file: apiMatch.episode?.relative_path || '',
    },
    video: {
      file: apiMatch.video?.file?.relative_path || '',
      name: '',
      type: 'video',
    },
    audio_tracks: (apiMatch.audio_files || []).map((track) => ({
      name: track?.name || '',
      file: track?.file?.relative_path || '',
      type: 'audio',
    })),
    subtitle_tracks: (apiMatch.subtitles || []).map((track) => ({
      name: track?.name || '',
      file: track?.file?.relative_path || '',
      type: 'subtitle',
    })),
  };
}

// Основная функция для массива
export function mapContentMatches(apiMatches: ApiContentMatches[]): UiContentMatches[] {
  return apiMatches.map(mapContentMatch);
}
