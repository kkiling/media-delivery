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
	contents, err := s.storage.GetVideoContents(ctx, params.ContentID)
	if err != nil {
		return nil, fmt.Errorf("storage.GetVideoContent: %w", err)
	}
	if len(contents) != 1 {
		return nil, ucerr.NotFound
	}
	content := contents[0]
	// Проверяем что он находится в правильном статусе
	if content.DeliveryStatus != DeliveryStatusNew {
		return nil, fmt.Errorf("video content is not in new status: %w", ucerr.InvalidArgument)
	}

	var state *tvshowdeliverystate.State

	options := tvshowdeliverystate.CreateOptions{
		TVShowID: *params.ContentID.TVShow,
	}

	// TODO: одна транзакция
	{
		// Создали стейт доставки видео контента
		state, err = s.tvShowDeliveryState.Create(ctx, options)
		if err != nil {
			return nil, fmt.Errorf("tvShowDeliveryState.Create: %w", err)
		}

		// Обновить модель стейта videoContent
		updateVideoContent := UpdateVideoContent{
			DeliveryStatus: content.DeliveryStatus,
			States: append(content.States, State{
				StateID: state.ID,
				Type:    runners.TVShowDelivery,
			}),
		}

		if err = s.storage.UpdateVideoContent(ctx, content.ID, &updateVideoContent); err != nil {
			return nil, fmt.Errorf("storage.UpdateVideoContent: %w", err)
		}
	}

	return state, nil
}

func (s *Service) GetDeliveryData(ctx context.Context, contentID common.ContentID) (*tvshowdeliverystate.State, error) {
	if err := contentID.Validate(); err != nil {
		return nil, err
	}

	stateID, err := s.getStateID(ctx, contentID, runners.TVShowDelivery)
	if err != nil {
		return nil, fmt.Errorf("getStateID: %w", err)
	}

	result, err := s.tvShowDeliveryState.GetStateByID(ctx, stateID)
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
	stateID, err := s.getStateID(ctx, contentID, runners.TVShowDelivery)
	if err != nil {
		return nil, fmt.Errorf("getStateID: %w", err)
	}
	newState, executeErr, err := s.tvShowDeliveryState.Complete(ctx, stateID, opts)
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

	stateID, err := s.getStateID(ctx, contentID, runners.TVShowDelivery)
	if err != nil {
		return nil, fmt.Errorf("getStateID: %w", err)
	}
	newState, executeErr, err := s.tvShowDeliveryState.Complete(ctx, stateID, opts)
	if err != nil {
		return nil, fmt.Errorf("tvShowDeliveryState.Complete: %w", err)
	}
	if executeErr != nil {
		s.logger.Errorf("tvShowDeliveryState.Complete: %v", executeErr)
	}
	return newState, nil
}
