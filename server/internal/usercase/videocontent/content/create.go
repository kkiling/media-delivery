package content

import (
	"context"
	"fmt"

	"github.com/samber/lo"

	"github.com/kkiling/media-delivery/internal/common"
	ucerr "github.com/kkiling/media-delivery/internal/usercase/err"
	"github.com/kkiling/media-delivery/internal/usercase/labels"
	"github.com/kkiling/media-delivery/internal/usercase/tvshowlibrary"
)

func (s *Service) validateCreateVideoContentParams(ctx context.Context, params CreateVideoContentParams) error {
	if err := params.ContentID.Validate(); err != nil {
		return err
	}
	// Фильмы пока не доступны
	if params.ContentID.MovieID != nil {
		return fmt.Errorf("movieID is not support: %w", ucerr.InvalidArgument)
	}

	// Ограничиваем одну раздачу на один фильм/сериал
	found, err := s.GetVideoContent(ctx, params.ContentID)
	if err != nil {
		return fmt.Errorf("getVideoContent: %w", err)
	}
	if len(found) > 0 {
		return ucerr.AlreadyExists
	}

	return nil
}

func (s *Service) checkContentExistInLibrary(ctx context.Context, contentID common.ContentID) error {
	// Проверяем наличие сезона сериала изи фильма
	// Получаем инфу о сериале
	tvShowInfo, err := s.tvShowLibrary.GetTVShowInfo(ctx, tvshowlibrary.GetTVShowParams{
		TVShowID: contentID.TVShow.ID,
	})
	if err != nil {
		return fmt.Errorf("tvShowLibrary.GetTVShowInfo: %w", err)
	}
	if tvShowInfo == nil {
		return fmt.Errorf("tvShow: %w", ucerr.NotFound)
	}

	// Проверяем что сезон тоже существует
	if !lo.ContainsBy(tvShowInfo.Result.Seasons, func(item tvshowlibrary.Season) bool {
		return item.SeasonNumber == contentID.TVShow.SeasonNumber
	}) {
		return fmt.Errorf("season: %w", ucerr.NotFound)
	}

	return nil
}

// CreateVideoContent создание видеоконтента (без доставки)
func (s *Service) CreateVideoContent(ctx context.Context, params CreateVideoContentParams) (*VideoContent, error) {
	if err := s.validateCreateVideoContentParams(ctx, params); err != nil {
		return nil, fmt.Errorf("validateCreateVideoContentParams: %w", err)
	}

	if err := s.checkContentExistInLibrary(ctx, params.ContentID); err != nil {
		return nil, fmt.Errorf("checkContentExistInLibrary: %w", err)
	}

	now := s.clock.Now()

	videoContent := VideoContent{
		ID:             s.uuidGenerator.New(),
		CreatedAt:      now,
		ContentID:      params.ContentID,
		DeliveryStatus: DeliveryStatusNew,
	}

	tvShow := tvshowlibrary.AddTVShowInLibraryParams{
		TVShowID:     params.ContentID.TVShow.ID,
		SeasonNumber: params.ContentID.TVShow.SeasonNumber,
	}

	contentID := common.ContentID{
		TVShow: &common.TVShowID{
			ID:           params.ContentID.TVShow.ID,
			SeasonNumber: params.ContentID.TVShow.SeasonNumber,
		},
	}

	labelContentInLibrary := labels.Label{
		ContentID: contentID,
		TypeLabel: labels.ContentInLibrary,
		CreatedAt: now,
	}

	// TODO: !!! !!! !!! подумать как обернуть в одну транзакцию
	{
		// Добавили сериал в библиотеку
		if err := s.tvShowLibrary.AddTVShowInLibrary(ctx, tvShow); err != nil {
			return nil, fmt.Errorf("tvShowLibrary.AddTVShowInLibrary: %w", err)
		}
		// Добавили лейбл что сериал в библиотеке
		if err := s.labels.AddLabel(ctx, labelContentInLibrary); err != nil {
			return nil, fmt.Errorf("labels.AddLabel: %w", err)
		}
		// Создали сущность видео контента
		if err := s.storage.SaveVideoContent(ctx, &videoContent); err != nil {
			return nil, fmt.Errorf("storage.SaveVideoContent: %w", err)
		}
	}

	return &videoContent, nil
}
