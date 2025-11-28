package content

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/kkiling/statemachine"
	"github.com/samber/lo"
	"go.uber.org/zap"
)

const (
	completeTimeSec      = 3
	getVideoContentLimit = 10
)

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

func (s *Service) completeActiveVideoContent(ctx context.Context) error {
	defer func() {
		if r := recover(); r != nil {
			s.logger.Error("panic recovered",
				zap.Any("panic_value", r),
				zap.Stack("stack"),
			)
		}
	}()

	var runnerList = []runnerCommon{
		deliveryRunner{s.tvShowDeliveryState},
		deleteRunner{s.tvShowDeleteState},
	}

	statusIn := lo.Map(runnerList, func(item runnerCommon, index int) DeliveryStatus {
		return item.TargetDeliveryStatus()
	})
	statusIn = lo.Uniq(statusIn)

	// Получаем контент в статусе InProgress
	contents, err := s.storage.GetVideoContentsByDeliveryStatus(ctx, statusIn, getVideoContentLimit)
	if err != nil {
		return fmt.Errorf("storage.GetVideoContents: %w", err)
	}

	for _, content := range contents {
		runner, find := lo.Find(runnerList, func(item runnerCommon) bool {
			return item.TargetDeliveryStatus() == content.DeliveryStatus
		})
		if !find {
			return fmt.Errorf("find item %v: %w", content, statemachine.ErrNotFound)
		}

		err = s.completeVideoContent(ctx, content, runner)
		if err != nil {
			return fmt.Errorf("completeVideoContent: %w", err)
		}
	}

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
		return false, fmt.Errorf("complete: %w", input)
	}

	if errors.Is(input, statemachine.ErrOptionsIsUndefined) {
		return false, nil
	}
	if errors.Is(input, statemachine.ErrInTerminalStatus) {
		// skip
	} else {
		return false, fmt.Errorf("complete: %w", input)
	}

	return getActualState, nil
}

func (s *Service) completeVideoContent(ctx context.Context, content VideoContent, runner runnerCommon) error {
	stateID := getLastState(content, runner.RunnerType())
	if stateID == nil {
		return fmt.Errorf("no state ID found for content %v", content.ContentID)
	}
	// Complete
	// Добиваем стейт
	newState, executeErr, err := runner.Complete(ctx, *stateID)
	getActualState, err := handleStateCompleteError(err)
	if err != nil {
		return err
	}
	// Ошибка в логике выполнения выпуска
	if executeErr != nil {
		s.logger.Errorf("executeError: %v", executeErr)
		return fmt.Errorf("runner.Complete executeError: %w", executeErr)
	}
	if getActualState {
		// Подтягиваем актуальный стейт из базы
		newState, err = runner.GetStateByID(ctx, *stateID)
		if err != nil {
			return fmt.Errorf("runner.GetStateByID: %w", err)
		}
	}

	if newState.status == statemachine.CompletedStatus {
		s.logger.Debugf("state is completed")
	} else if newState.status == statemachine.FailedStatus {
		s.logger.Debugf("state is failed")
	} else {
		s.logger.Debugf("state step: %s", newState.step)
	}

	newDeliveryStatus := runner.ToDeliveryStatus(newState.status)
	if newDeliveryStatus != content.DeliveryStatus {
		updateVideoContent := UpdateVideoContent{
			DeliveryStatus: newDeliveryStatus,
			States:         content.States,
		}
		// Обновляем VideoContent статус
		s.logger.Debugf("update video content: %d (season %d)", content.ContentID.TVShow.ID, content.ContentID.TVShow.SeasonNumber)
		err = s.storage.UpdateVideoContent(ctx, content.ID, &updateVideoContent)
		if err != nil {
			return fmt.Errorf("storage.UpdateVideoContent: %w", err)
		}
	}

	return nil
}
