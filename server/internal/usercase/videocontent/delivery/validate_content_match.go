package delivery

import (
	"fmt"

	"github.com/samber/lo"
)

// ValidateContentMatch валидация
func (s *Service) ValidateContentMatch(oldContentMatch, newContentMatch *ContentMatches) error {
	if err := validateEpisodesStructure(oldContentMatch, newContentMatch); err != nil {
		return err
	}

	if err := validateTracksTypes(newContentMatch); err != nil {
		return err
	}

	if !hasVideoWithAudio(newContentMatch) {
		return fmt.Errorf("no video with audio tracks")
	}

	if err := validateTracksConsistency(oldContentMatch, newContentMatch); err != nil {
		return err
	}

	if err := validateDefaultTracks(newContentMatch); err != nil {
		return err
	}

	return nil
}

// validateEpisodesStructure проверяет структуру эпизодов
func validateEpisodesStructure(oldContentMatch, newContentMatch *ContentMatches) error {
	if len(oldContentMatch.Matches) != len(newContentMatch.Matches) {
		return fmt.Errorf("old and new contents matches do not match")
	}

	for index, oldMatch := range oldContentMatch.Matches {
		newMatch := newContentMatch.Matches[index]
		if oldMatch.Episode.SeasonNumber != newMatch.Episode.SeasonNumber ||
			oldMatch.Episode.EpisodeNumber != newMatch.Episode.EpisodeNumber ||
			oldMatch.Episode.FullPath != newMatch.Episode.FullPath ||
			oldMatch.Episode.RelativePath != newMatch.Episode.RelativePath {
			return fmt.Errorf("episode mismatch")
		}
	}

	return nil
}

// validateTracksTypes проверяет типы треков
func validateTracksTypes(contentMatch *ContentMatches) error {
	for _, match := range contentMatch.Matches {
		if match.Video == nil && (len(match.AudioTracks) > 0 || len(match.Subtitles) > 0) {
			return fmt.Errorf("episode has audio or subtitle tracks but no video assigned")
		}

		if match.Video != nil && match.Video.Type != TrackTypeVideo {
			return fmt.Errorf("track type mismatch")
		}

		if err := validateTrackSliceType(match.AudioTracks, TrackTypeAudio); err != nil {
			return err
		}

		if err := validateTrackSliceType(match.Subtitles, TrackTypeSubtitle); err != nil {
			return err
		}
	}

	return nil
}

// validateTrackSliceType проверяет тип треков в слайсе
func validateTrackSliceType(tracks []Track, expectedType TrackType) error {
	for _, track := range tracks {
		if track.Type != expectedType {
			return fmt.Errorf("track type mismatch")
		}
	}
	return nil
}

// hasVideoWithAudio проверяет наличие видео с аудио треками
func hasVideoWithAudio(contentMatch *ContentMatches) bool {
	for _, match := range contentMatch.Matches {
		if match.Video != nil && len(match.AudioTracks) > 0 {
			return true
		}
	}
	return false
}

// validateTracksConsistency проверяет консистентность треков
func validateTracksConsistency(oldContentMatch, newContentMatch *ContentMatches) error {
	oldTrackMap := createOldTrackMap(oldContentMatch)
	newTracks := collectAllTracks(newContentMatch)

	for _, newTrack := range newTracks {
		oldTrack, exists := oldTrackMap[newTrack.File.RelativePath]
		if !exists {
			return fmt.Errorf("episode has unallocated file")
		}

		if err := compareTracks(oldTrack, newTrack); err != nil {
			return err
		}
	}

	return nil
}

// createOldTrackMap создает мапу старых треков
func createOldTrackMap(contentMatch *ContentMatches) map[string]Track {
	trackMap := make(map[string]Track)

	for _, match := range contentMatch.Matches {
		if match.Video != nil {
			trackMap[match.Video.File.RelativePath] = *match.Video
		}
		for _, audioTrack := range match.AudioTracks {
			trackMap[audioTrack.File.RelativePath] = audioTrack
		}
		for _, subtitleTrack := range match.Subtitles {
			trackMap[subtitleTrack.File.RelativePath] = subtitleTrack
		}
	}

	for _, track := range contentMatch.Unallocated {
		trackMap[track.File.RelativePath] = track
	}

	return trackMap
}

// collectAllTracks собирает все треки из контент матча
func collectAllTracks(contentMatch *ContentMatches) []Track {
	tracks := make([]Track, 0)

	for _, match := range contentMatch.Matches {
		if match.Video != nil {
			tracks = append(tracks, *match.Video)
		}
		tracks = append(tracks, match.AudioTracks...)
		tracks = append(tracks, match.Subtitles...)
	}

	tracks = append(tracks, contentMatch.Unallocated...)

	return tracks
}

// compareTracks сравнивает два трека
func compareTracks(oldTrack, newTrack Track) error {
	if lo.FromPtrOr(oldTrack.Name, "") != lo.FromPtrOr(newTrack.Name, "") {
		return fmt.Errorf("track name mismatch")
	}
	if oldTrack.Type != newTrack.Type {
		return fmt.Errorf("track mismatch")
	}
	if oldTrack.File.RelativePath != newTrack.File.RelativePath {
		return fmt.Errorf("track mismatch")
	}
	if oldTrack.File.FullPath != newTrack.File.FullPath {
		return fmt.Errorf("track mismatch")
	}

	return nil
}

// validateDefaultTracks проверяет существование дефолтных треков
func validateDefaultTracks(contentMatch *ContentMatches) error {
	allTracks := collectAllTracks(contentMatch)

	if contentMatch.Options.DefaultAudioTrackName != nil {
		defaultAudioName := *contentMatch.Options.DefaultAudioTrackName
		if !containsTrackByNameAndType(allTracks, defaultAudioName, TrackTypeAudio) {
			return fmt.Errorf("default audio does not exist")
		}
	}

	if contentMatch.Options.DefaultSubtitleTrack != nil {
		defaultSubtitleName := *contentMatch.Options.DefaultSubtitleTrack
		if !containsTrackByNameAndType(allTracks, defaultSubtitleName, TrackTypeSubtitle) {
			return fmt.Errorf("default subtitle does not exist")
		}
	}

	return nil
}

// containsTrackByNameAndType проверяет наличие трека по имени и типу
func containsTrackByNameAndType(tracks []Track, name string, trackType TrackType) bool {
	return lo.ContainsBy(tracks, func(track Track) bool {
		return track.Type == trackType && lo.FromPtr(track.Name) == name
	})
}
