package delivery

import (
	"context"
	"fmt"

	"github.com/samber/lo"

	"github.com/kkiling/media-delivery/internal/adapter/matchtvshow"
	"github.com/kkiling/media-delivery/internal/usercase/videocontent/common"
)

type PreparingFileMatchesParams struct {
	TorrentFiles []FileInfo
	Episodes     []EpisodeInfo
	TVShowID     common.TVShowID
}

func mapContentMatchesFromPrepareTVShowSeason(
	torrentFiles []FileInfo,
	episodes []EpisodeInfo,
	prepareResult *matchtvshow.PrepareTVShowSeason,
) ([]ContentMatches, error) {
	epMap := make(map[int]EpisodeInfo)
	for _, episode := range episodes {
		epMap[episode.EpisodeNumber] = episode
	}

	fileMap := make(map[string]FileInfo)
	for _, file := range torrentFiles {
		fileMap[file.RelativePath] = file
	}

	toTrack := func(tracks []matchtvshow.PrepareTrack) []Track {
		return lo.Map(tracks, func(item matchtvshow.PrepareTrack, _ int) Track {
			return Track{
				Name:     item.Name,
				Language: item.Language,
				File:     fileMap[item.File.RelativePath],
			}
		})
	}

	result := make([]ContentMatches, 0, len(prepareResult.Episodes))
	for _, prepareEpisode := range prepareResult.Episodes {
		episode := epMap[prepareEpisode.Episode.EpisodeNumber]

		var videoFile FileInfo
		if prepareEpisode.VideoFile != nil {
			videoFile = fileMap[prepareEpisode.VideoFile.File.RelativePath]
			episode.FileName += videoFile.Extension
			episode.RelativePath += videoFile.Extension
		}

		content := ContentMatches{
			Episode: episode,
			Video: VideoFile{
				File: videoFile,
			},
			AudioFiles: toTrack(prepareEpisode.AudioFiles),
			Subtitles:  toTrack(prepareEpisode.Subtitles),
		}

		result = append(result, content)
	}

	return result, nil
}

func mapToPrepareTvShowPrams(torrentFiles []FileInfo, episodes []EpisodeInfo) (*matchtvshow.PrepareTvShowPrams, error) {
	return &matchtvshow.PrepareTvShowPrams{
		Episodes: lo.Map(episodes, func(episode EpisodeInfo, _ int) matchtvshow.EpisodeInfo {
			return matchtvshow.EpisodeInfo{
				EpisodeNumber: episode.EpisodeNumber,
			}
		}),
		TorrentFiles: lo.Map(torrentFiles, func(item FileInfo, _ int) matchtvshow.TorrentFile {
			return matchtvshow.TorrentFile{
				RelativePath: item.RelativePath,
			}
		}),
	}, nil
}

// PrepareFileMatches получение информации о файлах раздачи
func (s *Service) PrepareFileMatches(ctx context.Context, params PreparingFileMatchesParams) ([]ContentMatches, error) {
	// Подготавливаем параметры для преобразования файлов
	prepareParams, err := mapToPrepareTvShowPrams(params.TorrentFiles, params.Episodes)

	if err != nil {
		return nil, fmt.Errorf("mapToPrepareTvShowPrams: %w", err)
	}

	prepareResult, err := s.prepareTVShow.PrepareTvShowSeason(prepareParams)
	if err != nil {
		return nil, fmt.Errorf("prepareTVShow.PrepareTvShowSeason: %w", err)
	}

	// Получаем инфу о метчах
	result, err := mapContentMatchesFromPrepareTVShowSeason(params.TorrentFiles, params.Episodes, prepareResult)
	if err != nil {
		return nil, fmt.Errorf("mapContentMatchesFromPrepareTVShowSeason: %w", err)
	}

	return result, nil
}
