package content

import (
	"context"
	"errors"
	"fmt"

	"github.com/kkiling/goplatform/log"
	"github.com/kkiling/goplatform/storagebase"

	"github.com/kkiling/media-delivery/internal/common"
	ucerr "github.com/kkiling/media-delivery/internal/usercase/err"
)

type Service struct {
	logger              log.Logger
	storage             Storage
	tvShowLibrary       TVShowLibrary
	tvShowDeliveryState TVShowDeliveryState
	tvShowDeleteState   TVShowDeleteState
	labels              Labels
	clock               Clock
	uuidGenerator       UUIDGenerator
}

func NewService(
	logger log.Logger,
	storage Storage,
	tvShowLibrary TVShowLibrary,
	tvShowDeliveryState TVShowDeliveryState,
	tvShowDeleteState TVShowDeleteState,
	labels Labels,
) *Service {
	return &Service{
		storage:             storage,
		tvShowLibrary:       tvShowLibrary,
		tvShowDeliveryState: tvShowDeliveryState,
		tvShowDeleteState:   tvShowDeleteState,
		labels:              labels,
		clock:               &common.RealClock{},
		uuidGenerator:       &common.UUIDGenerator{},
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
