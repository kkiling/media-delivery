package content

import (
	"context"
	"fmt"

	"github.com/kkiling/media-delivery/internal/common"
	ucerr "github.com/kkiling/media-delivery/internal/usercase/err"
	"github.com/kkiling/media-delivery/internal/usercase/videocontent/runners"
	"github.com/kkiling/media-delivery/internal/usercase/videocontent/runners/tvshowdeletestate"
	"github.com/kkiling/media-delivery/internal/usercase/videocontent/tvshowdelete"
)

func (s *Service) validateDeleteVideoContentFilesParams(ctx context.Context, params DeleteVideoContentFilesParams) error {
	if err := params.ContentID.Validate(); err != nil {
		return err
	}
	// Фильмы пока не доступны
	if params.ContentID.MovieID != nil {
		return fmt.Errorf("movieID is not support: %w", ucerr.InvalidArgument)
	}

	return nil
}

func (s *Service) CreateDeleteState(ctx context.Context, params DeleteVideoContentFilesParams) (*tvshowdeletestate.State, error) {
	if err := s.validateDeleteVideoContentFilesParams(ctx, params); err != nil {
		return nil, fmt.Errorf("validateDeleteVideoContentFilesParams: %w", err)
	}

	content, err := s.getVideoContent(ctx, params.ContentID)
	if err != nil {
		return nil, fmt.Errorf("getVideoContent: %w", err)
	}

	switch content.DeliveryStatus {
	case DeliveryStatusDelivered:
	default:
		return nil, fmt.Errorf("video content is in invalid status: %w", ucerr.InvalidArgument)
	}

	stateID := getLastState(content, runners.TVShowDelivery)
	if stateID == nil {
		return nil, fmt.Errorf("TVShowDelivery: %w", ucerr.NotFound)
	}

	// Достаем инфу о стейте доставки, что бы вытащить от туда нужную инфу
	deliveryState, err := s.tvShowDeliveryState.GetStateByID(ctx, *stateID)
	if err != nil {
		return nil, fmt.Errorf("tvShowDeliveryState.GetStateByID: %w", err)
	}

	options := tvshowdeletestate.CreateOptions{
		TVShowID:    *params.ContentID.TVShow,
		MagnetHash:  deliveryState.Data.Torrent.MagnetLink.Hash,
		TorrentPath: deliveryState.Data.TorrentFilesData.ContentFullPath,
		TVShowCatalogPath: tvshowdelete.TVShowCatalogPath{
			TVShowPath: deliveryState.Data.EpisodesData.TVShowCatalogPath.TVShowPath,
			SeasonPath: deliveryState.Data.EpisodesData.TVShowCatalogPath.SeasonPath,
		},
	}

	var result *tvshowdeletestate.State
	//  TODO: одна транзакция
	{
		result, err = s.tvShowDeleteState.Create(ctx, options)
		if err != nil {
			return nil, fmt.Errorf("tvShowDeleteState.Create: %w", err)
		}

		// Обновить модель стейта videoContent
		updateVideoContent := UpdateVideoContent{
			DeliveryStatus: DeliveryStatusDeleting,
			States: append(content.States, State{
				StateID:   result.ID,
				CreatedAt: result.CreatedAt,
				Type:      runners.TVShowDelete,
			}),
		}

		if err = s.storage.UpdateVideoContent(ctx, content.ID, &updateVideoContent); err != nil {
			return nil, fmt.Errorf("storage.UpdateVideoContent: %w", err)
		}
	}

	return result, nil
}

func (s *Service) GetDeleteData(ctx context.Context, contentID common.ContentID) (*tvshowdeletestate.State, error) {
	if err := contentID.Validate(); err != nil {
		return nil, err
	}

	content, err := s.getVideoContent(ctx, contentID)
	if err != nil {
		return nil, fmt.Errorf("getVideoContent: %w", err)
	}

	stateID := getLastState(content, runners.TVShowDelete)
	if stateID == nil {
		return nil, fmt.Errorf("TVShowDelete: %w", ucerr.NotFound)
	}

	result, err := s.tvShowDeleteState.GetStateByID(ctx, *stateID)
	if err != nil {
		return nil, fmt.Errorf("s.GetStateByID: %w", err)
	}

	return result, nil
}
