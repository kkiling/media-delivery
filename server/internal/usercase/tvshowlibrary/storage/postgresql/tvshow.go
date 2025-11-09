package postgresql

import (
	"context"
	"fmt"

	"github.com/kkiling/media-delivery/internal/usercase/tvshowlibrary"
	"github.com/kkiling/media-delivery/internal/usercase/tvshowlibrary/storage/db"
)

func (s *Storage) saveImage(ctx context.Context, img *tvshowlibrary.Image) error {
	queries := s.getQueries(ctx)

	err := queries.SaveImage(ctx, db.SaveImageParams{
		ID:       img.ID,
		W92:      img.W92,
		W185:     img.W185,
		W342:     img.W342,
		Original: img.Original,
	})

	if err != nil {
		return s.base.HandleError(err)
	}

	return nil
}

func (s *Storage) saveSeason(ctx context.Context, tvID uint64, season tvshowlibrary.Season) error {
	queries := s.getQueries(ctx)

	var posterID *string
	if season.Poster != nil {
		posterID = &season.Poster.ID
		if err := s.saveImage(ctx, season.Poster); err != nil {
			return s.base.HandleError(err)
		}
	}

	err := queries.SaveSeason(ctx, db.SaveSeasonParams{
		TvShowID:     int64(tvID),
		SeasonNumber: int(season.SeasonNumber),
		AirDate:      season.AirDate,
		EpisodeCount: int(season.EpisodeCount),
		Name:         season.Name,
		Overview:     season.Overview,
		PosterID:     posterID,
		VoteAverage:  season.VoteAverage,
	})

	if err != nil {
		return s.base.HandleError(err)
	}

	return nil
}

func (s *Storage) SaveTVShow(ctx context.Context, tvShow *tvshowlibrary.TVShow) error {
	return s.RunTransaction(ctx, func(tCtx context.Context) error {
		queries := s.getQueries(tCtx)

		var posterID *string
		var backdropID *string

		if tvShow.Poster != nil {
			posterID = &tvShow.Poster.ID
			if err := s.saveImage(tCtx, tvShow.Poster); err != nil {
				return fmt.Errorf("saveImage: %w", err)
			}
		}

		if tvShow.Backdrop != nil {
			backdropID = &tvShow.Backdrop.ID
			if err := s.saveImage(tCtx, tvShow.Backdrop); err != nil {
				return fmt.Errorf("saveImage: %w", err)
			}
		}

		err := queries.SaveTVShow(tCtx, db.SaveTVShowParams{
			ID:               int64(tvShow.ID),
			Name:             tvShow.Name,
			OriginalName:     tvShow.OriginalName,
			Overview:         tvShow.Overview,
			PosterID:         posterID,
			FirstAirDate:     tvShow.FirstAirDate,
			VoteAverage:      tvShow.VoteAverage,
			VoteCount:        int(tvShow.VoteCount),
			Popularity:       tvShow.Popularity,
			BackdropID:       backdropID,
			Genres:           tvShow.Genres,
			LastAirDate:      tvShow.LastAirDate,
			NumberOfEpisodes: int(tvShow.NumberOfEpisodes),
			NumberOfSeasons:  int(tvShow.NumberOfSeasons),
			OriginCountry:    tvShow.OriginCountry,
			Status:           tvShow.Status,
			Tagline:          tvShow.Tagline,
			Type:             tvShow.Type,
		})
		if err != nil {
			return s.base.HandleError(err)
		}

		// Сохраняем сезоны
		for _, season := range tvShow.Seasons {
			err = s.saveSeason(tCtx, tvShow.ID, season)
			if err != nil {
				return fmt.Errorf("saveSeason: %w", err)
			}
		}
		return nil
	})
}

func (s *Storage) SaveEpisodes(ctx context.Context, tvID uint64, seasonNumber uint8, episodes []tvshowlibrary.Episode) error {
	return s.RunTransaction(ctx, func(tCtx context.Context) error {
		queries := s.getQueries(tCtx)

		for _, episode := range episodes {
			var stillID *string
			if episode.Still != nil {
				stillID = &episode.Still.ID
				if err := s.saveImage(tCtx, episode.Still); err != nil {
					return fmt.Errorf("saveImage: %w", err)
				}
			}

			err := queries.SaveEpisode(tCtx, db.SaveEpisodeParams{
				TvShowID:      int64(tvID),
				SeasonNumber:  int(seasonNumber),
				AirDate:       episode.AirDate,
				EpisodeNumber: episode.EpisodeNumber,
				EpisodeType:   episode.EpisodeType,
				Name:          episode.Name,
				Overview:      episode.Overview,
				Runtime:       int(episode.Runtime),
				StillID:       stillID,
				VoteAverage:   episode.VoteAverage,
				VoteCount:     int(episode.VoteCount),
			})
			if err != nil {
				return s.base.HandleError(err)
			}
		}

		return nil
	})
}

func (s *Storage) getImage(ctx context.Context, imageID string) (*tvshowlibrary.Image, error) {
	queries := s.getQueries(ctx)

	res, err := queries.GetImage(ctx, imageID)
	if err != nil {
		return nil, s.base.HandleError(err)
	}

	return &tvshowlibrary.Image{
		ID:       res.ID,
		W92:      res.W92,
		W185:     res.W185,
		W342:     res.W342,
		Original: res.Original,
	}, nil
}

func (s *Storage) getSeasons(ctx context.Context, tvShowID uint64) ([]tvshowlibrary.Season, error) {
	queries := s.getQueries(ctx)

	res, err := queries.GetSeasons(ctx, int64(tvShowID))
	if err != nil {
		return nil, s.base.HandleError(err)
	}

	results := make([]tvshowlibrary.Season, 0, len(res))
	for _, item := range res {
		var poster *tvshowlibrary.Image
		var errPoster error
		if item.PosterID != nil {
			poster, errPoster = s.getImage(ctx, *item.PosterID)
			if errPoster != nil {
				return nil, s.base.HandleError(errPoster)
			}
		}

		results = append(results, tvshowlibrary.Season{
			AirDate:      item.AirDate,
			EpisodeCount: uint32(item.EpisodeCount),
			Name:         item.Name,
			Overview:     item.Overview,
			Poster:       poster,
			SeasonNumber: uint8(item.SeasonNumber),
			VoteAverage:  item.VoteAverage,
		})
	}

	return results, nil
}

func (s *Storage) GetTVShow(ctx context.Context, tvShowID uint64) (*tvshowlibrary.TVShow, error) {
	queries := s.getQueries(ctx)

	res, err := queries.GetTVShow(ctx, int64(tvShowID))
	if err != nil {
		return nil, s.base.HandleError(err)
	}

	var poster, backdrop *tvshowlibrary.Image
	if res.PosterID != nil {
		poster, err = s.getImage(ctx, *res.PosterID)
		if err != nil {
			return nil, fmt.Errorf("getImage: %w", err)
		}
	}
	if res.BackdropID != nil {
		backdrop, err = s.getImage(ctx, *res.BackdropID)
		if err != nil {
			return nil, fmt.Errorf("getImage: %w", err)
		}
	}

	seasons, err := s.getSeasons(ctx, tvShowID)
	if err != nil {
		return nil, fmt.Errorf("getSeasons: %w", err)
	}

	return &tvshowlibrary.TVShow{
		TVShowShort: tvshowlibrary.TVShowShort{
			ID:           uint64(res.ID),
			Name:         res.Name,
			OriginalName: res.OriginalName,
			Overview:     res.Overview,
			Poster:       poster,
			FirstAirDate: res.FirstAirDate,
			VoteAverage:  res.VoteAverage,
			VoteCount:    uint32(res.VoteCount),
			Popularity:   res.Popularity,
		},
		Backdrop:         backdrop,
		Genres:           res.Genres,
		LastAirDate:      res.LastAirDate,
		NumberOfEpisodes: uint32(res.NumberOfEpisodes),
		NumberOfSeasons:  uint32(res.NumberOfSeasons),
		OriginCountry:    res.OriginCountry,
		Status:           res.Status,
		Tagline:          res.Tagline,
		Type:             res.Type,
		Seasons:          seasons,
	}, nil
}

func (s *Storage) GetTVShows(ctx context.Context) ([]tvshowlibrary.TVShowShort, error) {
	queries := s.getQueries(ctx)

	res, err := queries.GetTVShows(ctx)
	if err != nil {
		return nil, s.base.HandleError(err)
	}

	results := make([]tvshowlibrary.TVShowShort, 0, len(res))
	for _, item := range res {
		var poster *tvshowlibrary.Image
		if item.PosterID != nil {
			poster, err = s.getImage(ctx, *item.PosterID)
			if err != nil {
				return nil, fmt.Errorf("getImage: %w", err)
			}
		}

		results = append(results, tvshowlibrary.TVShowShort{
			ID:           uint64(item.ID),
			Name:         item.Name,
			OriginalName: item.OriginalName,
			Overview:     item.Overview,
			Poster:       poster,
			FirstAirDate: item.FirstAirDate,
			VoteAverage:  item.VoteAverage,
			VoteCount:    uint32(item.VoteCount),
			Popularity:   item.Popularity,
		})
	}

	return results, nil
}

func (s *Storage) GetEpisodes(ctx context.Context, tvShowID uint64, seasonNumber uint8) ([]tvshowlibrary.Episode, error) {
	queries := s.getQueries(ctx)

	res, err := queries.GetEpisodes(ctx, db.GetEpisodesParams{
		TvShowID:     int64(tvShowID),
		SeasonNumber: int(seasonNumber),
	})
	if err != nil {
		return nil, s.base.HandleError(err)
	}

	results := make([]tvshowlibrary.Episode, 0, len(res))
	for _, item := range res {
		var still *tvshowlibrary.Image
		if item.StillID != nil {
			still, err = s.getImage(ctx, *item.StillID)
			if err != nil {
				return nil, fmt.Errorf("getImage: %w", err)
			}
		}

		results = append(results, tvshowlibrary.Episode{
			AirDate:       item.AirDate,
			EpisodeNumber: item.EpisodeNumber,
			EpisodeType:   item.EpisodeType,
			Name:          item.Name,
			Overview:      item.Overview,
			Runtime:       uint32(item.Runtime),
			Still:         still,
			VoteAverage:   item.VoteAverage,
			VoteCount:     uint32(item.VoteCount),
		})
	}

	return results, nil
}
