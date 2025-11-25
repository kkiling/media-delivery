package content

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/google/uuid"
	"github.com/kkiling/statemachine"
	"github.com/samber/lo"
	"go.uber.org/zap"

	ucerr "github.com/kkiling/media-delivery/internal/usercase/err"
	"github.com/kkiling/media-delivery/internal/usercase/videocontent/runners"
	"github.com/kkiling/media-delivery/internal/usercase/videocontent/runners/tvshowdeletestate"
	"github.com/kkiling/media-delivery/internal/usercase/videocontent/runners/tvshowdeliverystate"
)

const (
	completeTimeSec      = 3
	getVideoContentLimit = 10
)

func (s *Service) completeActiveVideoContent(ctx context.Context) error {
	defer func() {
		if r := recover(); r != nil {
			s.logger.Error("panic recovered",
				zap.Any("panic_value", r),
				zap.Stack("stack"),
			)
		}
	}()

	// Получаем контент в статусе InProgress
	contents, err := s.storage.GetVideoContentsByDeliveryStatus(ctx, []DeliveryStatus{
		DeliveryStatusInProgress,
		DeliveryStatusDeleting,
	}, getVideoContentLimit)
	if err != nil {
		return fmt.Errorf("storage.GetVideoContents: %w", err)
	}
	for _, content := range contents {
		err = s.completeVideoContent(ctx, content)
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
			if err := s.completeActiveVideoContent(ctx); err != nil {
				s.logger.Errorf("completeActiveVideoContent: %v", err)
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

func handleStateCompleteError(input error) (getActualState bool, err error) {
	if input == nil {
		return false, nil
	}
	getActualState = false
	switch {
	// Можем попасть в такую ситуацию, что стейт ожидал опцию на входе добивания
	// А мы начали добивать без ожидаемой опции, не считаем это ошибкой
	case errors.Is(input, statemachine.ErrOptionsIsUndefined):
		getActualState = true
	// Пока добивали, оказалось что стейт уже в терминальном статусе
	// Заменяем стейт на стейт из базы
	case errors.Is(input, statemachine.ErrInTerminalStatus):
		getActualState = true
	default:
		return false, fmt.Errorf("tvShowDeliveryState.Complete: %w", input)
	}

	if errors.Is(input, statemachine.ErrOptionsIsUndefined) {
		return false, nil
	}
	if errors.Is(input, statemachine.ErrInTerminalStatus) {
		// skip
	} else {
		return false, fmt.Errorf("tvShowDeliveryState.Complete: %w", input)
	}

	return getActualState, nil
}

func (s *Service) completeDeliveryState(ctx context.Context, stateID uuid.UUID) (*tvshowdeliverystate.State, error) {
	// Добиваем стейт
	newState, executeErr, err := s.tvShowDeliveryState.Complete(ctx, stateID)
	getActualState, err := handleStateCompleteError(err)
	if err != nil {
		return nil, err
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

func (s *Service) completeDeleteState(ctx context.Context, stateID uuid.UUID) (*tvshowdeletestate.State, error) {
	// Добиваем стейт
	newState, executeErr, err := s.tvShowDeleteState.Complete(ctx, stateID)
	getActualState, err := handleStateCompleteError(err)
	if err != nil {
		return nil, err
	}

	// Ошибка в логике выполнения выпуска
	if executeErr != nil {
		s.logger.Errorf("executeError: %v", executeErr)
		return nil, fmt.Errorf("tvShowDeliveryState.Complete executeError: %w", executeErr)
	}

	if getActualState {
		// Подтягиваем актуальный стейт из базы
		newState, err = s.tvShowDeleteState.GetStateByID(ctx, stateID)
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

func (s *Service) completeVideoContent(ctx context.Context, content VideoContent) error {
	sort.Slice(content.States, func(i, j int) bool {
		return content.States[i].CreatedAt.Before(content.States[j].CreatedAt)
	})
	lastState, find := lo.Last(content.States)
	if !find {
		return ucerr.NotFound
	}

	needUpdate := false
	updateVideoContent := UpdateVideoContent{
		DeliveryStatus: content.DeliveryStatus,
		States:         content.States,
	}

	s.logger.Debugf("start complete content: %d (season %d) - status %d",
		content.ContentID.TVShow.ID, content.ContentID.TVShow.SeasonNumber, content.DeliveryStatus)

	if content.DeliveryStatus == DeliveryStatusInProgress {
		// Delivery
		deliveryState, err := s.completeDeliveryState(ctx, lastState.StateID)
		if err != nil {
			return fmt.Errorf("completeDeliveryState: %w", err)
		}

		if deliveryState == nil || deliveryState.Type != runners.TVShowDelete {
			// Такого не должно быть
			return ucerr.NotFound
		}

		if deliveryState.Status == statemachine.CompletedStatus {
			needUpdate = true
			updateVideoContent.DeliveryStatus = DeliveryStatusDelivered
		} else if deliveryState.Status == statemachine.FailedStatus {
			needUpdate = true
			updateVideoContent.DeliveryStatus = DeliveryStatusFailed
		}
	} else if content.DeliveryStatus == DeliveryStatusDeleting {
		// Delete
		deliveryState, err := s.completeDeleteState(ctx, lastState.StateID)
		if err != nil {
			return fmt.Errorf("completeDeleteState: %w", err)
		}

		if deliveryState == nil || deliveryState.Type != runners.TVShowDelete {
			// Такого не должно быть
			return ucerr.NotFound
		}

		if deliveryState.Status == statemachine.CompletedStatus {
			needUpdate = true
			updateVideoContent.DeliveryStatus = DeliveryStatusDeleted
		} else if deliveryState.Status == statemachine.FailedStatus {
			needUpdate = true
			updateVideoContent.DeliveryStatus = DeliveryStatusFailed
		}
	} else {
		// Такого не должно быть
		return ucerr.InvalidArgument
	}

	if needUpdate {
		// Обновляем VideoContent статус
		s.logger.Debugf("update video content: %d (season %d)", content.ContentID.TVShow.ID, content.ContentID.TVShow.SeasonNumber)
		err := s.storage.UpdateVideoContent(ctx, content.ID, &updateVideoContent)
		if err != nil {
			return fmt.Errorf("storage.UpdateVideoContent: %w", err)
		}
	}

	return nil
}
