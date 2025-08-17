package container

import (
	"fmt"

	"github.com/kkiling/goplatform/log"
	"github.com/kkiling/goplatform/storagebase/sqlitebase"
	"github.com/kkiling/statemachine"

	"github.com/kkiling/media-delivery/internal/adapter/emby"
	prepareTVShow "github.com/kkiling/media-delivery/internal/adapter/matchtvshow"
	"github.com/kkiling/media-delivery/internal/adapter/mkvmerge"
	mkvsqlite "github.com/kkiling/media-delivery/internal/adapter/mkvmerge/storage/sqlite"
	"github.com/kkiling/media-delivery/internal/adapter/qbittorrent"
	"github.com/kkiling/media-delivery/internal/adapter/rutracker"
	"github.com/kkiling/media-delivery/internal/adapter/themoviedb"
	"github.com/kkiling/media-delivery/internal/config"
	"github.com/kkiling/media-delivery/internal/usercase/tvshowlibrary"
	tvShowLibrarySqlite "github.com/kkiling/media-delivery/internal/usercase/tvshowlibrary/storage/sqlite"
	contentDelivery "github.com/kkiling/media-delivery/internal/usercase/videocontent/content"
	contentSqlite "github.com/kkiling/media-delivery/internal/usercase/videocontent/content/storage/sqlite"
	"github.com/kkiling/media-delivery/internal/usercase/videocontent/delivery"
	"github.com/kkiling/media-delivery/internal/usercase/videocontent/runners/tvshowdeliverystate"
)

type Container struct {
	logger           log.Logger
	tvShowLibrary    *tvshowlibrary.Service
	contentDelivery  *contentDelivery.Service
	mkvMergePipeline *mkvmerge.Pipeline
}

func NewContainer(cfg *config.AppConfig) (*Container, error) {
	logger := log.NewLogger(log.Level(cfg.Server.LogLevel))

	// Storage
	tvShowLibraryStorage, err := tvShowLibrarySqlite.NewStorage(sqlitebase.Config{
		DSN: cfg.Sqlite.SqliteDsn,
	}, logger)
	if err != nil {
		return nil, fmt.Errorf("sqlite.NewStorage: %w", err)
	}

	stateStorage, err := statemachine.NewSqliteStorage(statemachine.SqliteConfig{
		DSN: cfg.Sqlite.SqliteDsn,
	}, logger)
	if err != nil {
		return nil, fmt.Errorf("sqlite.NewStorage: %w", err)
	}

	contentStorage, err := contentSqlite.NewStorage(sqlitebase.Config{
		DSN: cfg.Sqlite.SqliteDsn,
	}, logger)
	if err != nil {
		return nil, fmt.Errorf("sqlite.NewStorage: %w", err)
	}

	mkvPipelineStorage, err := mkvsqlite.NewStorage(sqlitebase.Config{
		DSN: cfg.Sqlite.SqliteDsn,
	}, logger)
	if err != nil {
		return nil, fmt.Errorf("sqlite.NewStorage: %w", err)
	}

	// Adapter
	themoviedbApi, err := themoviedb.NewApi(
		cfg.TheMovieDb.ApiKey,
		logger,
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
		cfg.Rutracker.Username,
		cfg.Rutracker.Password,
		cfg.Rutracker.CookieDir,
	)
	if err != nil {
		return nil, fmt.Errorf("rutracker.NewApi: %w", err)
	}

	mkvMerge := mkvmerge.NewMerge(logger)
	mkvPipeline := mkvmerge.NewPipeline(mkvMerge, mkvPipelineStorage, logger)

	prepareTVShowService := prepareTVShow.NewService(mkvMerge)

	// UserCase
	tvShowLibrary := tvshowlibrary.NewService(tvShowLibraryStorage, themoviedbApi)

	deliveryService := delivery.NewService(
		delivery.Config{
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
	tvShowDeliveryStateMachine := tvshowdeliverystate.NewState(deliveryService, stateStorage)

	deliveryContent := contentDelivery.NewService(logger, contentStorage, tvShowLibrary, tvShowDeliveryStateMachine)

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
