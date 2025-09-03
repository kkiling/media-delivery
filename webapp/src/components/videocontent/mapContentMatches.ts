import { ContentMatches as ApiContentMatches, Track as ApiTrack } from '@/api/api';
import { ContentMatches as UiContentMatches } from './ContentMatches';

// Маппинг одного Track
function mapTrack(track?: ApiTrack): { name: string; file: string } {
  return {
    name: track?.name || '',
    file: track?.file?.relative_path || '',
  };
}

// Маппинг одного ContentMatches
function mapContentMatch(apiMatch: ApiContentMatches): UiContentMatches {
  return {
    episode: {
      episode_number: apiMatch.episode?.episode_number ?? 0,
      season_number: apiMatch.episode?.season_number ?? 0,
      episode_file: apiMatch.episode?.relative_path || '',
    },
    media_files: {
      video_file: apiMatch.video?.file?.relative_path || '',
      audio_files: (apiMatch.audio_files || []).map(mapTrack),
      subtitles_files: (apiMatch.subtitles || []).map(mapTrack),
    },
  };
}

// Основная функция для массива
export function mapContentMatches(apiMatches: ApiContentMatches[]): UiContentMatches[] {
  return apiMatches.map(mapContentMatch);
}
