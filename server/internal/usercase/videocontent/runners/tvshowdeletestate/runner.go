package tvshowdeletestate

import (
	"context"
	"fmt"

	"github.com/kkiling/statemachine"

	"github.com/kkiling/media-delivery/internal/common"
	"github.com/kkiling/media-delivery/internal/usercase/videocontent/runners"
)

type Runner struct {
	contentDeleted ContentDeleted
}

func NewTaskRunner(contentDelivery ContentDeleted) *Runner {
	return &Runner{
		contentDeleted: contentDelivery,
	}
}

func (r *Runner) Create(_ context.Context, options CreateOptions) (CreateState, error) {
	// Логика создания задачи
	data := TVShowDeleteData{
		MagnetHash:        options.MagnetHash,
		TorrentPath:       options.TorrentPath,
		TVShowCatalogPath: options.TVShowCatalogPath,
	}

	return CreateState{
		FirstStep: StartDeleteTVShowSeason,
		Data:      data,
		MetaData: runners.Metadata{
			ContentID: common.ContentID{
				TVShow: &options.TVShowID,
			},
		},
	}, nil
}

func (r *Runner) Type() runners.Type {
	return runners.TVShowDelivery
}

func (r *Runner) StepRegistration(_ statemachine.StepRegistrationParams) StepRegistration {
	return StepRegistration{
		Steps: map[StepDelete]Step{
			StartDeleteTVShowSeason: {
				OnStep: func(ctx context.Context, stepContext StepContext) *StepResult {
					// начальный шаг
					return stepContext.Next(DeleteTorrentFromTorrentClient)
				},
			},
			DeleteTorrentFromTorrentClient: {
				OnStep: func(ctx context.Context, stepContext StepContext) *StepResult {
					// удаление торрент раздачи из торрент клиента
					data := stepContext.State.Data
					err := r.contentDeleted.DeleteTorrentFromTorrentClient(ctx, data.MagnetHash)
					if err != nil {
						return stepContext.Error(fmt.Errorf("DeleteTorrentFromTorrentClient: %w", err))
					}
					return stepContext.Next(DeleteTorrentFiles)
				},
			},
			DeleteTorrentFiles: {
				OnStep: func(ctx context.Context, stepContext StepContext) *StepResult {
					//  удаление файлов раздачи с диска
					data := stepContext.State.Data
					err := r.contentDeleted.DeleteTorrentFiles(ctx, data.TorrentPath)
					if err != nil {
						return stepContext.Error(fmt.Errorf("DeleteTorrentFiles: %w", err))
					}
					return stepContext.Next(DeleteSeasonFiles)
				},
			},
			DeleteSeasonFiles: {
				OnStep: func(ctx context.Context, stepContext StepContext) *StepResult {
					//  удаление файлов сезона с медиасервера
					data := stepContext.State.Data
					err := r.contentDeleted.DeleteSeasonFiles(ctx, data.TVShowCatalogPath)
					if err != nil {
						return stepContext.Error(fmt.Errorf("DeleteSeasonFiles: %w", err))
					}
					return stepContext.Next(DeleteSeasonFromMediaServer)
				},
			},
			DeleteSeasonFromMediaServer: {
				OnStep: func(ctx context.Context, stepContext StepContext) *StepResult {
					// удаление сезона сериала с медиасервера
					data := stepContext.State.Data
					err := r.contentDeleted.DeleteSeasonFromMediaServer(ctx, data.TVShowCatalogPath)
					if err != nil {
						return stepContext.Error(fmt.Errorf("DeleteSeasonFromMediaServer: %w", err))
					}
					return stepContext.Next(DeleteLabel)
				},
			},
			DeleteLabel: {
				OnStep: func(ctx context.Context, stepContext StepContext) *StepResult {
					data := stepContext.State.MetaData
					err := r.contentDeleted.DeleteLabelHasVideoContentFiles(ctx, data.ContentID)
					if err != nil {
						return stepContext.Error(fmt.Errorf("DeleteLabelHasVideoContentFiles: %w", err))
					}
					return stepContext.Complete()
				},
			},
		},
	}
}
