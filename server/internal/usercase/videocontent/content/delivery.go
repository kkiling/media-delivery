package content

import (
	"context"
	"fmt"

	"github.com/kkiling/media-delivery/internal/common"
	ucerr "github.com/kkiling/media-delivery/internal/usercase/err"
	"github.com/kkiling/media-delivery/internal/usercase/videocontent/runners"
	"github.com/kkiling/media-delivery/internal/usercase/videocontent/runners/tvshowdeliverystate"
)

func (s *Service) validateDeliveryVideoContentParams(ctx context.Context, params DeliveryVideoContentParams) error {
	if err := params.ContentID.Validate(); err != nil {
		return err
	}
	// Фильмы пока не доступны
	if params.ContentID.MovieID != nil {
		return fmt.Errorf("movieID is not support: %w", ucerr.InvalidArgument)
	}

	return nil
}

// CreateDeliveryState создание файловой раздачи
func (s *Service) CreateDeliveryState(ctx context.Context, params DeliveryVideoContentParams) (*tvshowdeliverystate.State, error) {
	if err := s.validateDeliveryVideoContentParams(ctx, params); err != nil {
		return nil, fmt.Errorf("validateDeliveryVideoContentParams: %w", err)
	}
	// Достаем videoContent
	content, err := s.getVideoContent(ctx, params.ContentID)
	if err != nil {
		return nil, fmt.Errorf("getVideoContent: %w", err)
	}

	// Проверяем что он находится в правильном статусе
	switch content.DeliveryStatus {
	case DeliveryStatusNew:
	case DeliveryStatusDeleted:
	default:
		return nil, fmt.Errorf("video content is in invalid status: %w", ucerr.InvalidArgument)
	}

	options := tvshowdeliverystate.CreateOptions{
		TVShowID: *params.ContentID.TVShow,
		Index:    len(content.States),
	}
	var result *tvshowdeliverystate.State

	// TODO: одна транзакция
	{
		// Создали стейт доставки видео контента
		result, err = s.tvShowDeliveryState.Create(ctx, options)
		if err != nil {
			return nil, fmt.Errorf("tvShowDeliveryState.Create: %w", err)
		}

		// Обновить модель стейта videoContent
		updateVideoContent := UpdateVideoContent{
			DeliveryStatus: DeliveryStatusInProgress,
			States: append(content.States, State{
				StateID:   result.ID,
				CreatedAt: result.CreatedAt,
				Type:      runners.TVShowDelivery,
			}),
		}

		if err = s.storage.UpdateVideoContent(ctx, content.ID, &updateVideoContent); err != nil {
			return nil, fmt.Errorf("storage.UpdateVideoContent: %w", err)
		}
	}

	return result, nil
}

func (s *Service) GetDeliveryData(ctx context.Context, contentID common.ContentID) (*tvshowdeliverystate.State, error) {
	if err := contentID.Validate(); err != nil {
		return nil, err
	}

	content, err := s.getVideoContent(ctx, contentID)
	if err != nil {
		return nil, fmt.Errorf("getVideoContent: %w", err)
	}

	stateID := getLastState(content, runners.TVShowDelivery)
	if stateID == nil {
		return nil, fmt.Errorf("TVShowDelivery: %w", ucerr.NotFound)
	}

	result, err := s.tvShowDeliveryState.GetStateByID(ctx, *stateID)
	if err != nil {
		return nil, fmt.Errorf("s.GetStateByID: %w", err)
	}

	return result, nil
}

func (s *Service) ChoseTorrentOptions(ctx context.Context,
	contentID common.ContentID,
	opts tvshowdeliverystate.ChoseTorrentOptions,
) (*tvshowdeliverystate.State, error) {
	if err := contentID.Validate(); err != nil {
		return nil, err
	}

	content, err := s.getVideoContent(ctx, contentID)
	if err != nil {
		return nil, fmt.Errorf("getVideoContent: %w", err)
	}

	stateID := getLastState(content, runners.TVShowDelivery)
	if stateID == nil {
		return nil, fmt.Errorf("TVShowDelivery: %w", ucerr.NotFound)
	}

	newState, executeErr, err := s.tvShowDeliveryState.Complete(ctx, *stateID, opts)
	if err != nil {
		return nil, fmt.Errorf("tvShowDeliveryState.Complete: %w", err)
	}
	if executeErr != nil {
		s.logger.Errorf("tvShowDeliveryState.Complete: %v", executeErr)
	}
	return newState, nil
}

func (s *Service) ChoseFileMatchesOptions(ctx context.Context,
	contentID common.ContentID,
	opts tvshowdeliverystate.ChoseFileMatchesOptions,
) (*tvshowdeliverystate.State, error) {
	if err := contentID.Validate(); err != nil {
		return nil, err
	}

	content, err := s.getVideoContent(ctx, contentID)
	if err != nil {
		return nil, fmt.Errorf("getVideoContent: %w", err)
	}

	stateID := getLastState(content, runners.TVShowDelivery)
	if stateID == nil {
		return nil, fmt.Errorf("TVShowDelivery: %w", ucerr.NotFound)
	}

	newState, executeErr, err := s.tvShowDeliveryState.Complete(ctx, *stateID, opts)
	if err != nil {
		return nil, fmt.Errorf("tvShowDeliveryState.Complete: %w", err)
	}

	if executeErr != nil {
		s.logger.Errorf("tvShowDeliveryState.Complete: %v", executeErr)
	}
	return newState, nil
}
