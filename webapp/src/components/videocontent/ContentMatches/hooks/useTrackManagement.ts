import { useState } from 'react';
import { Track } from '@/api/api';
import { sortTracks } from '../utils/sortTracks';

export const useTrackManagement = (initialUnallocated: Track[]) => {
  const [unallocated, setUnallocated] = useState<Track[]>(initialUnallocated);

  const addToUnallocated = (track: Track) => {
    // Проверяем существование трека с таким же relative_path
    setUnallocated((prev) => {
      if (prev.some((t) => t.relative_path === track.relative_path)) {
        return prev; // Если трек уже существует, возвращаем предыдущий массив
      }
      return sortTracks([...prev, track]); // Если уникальный - добавляем и сортируем
    });
  };

  const removeFromUnallocated = (relative_path?: string) => {
    if (!relative_path) return;
    setUnallocated((prev) => prev.filter((t) => t.relative_path !== relative_path));
  };

  return { unallocated, setUnallocated, addToUnallocated, removeFromUnallocated };
};
