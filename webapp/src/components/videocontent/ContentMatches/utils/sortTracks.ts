import { Track, TrackType } from '@/api/api';

const TRACK_TYPE_ORDER = {
  [TrackType.TRACK_TYPE_VIDEO]: 0,
  [TrackType.TRACK_TYPE_AUDIO]: 1,
  [TrackType.TRACK_TYPE_SUBTITLE]: 2,
};

export const sortTracks = (tracks: Track[]): Track[] => {
  return [...tracks].sort((a, b) => {
    const aType = a.type ?? '';
    const bType = b.type ?? '';

    if (aType !== bType) {
      return (
        (TRACK_TYPE_ORDER[aType as keyof typeof TRACK_TYPE_ORDER] ?? 3) -
        (TRACK_TYPE_ORDER[bType as keyof typeof TRACK_TYPE_ORDER] ?? 3)
      );
    }

    if (a.name === b.name) {
      return (a.relative_path || '').localeCompare(b.relative_path || '');
    }

    return (a.name || '').localeCompare(b.name || '');
  });
};
