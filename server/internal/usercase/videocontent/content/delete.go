package content

import (
	"context"
	"fmt"

	"github.com/kkiling/media-delivery/internal/common"
	ucerr "github.com/kkiling/media-delivery/internal/usercase/err"
	"github.com/kkiling/media-delivery/internal/usercase/videocontent/runners"
	"github.com/kkiling/media-delivery/internal/usercase/videocontent/runners/tvshowdeletestate"
	"github.com/kkiling/media-delivery/internal/usercase/videocontent/tvshowdelete"
	"github.com/kkiling/statemachine"
	"github.com/samber/lo"
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

func (s *Service) DeleteVideoContentFiles(ctx context.Context, params DeleteVideoContentFilesParams) (*tvshowdeletestate.State, error) {
	if err := s.validateDeleteVideoContentFilesParams(ctx, params); err != nil {
		return nil, fmt.Errorf("validateDeleteVideoContentFilesParams: %w", err)
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
	if content.DeliveryStatus != DeliveryStatusDelivered {
		return nil, fmt.Errorf("video content is not delivered: %w", ucerr.InvalidArgument)
	}

	// Достаем ID стейта доставки
	contentState, find := lo.Find(content.States, func(item State) bool {
		return item.Type == runners.TVShowDelivery
	})
	if !find {
		return nil, ucerr.NotFound
	}
	// Достаем инфу о стейте доставки, что бы вытащить от туда нужную инфу
	deliveryState, err := s.tvShowDeliveryState.GetStateByID(ctx, contentState.StateID)
	if err != nil {
		return nil, fmt.Errorf("tvShowDeliveryState.GetStateByID: %w", err)
	}
	if deliveryState == nil {
		// Нечего удалить, так как не было стейта доставки
		return nil, ucerr.NotFound
	}
	if deliveryState.Status != statemachine.CompletedStatus {
		return nil, fmt.Errorf("delivery state is not completed: %w", ucerr.InvalidArgument)
	}

	// Хеш раздачи (что бы удалить раздачу в торрент клиенте)
	fmt.Printf("href: %s\n", deliveryState.Data.Torrent.MagnetLink.Hash)
	// Путь до торрент файлов
	fmt.Printf("torrent files path: %s\n", deliveryState.Data.TorrentFilesData.ContentFullPath)
	// Путь до файлов медиа сервера (сезона)
	fmt.Printf("media server tvshow season path: %s\n", deliveryState.Data.EpisodesData.TVShowCatalogPath.SeasonPath)
	// ID сезона и сериала
	fmt.Printf(" :%v", deliveryState.MetaData.ContentID.TVShow)

	options := tvshowdeletestate.CreateOptions{
		TVShowID:    *params.ContentID.TVShow,
		MagnetHash:  "", // TODO:
		TorrentPath: "", // TODO:
		TVShowCatalogPath: tvshowdelete.TVShowCatalogPath{
			TVShowPath: "",
			SeasonPath: "",
		}, // TODO:
	}

	var state *tvshowdeletestate.State
	//  TODO: одна транзакция
	{
		state, err = s.tvShowDeleteState.Create(ctx, options)
		if err != nil {
			return nil, fmt.Errorf("tvShowDeleteState.Create: %w", err)
		}

		// Обновить модель стейта videoContent
		updateVideoContent := UpdateVideoContent{
			DeliveryStatus: content.DeliveryStatus,
			States: append(content.States, State{
				StateID: state.ID,
				Type:    runners.TVShowDelete,
			}),
		}

		if err = s.storage.UpdateVideoContent(ctx, content.ID, &updateVideoContent); err != nil {
			return nil, fmt.Errorf("storage.UpdateVideoContent: %w", err)
		}
	}

	return state, nil
}

func (s *Service) GetDeleteData(ctx context.Context, contentID common.ContentID) (*tvshowdeletestate.State, error) {
	if err := contentID.Validate(); err != nil {
		return nil, err
	}

	stateID, err := s.getStateID(ctx, contentID, runners.TVShowDelete)
	if err != nil {
		return nil, fmt.Errorf("getStateID: %w", err)
	}

	result, err := s.tvShowDeleteState.GetStateByID(ctx, stateID)
	if err != nil {
		return nil, fmt.Errorf("s.GetStateByID: %w", err)
	}

	return result, nil
}
