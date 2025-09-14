import { useState, useCallback } from 'react';
import { ContentMatch } from '@/api/api';

export const useMatchesState = (initialMatches: ContentMatch[]) => {
  const [matches, setMatches] = useState<ContentMatch[]>(initialMatches);

  const updateMatch = useCallback(
    (index: number, updater: (match: ContentMatch) => ContentMatch) => {
      setMatches((prevMatches) => {
        const newMatches = prevMatches.slice();
        newMatches[index] = updater(prevMatches[index]);
        return newMatches;
      });
    },
    []
  );

  return {
    matches,
    setMatches: useCallback((newMatches: ContentMatch[]) => setMatches(newMatches), []),
    updateMatch,
  };
};
