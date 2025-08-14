package delivery

import (
	"context"
	"fmt"
	"path/filepath"

	ucerr "github.com/kkiling/torrent-to-media-server/internal/usercase/err"
	"github.com/kkiling/torrent-to-media-server/internal/usercase/tvshowlibrary"
	"github.com/kkiling/torrent-to-media-server/internal/usercase/videocontent/common"
	"github.com/samber/lo"
)

type GetEpisodesDataParams struct {
	TVShowID common.TVShowID
}

// GetEpisodesData получение информацию о эпизодах сериала и формируем имена каталогов и файлов
func (s *Service) GetEpisodesData(ctx context.Context, params GetEpisodesDataParams) (*EpisodesData, error) {
	// Получаем инфу о сезоне сериала
	tvShowInfo, err := s.tvShowLibrary.GetTVShowInfo(ctx, tvshowlibrary.GetTVShowParams{
		TVShowID: params.TVShowID.ID,
	})
	if err != nil {
		return nil, fmt.Errorf("tvShowLibrary.GetTVShowInfo: %w", err)
	}
	if tvShowInfo == nil {
		return nil, fmt.Errorf("tvShowInfo not found: %w", ucerr.NotFound)
	}

	season, find := lo.Find(tvShowInfo.Result.Seasons, func(item tvshowlibrary.Season) bool {
		return item.SeasonNumber == params.TVShowID.SeasonNumber
	})
	if !find {
		return nil, fmt.Errorf("season not found: %w", ucerr.NotFound)
	}

	// Достаем инфу о эпизодах
	episodes, err := s.tvShowLibrary.GetSeasonEpisodes(ctx, tvshowlibrary.GetSeasonEpisodesParams{
		TVShowID:     params.TVShowID.ID,
		SeasonNumber: params.TVShowID.SeasonNumber,
	})
	if err != nil {
		return nil, fmt.Errorf("tvShowLibrary.GetSeasonEpisodes: %w", err)
	}

	// Название сезона
	/*
		Series Name/
		  Season 01/
		    S01E01 - Episode Name.mp4
	*/
	tvShowName := fmt.Sprintf("%s (%d)" /*tvShowInfo.Result.Name*/, "Vinland", tvShowInfo.Result.FirstAirDate.Year())
	seasonName := fmt.Sprintf("S%02d %s", season.SeasonNumber, "Season 2") //season.Name)
	tvShowsPath := filepath.Join(s.config.BasePath, s.config.TVShowMediaSaveTvShowsPath, tvShowName)

	tvShowCatalogPath := TVShowCatalogPath{
		TVShowPath: tvShowsPath,
		SeasonPath: seasonName,
	}
	return &EpisodesData{
		TVShowCatalogPath: tvShowCatalogPath,
		Episodes: lo.Map(episodes.Items, func(item tvshowlibrary.Episode, _ int) EpisodeInfo {
			name := fmt.Sprintf("S%02dE%02d %s", season.SeasonNumber, item.EpisodeNumber, "Name") //item.Name)
			return EpisodeInfo{
				SeasonNumber:  season.SeasonNumber,
				EpisodeNumber: item.EpisodeNumber,
				EpisodeName:   item.Name,
				FileName:      filepath.Join(tvShowCatalogPath.FullSeasonPath(), name),
			}
		}),
	}, nil
}
