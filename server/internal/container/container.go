package container

import (
	"context"
	"fmt"

	"github.com/kkiling/goplatform/log"
	"github.com/kkiling/goplatform/storagebase/postgrebase"
	"github.com/kkiling/media-delivery/internal/usercase/videocontent/runners/tvshowdeletestate"
	"github.com/kkiling/media-delivery/internal/usercase/videocontent/tvshowdelete"
	"github.com/kkiling/media-delivery/internal/usercase/videocontent/tvshowdelivery"
	"github.com/kkiling/statemachine"

	"github.com/kkiling/media-delivery/internal/adapter/emby"
	prepareTVShow "github.com/kkiling/media-delivery/internal/adapter/matchtvshow"
	"github.com/kkiling/media-delivery/internal/adapter/mkvmerge"
	mkvPostgresql "github.com/kkiling/media-delivery/internal/adapter/mkvmerge/storage/postgresql"
	"github.com/kkiling/media-delivery/internal/adapter/qbittorrent"
	"github.com/kkiling/media-delivery/internal/adapter/rutracker"
	"github.com/kkiling/media-delivery/internal/adapter/themoviedb"
	"github.com/kkiling/media-delivery/internal/config"
	"github.com/kkiling/media-delivery/internal/usercase/labels"
	labelsPostgreSql "github.com/kkiling/media-delivery/internal/usercase/labels/storage/postgresql"
	"github.com/kkiling/media-delivery/internal/usercase/tvshowlibrary"
	tvShowLibraryPostgreSql "github.com/kkiling/media-delivery/internal/usercase/tvshowlibrary/storage/postgresql"
	contentDelivery "github.com/kkiling/media-delivery/internal/usercase/videocontent/content"
	contentPostgreSql "github.com/kkiling/media-delivery/internal/usercase/videocontent/content/storage/postgresql"
	"github.com/kkiling/media-delivery/internal/usercase/videocontent/runners/tvshowdeliverystate"
)

type Container struct {
	logger           log.Logger
	tvShowLibrary    *tvshowlibrary.Service
	contentDelivery  *contentDelivery.Service
	mkvMergePipeline *mkvmerge.Pipeline
}

func NewContainer(ctx context.Context, cfg *config.AppConfig) (*Container, error) {
	logger := log.NewLogger(log.Level(cfg.Server.LogLevel))

	// *** *** ***
	// Storage
	pgPool, err := postgrebase.NewPgConn(ctx, postgrebase.Config{
		ConnString: cfg.Postgresql.ConnString,
	})
	if err != nil {
		return nil, fmt.Errorf("sqlite.NewStorage: %w", err)
	}

	stateStorage := statemachine.NewStorage(pgPool)
	tvShowLibraryStorage := tvShowLibraryPostgreSql.NewStorage(pgPool)
	contentStorage := contentPostgreSql.NewStorage(pgPool)
	mkvPipelineStorage := mkvPostgresql.NewStorage(pgPool)
	labelsStorage := labelsPostgreSql.NewStorage(pgPool)
	// *** *** ***

	// Adapter
	themoviedbApi, err := themoviedb.NewApi(
		logger,
		themoviedb.Config{
			APIKey:   cfg.TheMovieDb.ApiKey,
			ProxyURL: cfg.TheMovieDb.ProxyURL,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("themoviedb.NewApi: %w", err)
	}

	embyApi, err := emby.NewApi(cfg.Emby.ApiKey, cfg.Emby.ApiUrl, logger)
	if err != nil {
		return nil, fmt.Errorf("emby.NewApi: %w", err)
	}

	qBittorrentApi, err := qbittorrent.NewApi(
		logger,
		cfg.QBittorrent.ApiUrl,
		cfg.QBittorrent.Username,
		cfg.QBittorrent.Password,
		cfg.QBittorrent.CookieDir,
	)
	if err != nil {
		return nil, fmt.Errorf("qbittorrent.NewApi: %w", err)
	}

	rutrackerApi, err := rutracker.NewApi(
		logger,
		rutracker.Config{
			Username:   cfg.Rutracker.Username,
			Password:   cfg.Rutracker.Password,
			CookiesDir: cfg.Rutracker.CookiesDir,
			ProxyURL:   cfg.Rutracker.ProxyURL,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("rutracker.NewApi: %w", err)
	}

	mkvMerge := mkvmerge.NewMerge(logger)
	mkvPipeline := mkvmerge.NewPipeline(mkvMerge, mkvPipelineStorage, logger)

	prepareTVShowService := prepareTVShow.NewService()

	labelsService := labels.NewService(labelsStorage)

	// UserCase
	tvShowLibrary := tvshowlibrary.NewService(tvShowLibraryStorage, themoviedbApi)

	tvShowDeliveryService := tvshowdelivery.NewService(
		tvshowdelivery.Config{
			BasePath:                   cfg.DeliveryConfig.BasePath,
			TVShowTorrentSavePath:      cfg.DeliveryConfig.TVShowTorrentSavePath,
			TVShowMediaSaveTvShowsPath: cfg.DeliveryConfig.TVShowMediaSaveTvShowsPath,
		},
		tvShowLibrary,
		rutrackerApi,
		qBittorrentApi,
		embyApi,
		prepareTVShowService,
		mkvPipeline,
	)
	tvShowDeleteService := tvshowdelete.NewService(
		tvshowdelete.Config{
			BasePath:                   cfg.DeliveryConfig.BasePath,
			TVShowTorrentSavePath:      cfg.DeliveryConfig.TVShowTorrentSavePath,
			TVShowMediaSaveTvShowsPath: cfg.DeliveryConfig.TVShowMediaSaveTvShowsPath,
		},
		qBittorrentApi,
		embyApi,
	)

	tvShowDeliveryStateMachine := tvshowdeliverystate.NewState(tvShowDeliveryService, stateStorage)
	tvShowDeleteStateMachine := tvshowdeletestate.NewState(tvShowDeleteService, stateStorage)

	deliveryContent := contentDelivery.NewService(
		logger,
		contentStorage,
		tvShowLibrary,
		tvShowDeliveryStateMachine,
		tvShowDeleteStateMachine,
		labelsService,
	)

	return &Container{
		logger:           logger,
		tvShowLibrary:    tvShowLibrary,
		contentDelivery:  deliveryContent,
		mkvMergePipeline: mkvPipeline,
	}, nil
}

func (c *Container) GetTvShowLibrary() *tvshowlibrary.Service {
	return c.tvShowLibrary
}

func (c *Container) GetContentDelivery() *contentDelivery.Service {
	return c.contentDelivery
}

func (c *Container) MkvMergePipeline() *mkvmerge.Pipeline {
	return c.mkvMergePipeline
}

func (c *Container) GetLogger() log.Logger {
	return c.logger
}
