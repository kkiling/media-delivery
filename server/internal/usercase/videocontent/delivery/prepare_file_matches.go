package delivery

import (
	"context"
	"fmt"
	"path/filepath"
	"sort"
	"strings"

	"github.com/samber/lo"

	"github.com/kkiling/media-delivery/internal/adapter/matchtvshow"
)

type PreparingFileMatchesParams struct {
	TorrentFiles []FileInfo
	Episodes     []EpisodeInfo
	SeasonInfo   SeasonInfo
}

func toTrack(tracks []matchtvshow.Track, fileMap map[string]FileInfo) []Track {
	return lo.Map(tracks, func(item matchtvshow.Track, _ int) Track {
		return Track{
			Name:     item.Name,
			Language: item.Language,
			File:     fileMap[item.File],
		}
	})
}

func episodeToString(seasonNumber uint8, episodeNumber int) string {
	return fmt.Sprintf("%d-%d", seasonNumber, episodeNumber)
}

func mapResult(prepare []matchtvshow.Episode, torrentFiles []FileInfo, episodes []EpisodeInfo) ([]ContentMatches, error) {
	epMap := make(map[string]EpisodeInfo)
	for _, episode := range episodes {
		epMap[episodeToString(episode.SeasonNumber, episode.EpisodeNumber)] = episode
	}
	fileMap := make(map[string]FileInfo)
	for _, file := range torrentFiles {
		fileMap[file.RelativePath] = file
	}

	result := make([]ContentMatches, 0, len(prepare))
	for _, p := range prepare {
		episode, ok := epMap[episodeToString(p.SeasonNumber, p.EpisodeNumber)]
		if !ok {
			continue
		}
		ext := strings.ToLower(filepath.Ext(p.VideoFile))
		episode.FileName += ext
		episode.RelativePath += ext

		content := ContentMatches{
			Episode: episode,
			Video: VideoFile{
				File: fileMap[p.VideoFile],
			},
			AudioFiles: toTrack(p.AudioFiles, fileMap),
			Subtitles:  toTrack(p.Subtitles, fileMap),
		}

		result = append(result, content)
	}

	return result, nil
}

// PrepareFileMatches получение информации о файлах раздачи
func (s *Service) PrepareFileMatches(ctx context.Context, params PreparingFileMatchesParams) ([]ContentMatches, error) {
	// Подготавливаем параметры для преобразования файлов
	torrentFiles := lo.Map(params.TorrentFiles, func(item FileInfo, _ int) string {
		return item.RelativePath
	})

	prepareResult, err := s.prepareTVShow.MatchEpisodeFiles(torrentFiles)
	if err != nil {
		return nil, fmt.Errorf("prepareTVShow.MatchEpisodeFiles: %w", err)
	}

	// Получаем инфу о метчах
	result, err := mapResult(prepareResult, params.TorrentFiles, params.Episodes)
	if err != nil {
		return nil, fmt.Errorf("mapResult: %w", err)
	}

	sort.Slice(result, func(i, j int) bool {
		if result[i].Episode.SeasonNumber == result[j].Episode.SeasonNumber {
			return result[i].Episode.EpisodeNumber < result[j].Episode.EpisodeNumber
		}
		return result[i].Episode.SeasonNumber < result[j].Episode.SeasonNumber
	})

	return result, nil
}
