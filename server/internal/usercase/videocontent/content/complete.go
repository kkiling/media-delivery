package content

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-co-op/gocron"
	ucerr "github.com/kkiling/media-delivery/internal/usercase/err"
	"github.com/kkiling/media-delivery/internal/usercase/videocontent/runners"
	"github.com/kkiling/statemachine"
	"github.com/samber/lo"
	"go.uber.org/zap"
)

func (s *Service) completeTVShowDelivery(ctx context.Context, content VideoContent) error {
	stateInfo, find := lo.Find(content.State, func(item State) bool {
		return item.Type == runners.TVShowDelivery
	})
	if !find {
		return ucerr.NotFound
	}
	s.logger.Debugf("start complete content: %d (season %d) - status %s",
		content.ContentID.TVShow.ID, content.ContentID.TVShow.SeasonNumber, content.DeliveryStatus)

	if content.DeliveryStatus == DeliveryStatusInProgress {
		// Добиваем стейт
		newState, executeErr, err := s.tvShowDeliveryState.Complete(ctx, stateInfo.StateID)
		if err != nil && !errors.Is(err, statemachine.ErrOptionsIsUndefined) {
			if errors.Is(err, statemachine.ErrInTerminalStatus) {
				newState, err = s.tvShowDeliveryState.GetStateByID(ctx, stateInfo.StateID)
				if err != nil {
					return fmt.Errorf("tvShowDeliveryState.GetStateByID: %w", err)
				}
			} else {
				return fmt.Errorf("tvShowDeliveryState.Complete: %w", err)
			}
		}
		if executeErr != nil {
			s.logger.Errorf("executeError: %v", executeErr)
			return executeErr
		}

		if newState.Status == statemachine.CompletedStatus {
			s.logger.Debugf("state is completed")
		} else if newState.Status == statemachine.FailedStatus {
			s.logger.Debugf("state is failed")
		} else {
			s.logger.Debugf("state step: %s", newState.Step)
		}

		needUpdate := false
		updateVideoContent := UpdateVideoContent{
			DeliveryStatus: content.DeliveryStatus,
		}
		if newState.Status == statemachine.CompletedStatus {
			needUpdate = true
			updateVideoContent.DeliveryStatus = DeliveryStatusDelivered
		} else if newState.Status == statemachine.FailedStatus {
			needUpdate = true
			updateVideoContent.DeliveryStatus = DeliveryStatusFailed
		}

		if needUpdate {
			s.logger.Debugf("update video content: %d (season %d)", content.ContentID.TVShow.ID, content.ContentID.TVShow.SeasonNumber)
			err = s.storage.UpdateVideoContent(ctx, content.ID, &updateVideoContent)
			if err != nil {
				return fmt.Errorf("storage.UpdateVideoContent: %w", err)
			}
		}
	}

	//Трекаем обновление статуса в процессе доставки in_progress до доставлено delivered на основе стейта
	//	in_progress -> delivered
	//Трекаем по аналогии
	//	updating -> delivered
	//Трекаем по анлогии на основании стейта
	//	deleting -> deleted

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

	contents, err := s.storage.GetVideoContentsByStatus(ctx, DeliveryStatusInProgress, 10)
	if err != nil {
		return fmt.Errorf("storage.GetVideoContents: %w", err)
	}
	for _, content := range contents {
		err = s.completeTVShowDelivery(ctx, content)
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
	_, err := scheduler.Every(3).Seconds().Do(func() {
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
