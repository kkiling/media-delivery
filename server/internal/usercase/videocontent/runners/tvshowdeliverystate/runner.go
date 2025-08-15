package tvshowdeliverystate

import (
	"context"
	"fmt"
	"reflect"

	"github.com/kkiling/statemachine"
	"github.com/samber/lo"

	ucerr "github.com/kkiling/torrent-to-media-server/internal/usercase/err"
	"github.com/kkiling/torrent-to-media-server/internal/usercase/videocontent/common"
	"github.com/kkiling/torrent-to-media-server/internal/usercase/videocontent/delivery"
	"github.com/kkiling/torrent-to-media-server/internal/usercase/videocontent/runners"
)

type Runner struct {
	contentDelivery ContentDelivery
}

func NewTaskRunner(contentDelivery ContentDelivery) *Runner {
	return &Runner{
		contentDelivery: contentDelivery,
	}
}

func (r *Runner) Create(_ context.Context, options CreateOptions) (CreateState, error) {
	// Логика создания задачи
	data := TVShowDeliveryData{}

	return CreateState{
		FirstStep: GenerateSearchQuery,
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
		Steps: map[StepDelivery]Step{
			GenerateSearchQuery: {
				OnStep: func(ctx context.Context, stepContext StepContext) *StepResult {
					// Генерация поискового запроса
					data := stepContext.State.Data

					res, err := r.contentDelivery.GenerateSearchQuery(ctx, delivery.GenerateSearchQueryParams{
						TVShowID: *stepContext.State.MetaData.ContentID.TVShow,
					})
					if err != nil {
						return stepContext.Error(fmt.Errorf("GenerateSearchQuery: %w", err))
					}

					data.SearchQuery = res
					return stepContext.Next(SearchTorrents).WithData(data)
				},
			},
			SearchTorrents: {
				OnStep: func(ctx context.Context, stepContext StepContext) *StepResult {
					// ищем раздачи сезона сериала / фильма
					data := stepContext.State.Data
					res, err := r.contentDelivery.SearchTorrent(ctx, delivery.SearchTorrentParams{
						SearchQuery: data.SearchQuery.Query,
					})
					if err != nil {
						return stepContext.Error(fmt.Errorf("SearchTorrent: %w", err))
					}
					data.TorrentSearch = res
					return stepContext.Next(WaitingUserChoseTorrent).WithData(data)
				},
			},
			WaitingUserChoseTorrent: {
				OnStep: func(ctx context.Context, stepContext StepContext) *StepResult {
					// Ожидаем когда пользователь выберет раздачу
					// Или ожидаем что клиент изменит поисковый запрос, тогда прыгаем на SearchTorrents
					// Получение опций выполнения выпуска
					opts := ChoseTorrentOptions{}
					ok, err := stepContext.GetOptions(&opts)
					if err != nil {
						return stepContext.Error(err)
					}
					if !ok { // Пока не получили опцию, не идем дальше
						return stepContext.Empty()
					}
					if opts.Href == nil && opts.NewSearchQuery == nil {
						return stepContext.Error(fmt.Errorf("either Href or NewSearchQuery must be specified: %w", ucerr.InvalidArgument))
					} else if opts.Href != nil && opts.NewSearchQuery != nil {
						return stepContext.Error(fmt.Errorf("only one of Href or NewSearchQuery can be specified: %w", ucerr.InvalidArgument))
					}

					data := stepContext.State.Data

					// Пользователь изменил поисковый запрос
					if opts.NewSearchQuery != nil {
						// Снова производим поиск по раздачам
						data.SearchQuery.Query = *opts.NewSearchQuery
						return stepContext.Next(SearchTorrents).WithData(data)
					}
					if opts.Href != nil {
						// Пользователь выбрал раздачу для скачивания
						// Проверяем что клиент выбрал href из списка
						contains := lo.ContainsBy(data.TorrentSearch, func(item delivery.TorrentSearch) bool {
							return item.Href == *opts.Href
						})
						if !contains {
							return stepContext.Error(fmt.Errorf("no such href: %w", ucerr.InvalidArgument))
						}

						data.Torrent = &delivery.Torrent{
							Href: *opts.Href,
						}

						return stepContext.Next(GetMagnetLink).WithData(data)
					}
					return stepContext.Error(fmt.Errorf("unknow deliverystate: %w", ucerr.InvalidArgument))
				},
				OptionsType: reflect.TypeOf(ChoseTorrentOptions{}),
			},
			GetMagnetLink: {
				OnStep: func(ctx context.Context, stepContext StepContext) *StepResult {
					// Получение магнет ссылки
					data := stepContext.State.Data
					res, err := r.contentDelivery.GetMagnetLink(ctx, delivery.GetMagnetLinkParams{
						Href: data.Torrent.Href,
					})

					if err != nil {
						return stepContext.Error(fmt.Errorf("GetMagnetLink: %w", err))
					}

					data.Torrent.MagnetLink = res
					return stepContext.Next(AddTorrentToTorrentClient).WithData(data)
				},
			},
			AddTorrentToTorrentClient: {
				OnStep: func(ctx context.Context, stepContext StepContext) *StepResult {
					//  Добавление раздачи для скачивания торрент клиентом
					data := stepContext.State.Data
					err := r.contentDelivery.AddTorrentToTorrentClient(ctx, delivery.AddTorrentParams{
						TVShowID: *stepContext.State.MetaData.ContentID.TVShow,
						Magnet:   data.Torrent.MagnetLink.Magnet,
					})
					if err != nil {
						return stepContext.Error(fmt.Errorf("AddTorrentToTorrentClient: %w", err))
					}
					return stepContext.Next(WaitingTorrentFiles)
				},
			},
			// Ожидание когда появится информация о файлах в раздаче
			WaitingTorrentFiles: {
				OnStep: func(ctx context.Context, stepContext StepContext) *StepResult {
					//  Добавление раздачи для скачивания торрент клиентом
					data := stepContext.State.Data
					res, err := r.contentDelivery.WaitingTorrentFiles(ctx, delivery.WaitingTorrentFilesParams{
						Hash: data.Torrent.MagnetLink.Hash,
					})
					if err != nil {
						return stepContext.Error(fmt.Errorf("WaitingTorrentFiles: %w", err))
					}
					if res == nil {
						return stepContext.Empty()
					}
					data.TorrentFilesData = res
					return stepContext.Next(GetEpisodesData).WithData(data)
				},
			},
			GetEpisodesData: {
				OnStep: func(ctx context.Context, stepContext StepContext) *StepResult {
					//  Добавление раздачи для скачивания торрент клиентом
					data := stepContext.State.Data
					res, err := r.contentDelivery.GetEpisodesData(ctx, delivery.GetEpisodesDataParams{
						TVShowID: *stepContext.State.MetaData.ContentID.TVShow,
					})
					if err != nil {
						return stepContext.Error(fmt.Errorf("GetEpisodesData: %w", err))
					}
					data.EpisodesData = res
					return stepContext.Next(PrepareFileMatches).WithData(data)
				},
			},
			PrepareFileMatches: {
				OnStep: func(ctx context.Context, stepContext StepContext) *StepResult {
					// Получение информации о файлах раздачи
					data := stepContext.State.Data
					contentMatches, err := r.contentDelivery.PrepareFileMatches(ctx, delivery.PreparingFileMatchesParams{
						TorrentFiles: data.TorrentFilesData.Files,
						Episodes:     data.EpisodesData.Episodes,
						TVShowID:     *stepContext.State.MetaData.ContentID.TVShow,
					})
					if err != nil {
						return stepContext.Error(fmt.Errorf("PrepareFileMatches: %w", err))
					}
					if len(contentMatches) == 0 {
						return stepContext.Empty()
					}

					data.ContentMatches = contentMatches
					// Определение необходимости конвертации файлов
					if r.contentDelivery.NeedPrepareFileMatches(data.ContentMatches) {
						return stepContext.Next(WaitingChoseFileMatches).WithData(data)
					}
					// Если нечего конвертировать то идем сразу на ожидание скачивания
					return stepContext.Next(WaitingTorrentDownloadComplete).WithData(data)
				},
			},
			WaitingChoseFileMatches: {
				OnStep: func(ctx context.Context, stepContext StepContext) *StepResult {
					// ожидание подтверждения пользователем соответствий выбора файлов
					opts := ChoseFileMatchesOptions{}
					ok, err := stepContext.GetOptions(&opts)
					if err != nil {
						return stepContext.Error(err)
					}
					if !ok { // Пока не получили опцию, не идем дальше
						return stepContext.Empty()
					}
					if !opts.Approve {
						return stepContext.Empty()
					}
					// TODO: выбор пользовтелем другого сопоставления

					return stepContext.Next(WaitingTorrentDownloadComplete)
				},
				OptionsType: reflect.TypeOf(ChoseFileMatchesOptions{}),
			},
			WaitingTorrentDownloadComplete: {
				OnStep: func(ctx context.Context, stepContext StepContext) *StepResult {
					// Ожидание когда торрент докачается до конца
					data := stepContext.State.Data
					res, err := r.contentDelivery.WaitingTorrentDownloadComplete(ctx, delivery.WaitingTorrentDownloadCompleteParams{
						Hash: data.Torrent.MagnetLink.Hash,
					})
					if err != nil {
						return stepContext.Error(fmt.Errorf("PrepareFileMatches: %w", err))
					}
					data.TorrentDownloadStatus = res
					if res.IsComplete {
						return stepContext.Next(CreateVideoContentCatalogs).WithData(data)
					}
					return stepContext.Empty().WithData(data)
				},
			},
			CreateVideoContentCatalogs: {
				OnStep: func(ctx context.Context, stepContext StepContext) *StepResult {
					// Формирование каталогов и иерархии файлов
					data := stepContext.State.Data
					err := r.contentDelivery.CreateContentCatalogs(ctx, delivery.CreateContentCatalogsParams{
						TVShowCatalogPath: data.EpisodesData.TVShowCatalogPath,
					})
					if err != nil {
						return stepContext.Error(fmt.Errorf("CreateContentCatalogs: %w", err))
					}

					return stepContext.Next(DeterminingNeedConvertFiles)
				},
			},
			DeterminingNeedConvertFiles: {
				OnStep: func(ctx context.Context, stepContext StepContext) *StepResult {
					// Определение необходимости конвертации файлов
					data := stepContext.State.Data
					if r.contentDelivery.NeedPrepareFileMatches(data.ContentMatches) {
						return stepContext.Next(StartMergeVideoFiles).WithData(data)
					}
					return stepContext.Next(CreateHardLinkCopy)
				},
			},
			CreateHardLinkCopy: {
				OnStep: func(ctx context.Context, stepContext StepContext) *StepResult {
					// Копирование файлов из раздачи в каталог медиасервера (точнее создание симлинков)
					data := stepContext.State.Data
					if err := r.contentDelivery.CreateHardLinkCopyToMediaServer(ctx, delivery.CreateHardLinkCopyParams{
						ContentMatches: data.ContentMatches,
					}); err != nil {
						return stepContext.Error(fmt.Errorf("CreateHardLinkCopyToMediaServer: %w", err))
					}

					return stepContext.Next(GetCatalogsSize)
				},
			},
			StartMergeVideoFiles: {
				OnStep: func(ctx context.Context, stepContext StepContext) *StepResult {
					//  Конвертирование файлов - полученные файлы сразу сохраняются в каталог медиасервера
					data := stepContext.State.Data

					//  Конвертирование файлов - полученные файлы сразу сохраняются в каталог медиасервера
					mergeIDs, err := r.contentDelivery.StartMergeVideo(ctx, delivery.MergeVideoParams{
						IdempotencyKey: stepContext.State.ID.String(),
						ContentMatches: data.ContentMatches,
					})
					if err != nil {
						return stepContext.Error(fmt.Errorf("StartMergeVideoFiles: %w", err))
					}
					data.MergeIDs = mergeIDs
					return stepContext.Next(WaitingMergeVideoFiles).WithData(data)
				},
			},
			WaitingMergeVideoFiles: {
				OnStep: func(ctx context.Context, stepContext StepContext) *StepResult {
					data := stepContext.State.Data

					//  Конвертирование файлов - полученные файлы сразу сохраняются в каталог медиасервера
					status, err := r.contentDelivery.GetMergeVideoStatus(ctx, data.MergeIDs)
					if err != nil {
						return stepContext.Error(fmt.Errorf("WaitingMergeVideoFiles: %w", err))
					}

					data.MergeVideoStatus = status
					if status.IsComplete {
						if len(status.Errors) == 0 {
							return stepContext.Next(SetVideoFileGroup).WithData(data)
						}
						return stepContext.Error(fmt.Errorf("merge videos contains errors")).WithData(data)
					}
					return stepContext.Empty().WithData(data)
				},
			},
			SetVideoFileGroup: {
				OnStep: func(ctx context.Context, stepContext StepContext) *StepResult {
					data := stepContext.State.Data
					files := lo.Map(data.ContentMatches, func(item delivery.ContentMatches, _ int) string {
						return item.Episode.FileName
					})
					// Установка группы файлам
					err := r.contentDelivery.SetVideoFileGroup(ctx, files)
					if err != nil {
						return stepContext.Error(fmt.Errorf("SetVideoFileGroup: %w", err))
					}

					return stepContext.Next(GetCatalogsSize)
				},
			},
			GetCatalogsSize: {
				OnStep: func(ctx context.Context, stepContext StepContext) *StepResult {
					data := stepContext.State.Data

					data.TVShowCatalogInfo = &delivery.TVShowCatalog{
						TorrentPath:              data.TorrentFilesData.ContentFullPath,
						MediaServerPath:          data.EpisodesData.TVShowCatalogPath,
						IsCopyFilesInMediaServer: r.contentDelivery.NeedPrepareFileMatches(data.ContentMatches),
					}

					// Получение размера каталогов
					torrentTVShowSeasonSize, err := r.contentDelivery.GetCatalogSize(ctx, data.TVShowCatalogInfo.TorrentPath)
					if err != nil {
						return stepContext.Error(fmt.Errorf("contentDelivery.GetCatalogSize: %w", err))
					}
					mediaServerTVShowSize, err := r.contentDelivery.GetCatalogSize(ctx, data.TVShowCatalogInfo.MediaServerPath.FullSeasonPath())
					if err != nil {
						return stepContext.Error(fmt.Errorf("contentDelivery.GetCatalogSize: %w", err))
					}

					data.TVShowCatalogInfo.TorrentSize = torrentTVShowSeasonSize
					data.TVShowCatalogInfo.MediaServerSize = mediaServerTVShowSize

					return stepContext.Next(SetMediaMetaData).WithData(data)
				},
			},
			SetMediaMetaData: {
				OnStep: func(ctx context.Context, stepContext StepContext) *StepResult {
					data := stepContext.State.Data
					// Установка группы файлам
					err := r.contentDelivery.SetMediaMetaData(ctx, delivery.SetMediaMetaDataParams{
						TVShowPath: data.TVShowCatalogInfo.MediaServerPath.TVShowPath,
						TVShowID:   *stepContext.State.MetaData.ContentID.TVShow,
					})
					if err != nil {
						return stepContext.Error(fmt.Errorf("SetMediaMetaData: %w", err))
					}

					return stepContext.Complete()
				},
			},
		},
	}
}
