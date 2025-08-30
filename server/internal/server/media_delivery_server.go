package server

import (
	"context"

	"github.com/kkiling/goplatform/log"
	"github.com/kkiling/goplatform/server"

	"github.com/kkiling/media-delivery/internal/server/handler/tvshowlibrary"
	"github.com/kkiling/media-delivery/internal/server/handler/videocontent"
)

// MediaDeliveryServer сервер
type MediaDeliveryServer struct {
	*CustomServer
}

// NewMediaDeliveryServer новый сервер
func NewMediaDeliveryServer(
	logger log.Logger,
	cfg server.Config,
	tvShowLibrary tvshowlibrary.TVShowLibrary,
	videoContent videocontent.VideoContent,
) *MediaDeliveryServer {
	return &MediaDeliveryServer{
		CustomServer: NewCustomServer(
			logger,
			cfg,
			tvshowlibrary.NewHandler(logger, tvShowLibrary),
			videocontent.NewHandler(logger, videoContent),
		),
	}
}

// Start старт сервера
func (s *MediaDeliveryServer) Start(ctx context.Context) error {
	return s.CustomServer.Start(ctx, "api")
}
