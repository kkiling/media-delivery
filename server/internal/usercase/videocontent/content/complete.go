package content

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/google/uuid"
	"github.com/kkiling/media-delivery/internal/usercase/videocontent/runners/tvshowdeliverystate"
	"github.com/kkiling/statemachine"
	"github.com/samber/lo"
	"go.uber.org/zap"

	ucerr "github.com/kkiling/media-delivery/internal/usercase/err"
	"github.com/kkiling/media-delivery/internal/usercase/videocontent/runners"
)

const (
	completeTimeSec      = 3
	getVideoContentLimit = 10
)

func (s *Service) completeState(ctx context.Context, stateID uuid.UUID) (*tvshowdeliverystate.State, error) {
	// Добиваем стейт
	newState, executeErr, err := s.tvShowDeliveryState.Complete(ctx, stateID)

	getActualState := false
	if err != nil {
		switch {
		// Можем попасть в такую ситуацию, что стейт ожидал опцию на входе добивания
		// А мы начали добивать без ожидаемой опции, не считаем это ошибкой
		case errors.Is(err, statemachine.ErrOptionsIsUndefined):
			getActualState = true
		// Пока добивали, оказалось что стейт уже в терминальном статусе
		// Заменяем стейт на стейт из базы
		case errors.Is(err, statemachine.ErrInTerminalStatus):
			getActualState = true
		default:
			return nil, fmt.Errorf("tvShowDeliveryState.Complete: %w", err)
		}

		if errors.Is(err, statemachine.ErrOptionsIsUndefined) {
			return nil, nil
		}
		if errors.Is(err, statemachine.ErrInTerminalStatus) {

		} else {
			return nil, fmt.Errorf("tvShowDeliveryState.Complete: %w", err)
		}
	}

	// Ошибка в логике выполнения выпуска
	if executeErr != nil {
		s.logger.Errorf("executeError: %v", executeErr)
		return nil, fmt.Errorf("tvShowDeliveryState.Complete executeError: %w", executeErr)
	}

	if getActualState {
		// Подтягиваем актуальный стейт из базы
		newState, err = s.tvShowDeliveryState.GetStateByID(ctx, stateID)
		if err != nil {
			return nil, fmt.Errorf("tvShowDeliveryState.GetStateByID: %w", err)
		}
	}

	if newState.Status == statemachine.CompletedStatus {
		s.logger.Debugf("state is completed")
	} else if newState.Status == statemachine.FailedStatus {
		s.logger.Debugf("state is failed")
	} else {
		s.logger.Debugf("state step: %s", newState.Step)
	}

	return newState, nil
}

func (s *Service) completeInProgressTVShowDelivery(ctx context.Context, content VideoContent) error {
	stateInfo, find := lo.Find(content.States, func(item State) bool {
		return item.Type == runners.TVShowDelivery
	})
	if !find {
		// Не должно быть такой ситуации что видео контент в статусе InProgress а мы не нашли стейт доставки
		s.logger.Warn("can't find TVShowDelivery state for tvShow: %d (season %d)",
			content.ContentID.TVShow.ID, content.ContentID.TVShow.SeasonNumber)
		return ucerr.NotFound
	}
	s.logger.Debugf("start complete content: %d (season %d) - status %d",
		content.ContentID.TVShow.ID, content.ContentID.TVShow.SeasonNumber, content.DeliveryStatus)

	deliveryState, err := s.completeState(ctx, stateInfo.StateID)
	if err != nil {
		return fmt.Errorf("completeState: %w", err)
	}
	if deliveryState == nil {
		return ucerr.NotFound
	}

	needUpdate := false
	updateVideoContent := UpdateVideoContent{
		DeliveryStatus: content.DeliveryStatus,
		States:         content.States,
	}
	if deliveryState.Status == statemachine.CompletedStatus {
		needUpdate = true
		updateVideoContent.DeliveryStatus = DeliveryStatusDelivered
	} else if deliveryState.Status == statemachine.FailedStatus {
		needUpdate = true
		updateVideoContent.DeliveryStatus = DeliveryStatusFailed
	}

	if needUpdate {
		// Обновляем VideoContent статус
		s.logger.Debugf("update video content: %d (season %d)", content.ContentID.TVShow.ID, content.ContentID.TVShow.SeasonNumber)
		err = s.storage.UpdateVideoContent(ctx, content.ID, &updateVideoContent)
		if err != nil {
			return fmt.Errorf("storage.UpdateVideoContent: %w", err)
		}
	}

	return nil
}

func (s *Service) completeTVShowDeliveries(ctx context.Context) error {
	defer func() {
		if r := recover(); r != nil {
			s.logger.Error("panic recovered",
				zap.Any("panic_value", r),
				zap.Stack("stack"),
			)
		}
	}()

	// Получаем контент в статусе InProgress
	contents, err := s.storage.GetVideoContentsByDeliveryStatus(ctx, DeliveryStatusInProgress, getVideoContentLimit)
	if err != nil {
		return fmt.Errorf("storage.GetVideoContents: %w", err)
	}
	for _, content := range contents {
		err = s.completeInProgressTVShowDelivery(ctx, content)
		if err != nil {
			return fmt.Errorf("completeTVShowDelivery: %w", err)
		}
	}

	return nil
}

func (s *Service) Complete(ctx context.Context) error {
	scheduler := gocron.NewScheduler(time.UTC)

	// Настраиваем выполнение в 1 поток (по умолчанию и так последовательно)
	scheduler.SetMaxConcurrentJobs(1, gocron.WaitMode)

	// Запускаем задачу каждые 3 секунды
	_, err := scheduler.Every(completeTimeSec).Seconds().Do(func() {
		select {
		case <-ctx.Done(): // Если контекст отменён, выходим
			return
		default:
			if err := s.completeTVShowDeliveries(ctx); err != nil {
				s.logger.Errorf("completeTVShowDeliveries: %v", err)
			}
		}
	})
	if err != nil {
		return fmt.Errorf("failed run cron complete: %w", err)
	}

	// Запускаем планировщик (асинхронно)
	scheduler.StartAsync()

	// Ждём отмены контекста
	<-ctx.Done()

	// Останавливаем планировщик при завершении
	scheduler.Stop()
	return nil
}
