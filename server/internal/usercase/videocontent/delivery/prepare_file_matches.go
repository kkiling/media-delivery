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
}

func episodeToString(seasonNumber uint8, episodeNumber int) string {
	return fmt.Sprintf("%d-%d", seasonNumber, episodeNumber)
}

func mapResult(prepare []matchtvshow.Episode, torrentFiles []FileInfo, episodes []EpisodeInfo) (*ContentMatches, error) {
	epMap := make(map[string]EpisodeInfo)
	for _, episode := range episodes {
		epMap[episodeToString(episode.SeasonNumber, episode.EpisodeNumber)] = episode
	}
	torrentFilesMap := make(map[string]FileInfo)
	for _, file := range torrentFiles {
		torrentFilesMap[file.RelativePath] = file
	}

	toTrack := func(tracks []matchtvshow.Track, typeTrack TrackType) []Track {
		return lo.Map(tracks, func(item matchtvshow.Track, _ int) Track {
			return Track{
				Type:     typeTrack,
				Name:     &item.Name,
				Language: item.Language,
				File:     torrentFilesMap[item.File],
			}
		})
	}

	matches := make([]ContentMatch, 0, len(prepare))
	for _, p := range prepare {
		episode, ok := epMap[episodeToString(p.SeasonNumber, p.EpisodeNumber)]
		if !ok {
			continue
		}
		ext := strings.ToLower(filepath.Ext(p.VideoFile))
		episode.FileName += ext
		episode.RelativePath += ext

		content := ContentMatch{
			Episode: episode,
			Video: Track{
				Type: TrackTypeVideo,
				File: torrentFilesMap[p.VideoFile],
			},
			AudioFiles: toTrack(p.AudioFiles, TrackTypeAudio),
			Subtitles:  toTrack(p.Subtitles, TrackTypeSubtitle),
		}

		matches = append(matches, content)
	}

	// TODO: Разобраться с неопределенными файлами

	result := ContentMatches{
		Matches:     matches,
		Unallocated: []Track{},
		Options: ContentMatchesOptions{
			KeepOriginalAudio:     true,
			KeepOriginalSubtitles: true,
			DefaultAudioTrackName: func() *string {
				if len(matches) > 0 && len(matches[0].AudioFiles) > 0 {
					return matches[0].AudioFiles[0].Name
				}
				return nil
			}(),
			DefaultSubtitleTrack: nil,
		},
	}

	return &result, nil
}

// PrepareFileMatches получение информации о файлах раздачи
func (s *Service) PrepareFileMatches(ctx context.Context, params PreparingFileMatchesParams) (*ContentMatches, error) {
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

	// Сортируем эпизоды
	sort.Slice(result.Matches, func(i, j int) bool {
		if result.Matches[i].Episode.SeasonNumber == result.Matches[j].Episode.SeasonNumber {
			return result.Matches[i].Episode.EpisodeNumber < result.Matches[j].Episode.EpisodeNumber
		}
		return result.Matches[i].Episode.SeasonNumber < result.Matches[j].Episode.SeasonNumber
	})

	return result, nil
}
