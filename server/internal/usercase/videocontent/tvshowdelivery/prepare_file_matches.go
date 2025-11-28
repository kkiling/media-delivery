package tvshowdelivery

import (
	"context"
	"fmt"
	"sort"

	"github.com/samber/lo"

	"github.com/kkiling/media-delivery/internal/adapter/matchtvshow"
)

type PreparingFileMatchesParams struct {
	TorrentFiles []FileInfo
	Episodes     []EpisodeInfo
}

func episodeToString(seasonNumber uint8, episodeNumber int) string {
	return fmt.Sprintf("%d-%d", seasonNumber, episodeNumber)
}

func sortResult(result *ContentMatches) {

	sort.Slice(result.Matches, func(i, j int) bool {
		if result.Matches[i].Episode.SeasonNumber == result.Matches[j].Episode.SeasonNumber {
			return result.Matches[i].Episode.EpisodeNumber < result.Matches[j].Episode.EpisodeNumber
		}
		return result.Matches[i].Episode.SeasonNumber < result.Matches[j].Episode.SeasonNumber
	})

	// Сортируем нераспределенные файлы
	// Определяем порядок сортировки по типам
	var trackTypeOrder = map[TrackType]int{
		TrackTypeVideo:    0,
		TrackTypeAudio:    1,
		TrackTypeSubtitle: 2,
	}
	sort.Slice(result.Unallocated, func(i, j int) bool {
		a := result.Unallocated[i]
		b := result.Unallocated[j]

		// Сначала сортируем по типу
		aOrder := trackTypeOrder[a.Type]
		bOrder := trackTypeOrder[b.Type]

		if aOrder != bOrder {
			return aOrder < bOrder
		}

		// Если типы одинаковые, сортируем по имени
		if a.Name != nil && b.Name != nil && *a.Name != *b.Name {
			return *a.Name < *b.Name
		}

		// Если имена одинаковые или отсутствуют, сортируем по относительному пути
		return a.File.RelativePath < b.File.RelativePath
	})
}

func mapResult(prepare *matchtvshow.ContentMatches, torrentFiles []FileInfo, episodes []EpisodeInfo) (*ContentMatches, error) {
	matchesMap := make(map[string]*ContentMatch)
	unallocated := make([]Track, 0)
	matches := make([]ContentMatch, 0)

	for _, episode := range episodes {
		key := episodeToString(episode.SeasonNumber, episode.EpisodeNumber)
		matchesMap[key] = &ContentMatch{
			Episode: episode,
		}
	}

	torrentFilesMap := make(map[string]FileInfo)
	for _, file := range torrentFiles {
		torrentFilesMap[file.RelativePath] = file
	}

	toTrack := func(item matchtvshow.Track) Track {
		file := torrentFilesMap[item.File]
		return Track{
			Type:     TrackType(item.Type),
			Name:     &item.Name,
			Language: item.Language,
			File:     file,
		}
	}
	toTracks := func(tracks []matchtvshow.Track) []Track {
		return lo.Map(tracks, func(item matchtvshow.Track, _ int) Track {
			return toTrack(item)
		})
	}

	for _, prepareMatch := range prepare.Matches {
		key := episodeToString(prepareMatch.SeasonNumber, prepareMatch.EpisodeNumber)
		match, ok := matchesMap[key]
		if !ok {
			// Если не нашел то добавляем треки в нераспределенные
			if prepareMatch.Video != nil {
				unallocated = append(unallocated, toTrack(*prepareMatch.Video))
			}
			unallocated = append(unallocated, toTracks(prepareMatch.AudioTracks)...)
			unallocated = append(unallocated, toTracks(prepareMatch.Subtitles)...)
			continue
		}
		if prepareMatch.Video != nil {
			match.Video = lo.ToPtr(toTrack(*prepareMatch.Video))
		}
		match.AudioTracks = toTracks(prepareMatch.AudioTracks)
		match.Subtitles = toTracks(prepareMatch.Subtitles)
	}
	unallocated = append(unallocated, toTracks(prepare.Unallocated)...)

	// Сортируем эпизоды
	for _, track := range matchesMap {
		matches = append(matches, *track)
	}

	result := ContentMatches{
		Matches:     matches,
		Unallocated: unallocated,
		Options: ContentMatchesOptions{
			KeepOriginalAudio:     true,
			KeepOriginalSubtitles: true,
			DefaultAudioTrackName: func() *string {
				if len(matches) > 0 && len(matches[0].AudioTracks) > 0 {
					return matches[0].AudioTracks[0].Name
				}
				return nil
			}(),
			DefaultSubtitleTrack: nil,
		},
	}
	sortResult(&result)

	return &result, nil
}

// PrepareFileMatches получение информации о файлах раздачи
func (s *Service) PrepareFileMatches(ctx context.Context, params PreparingFileMatchesParams) (*ContentMatches, error) {
	// Подготавливаем параметры для преобразования файлов
	torrentFiles := lo.Map(params.TorrentFiles, func(item FileInfo, _ int) string {
		return item.RelativePath
	})

	prepare, err := s.prepareTVShow.MatchEpisodeFiles(torrentFiles)
	if err != nil {
		return nil, fmt.Errorf("prepareTVShow.MatchEpisodeFiles: %w", err)
	}

	// Получаем инфу о метчах
	result, err := mapResult(prepare, params.TorrentFiles, params.Episodes)
	if err != nil {
		return nil, fmt.Errorf("mapResult: %w", err)
	}

	return result, nil
}
