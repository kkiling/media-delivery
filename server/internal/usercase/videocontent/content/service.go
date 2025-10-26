package content

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/kkiling/goplatform/log"
	"github.com/kkiling/goplatform/storagebase"
	"github.com/kkiling/media-delivery/internal/common"
	ucerr "github.com/kkiling/media-delivery/internal/usercase/err"
	"github.com/kkiling/media-delivery/internal/usercase/videocontent/runners"
	"github.com/kkiling/media-delivery/internal/usercase/videocontent/runners/tvshowdeliverystate"
	"github.com/samber/lo"
)

type Service struct {
	logger              log.Logger
	storage             Storage
	tvShowLibrary       TVShowLibrary
	tvShowDeliveryState TVShowDeliveryState
	labels              Labels
	clock               Clock
	uuidGenerator       UUIDGenerator
}

func NewService(
	logger log.Logger,
	storage Storage,
	tvShowLibrary TVShowLibrary,
	tvShowDeliveryState TVShowDeliveryState,
	labels Labels,
) *Service {
	return &Service{
		storage:             storage,
		tvShowLibrary:       tvShowLibrary,
		tvShowDeliveryState: tvShowDeliveryState,
		labels:              labels,
		clock:               &realClock{},
		uuidGenerator:       &uuidGenerator{},
		logger:              logger.Named("content"),
	}
}

/*
	Какие кейсы:
		- Создание новой доставки файлов с прохождением полного флоу от поиска раздачи до доставки файлов до медиасервера
		- Можно оставить информацию о раздаче (Href и Magnet) но при этом удалить все файлы, что бы не занимали место на диске
		- Потом на основе (Href и Magnet) восстанавливать файлы и скачивать их снова, при этом не запрашивая больше инфу от клиента
			и все подтягивать из старых стейтов (что делать если раздача обновиться?)
       - Раздача может обновиться и запускается процесс обновления раздачи
*/

func (s *Service) GetVideoContent(ctx context.Context, contentID common.ContentID) ([]VideoContent, error) {
	if err := contentID.Validate(); err != nil {
		return nil, err
	}

	result, err := s.storage.GetVideoContents(ctx, contentID)
	switch {
	case err == nil:
	case errors.Is(err, storagebase.ErrNotFound): // Выпуск не найден
		return nil, ucerr.NotFound
	default:
		return nil, fmt.Errorf("storage.GetVideoContent: %w", err)
	}
	return result, nil
}

func (s *Service) getStateID(ctx context.Context, contentID common.ContentID, runersType runners.Type) (uuid.UUID, error) {
	contents, err := s.storage.GetVideoContents(ctx, contentID)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("storage.GetVideoContent: %w", err)
	}

	if len(contents) != 1 {
		return uuid.UUID{}, ucerr.NotFound
	}

	content := contents[0]
	state, find := lo.Find(content.State, func(item State) bool {
		return item.Type == runersType
	})
	if !find {
		return uuid.UUID{}, ucerr.NotFound
	}

	return state.StateID, nil
}

func (s *Service) GetTVShowDeliveryData(ctx context.Context, contentID common.ContentID) (*tvshowdeliverystate.State, error) {
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
