package tvshowlibrary

import (
	"context"
	"errors"
	"fmt"

	"github.com/kkiling/media-delivery/internal/adapter/apierr"
	"github.com/kkiling/media-delivery/internal/adapter/themoviedb"
	ucerr "github.com/kkiling/media-delivery/internal/usercase/err"
)

const (
	language       = themoviedb.LanguageRU
	perPageDefault = 20
)

type Service struct {
	theMovieDb TheMovieDb
	storage    Storage
}

func NewService(
	storage Storage,
	theMovieDb TheMovieDb,
) *Service {
	return &Service{
		storage:    storage,
		theMovieDb: theMovieDb,
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

	tvShow := mapTVShow(response)

	return &GetTVShowResult{
		Result: tvShow,
	}, err
}

// GetSeasonInfo получение информации о сериях сезона
func (s *Service) GetSeasonInfo(ctx context.Context, params GetSeasonInfoParams) (*GetSeasonInfoResult, error) {
	response, err := s.theMovieDb.GetSeasonInfo(ctx, params.TVShowID, params.SeasonNumber, language)
	if err != nil {
		if errors.Is(err, apierr.ContentNotFound) {
			return nil, ucerr.NotFound
		}
		return nil, fmt.Errorf("theMovieDb.GetSeasonInfo: %w", err)
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

func (s *Service) AddTVShowInLibrary(ctx context.Context, params AddTVShowInLibraryParams) error {
	responseTVShow, err := s.theMovieDb.GetTV(ctx, params.TVShowID, language)
	if err != nil {
		if errors.Is(err, apierr.ContentNotFound) {
			return ucerr.NotFound
		}
		return fmt.Errorf("theMovieDb.GetTV: %w", err)
	}

	tvShow := mapTVShow(responseTVShow)

	responseSeason, err := s.theMovieDb.GetSeasonInfo(ctx, params.TVShowID, params.SeasonNumber, language)
	if err != nil {
		if errors.Is(err, apierr.ContentNotFound) {
			return ucerr.NotFound
		}
		return fmt.Errorf("theMovieDb.GetSeasonInfo: %w", err)
	}

	episodes := mapEpisodes(responseSeason.Episodes)

	// TODO: транзакция
	if err = s.storage.SaveOrUpdateTVShow(ctx, tvShow); err != nil {
		return fmt.Errorf("s.storage.SaveTVShow: %w", err)
	}
	// Получение информации о эпизодах сериала, автоматически добавляет их в библиотеку
	if err = s.storage.SaveOrUpdateSeasonEpisode(ctx, params.TVShowID, params.SeasonNumber, episodes); err != nil {
		return fmt.Errorf("s.storage.SaveOrUpdateSeasonEpisode: %w", err)
	}

	return nil
}
