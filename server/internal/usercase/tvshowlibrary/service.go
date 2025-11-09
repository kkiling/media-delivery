package tvshowlibrary

import (
	"context"
	"errors"
	"fmt"

	"github.com/kkiling/media-delivery/internal/adapter/apierr"
	"github.com/kkiling/media-delivery/internal/adapter/themoviedb"
	"github.com/kkiling/media-delivery/internal/common"
	ucerr "github.com/kkiling/media-delivery/internal/usercase/err"
)

const (
	language       = themoviedb.LanguageRU
	perPageDefault = 20
)

type Service struct {
	theMovieDb TheMovieDb
	storage    Storage
	clock      Clock
}

func NewService(
	storage Storage,
	theMovieDb TheMovieDb,
) *Service {
	return &Service{
		theMovieDb: theMovieDb,
		storage:    storage,
		clock:      &common.RealClock{},
	}
}

// SearchTVShow поиск сериалов по названию
func (s *Service) SearchTVShow(ctx context.Context, params TVShowSearchParams) (*TVShowSearchResult, error) {
	response, err := s.theMovieDb.SearchTV(ctx, themoviedb.SearchQuery{
		Language: language,
		Query:    params.Query,
		Page:     1,
		PerPage:  perPageDefault,
	})
	if err != nil {
		return nil, fmt.Errorf("theMovieDb.SearchTV: %w", err)
	}

	return &TVShowSearchResult{
		Items: mapTVShowShorts(response.Results),
	}, nil
}

// GetTVShowInfo получение подробной информации о сериале
func (s *Service) GetTVShowInfo(ctx context.Context, params GetTVShowParams) (*GetTVShowResult, error) {
	response, err := s.theMovieDb.GetTV(ctx, params.TVShowID, language)
	if err != nil {
		if errors.Is(err, apierr.ContentNotFound) {
			return nil, ucerr.NotFound
		}
		return nil, fmt.Errorf("theMovieDb.GetTV: %w", err)
	}

	return &GetTVShowResult{
		Result: mapTVShow(response),
	}, err
}

// GetSeasonInfo получение информации о сериях сезона
func (s *Service) GetSeasonInfo(ctx context.Context, params GetSeasonInfoParams) (*GetSeasonInfoResult, error) {
	response, err := s.theMovieDb.GetSeason(ctx, params.TVShowID, params.SeasonNumber, language)
	if err != nil {
		if errors.Is(err, apierr.ContentNotFound) {
			return nil, ucerr.NotFound
		}
		return nil, fmt.Errorf("theMovieDb.GetSeason: %w", err)
	}

	return &GetSeasonInfoResult{
		Result: &SeasonWithEpisodes{
			Season:   *mapSeason(response.Season),
			Episodes: mapEpisodes(response.Episodes),
		},
	}, nil
}

// GetTVShowsFromLibrary получение списка сериалов из библиотеки
func (s *Service) GetTVShowsFromLibrary(ctx context.Context, _ GetTVShowsFromLibraryParams) (*GetTVShowsFromLibraryResult, error) {
	tvShows, err := s.storage.GetTVShows(ctx)
	if err != nil {
		return nil, fmt.Errorf("storage.GetTVShow: %w", err)
	}
	return &GetTVShowsFromLibraryResult{
		Items: tvShows,
	}, nil
}

// AddTVShowInLibrary добавление сериала в библиотеку
func (s *Service) AddTVShowInLibrary(ctx context.Context, params AddTVShowInLibraryParams) error {
	// Проверка наличия сериала
	responseTVShow, err := s.theMovieDb.GetTV(ctx, params.TVShowID, language)
	if err != nil {
		if errors.Is(err, apierr.ContentNotFound) {
			return ucerr.NotFound
		}
		return fmt.Errorf("theMovieDb.GetTV: %w", err)
	}

	// Проверка наличия сезона
	responseSeason, err := s.theMovieDb.GetSeason(ctx, params.TVShowID, params.SeasonNumber, language)
	if err != nil {
		if errors.Is(err, apierr.ContentNotFound) {
			return ucerr.NotFound
		}
		return fmt.Errorf("theMovieDb.GetSeason: %w", err)
	}

	// TODO: транзакция
	{
		// Сохранение сериала
		tvShow := mapTVShow(responseTVShow)
		if err = s.storage.SaveTVShow(ctx, tvShow); err != nil {
			return fmt.Errorf("s.storage.SaveTVShow: %w", err)
		}
		// Сохранение серий сезона
		episodes := mapEpisodes(responseSeason.Episodes)
		if err = s.storage.SaveEpisodes(ctx, params.TVShowID, params.SeasonNumber, episodes); err != nil {
			return fmt.Errorf("s.storage.SaveEpisodes: %w", err)
		}
	}

	return nil
}
