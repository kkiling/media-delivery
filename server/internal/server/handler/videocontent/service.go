package videocontent

import (
	"context"
	"net/http"

	"github.com/kkiling/goplatform/log"
	"github.com/kkiling/goplatform/server"
	"google.golang.org/grpc"

	"github.com/kkiling/media-delivery/internal/usercase/videocontent"
	desc "github.com/kkiling/media-delivery/pkg/gen/media-delivery"
)

// VideoContent юзеркейс работы с доставкой видео файлов
type VideoContent interface {
	CreateVideoContent(ctx context.Context, params videocontent.CreateVideoContentParams) (*videocontent.VideoContent, error)
	GetVideoContent(ctx context.Context, contentID videocontent.ContentID) ([]videocontent.VideoContent, error)
	CreateDeliveryState(ctx context.Context, params videocontent.DeliveryVideoContentParams) (*videocontent.TVShowDeliveryState, error)
	GetDeliveryData(ctx context.Context, contentID videocontent.ContentID) (*videocontent.TVShowDeliveryState, error)
	ChoseTorrentOptions(ctx context.Context, contentID videocontent.ContentID, opts videocontent.ChoseTorrentOptions) (*videocontent.TVShowDeliveryState, error)
	ChoseFileMatchesOptions(ctx context.Context, contentID videocontent.ContentID, opts videocontent.ChoseFileMatchesOptions) (*videocontent.TVShowDeliveryState, error)
	CreateDeleteState(ctx context.Context, params videocontent.CreateDeleteStateParams) (*videocontent.TVShowDeleteState, error)
	GetDeleteData(ctx context.Context, contentID videocontent.ContentID) (*videocontent.TVShowDeleteState, error)
}

type Handler struct {
	desc.VideoContentServiceServer
	logger       log.Logger
	videoContent VideoContent
}

// NewHandler новый хендлер
func NewHandler(logger log.Logger,
	videoContent VideoContent,
) *Handler {
	return &Handler{
		logger:       logger.Named("video_content"),
		videoContent: videoContent,
	}
}

// RegistrationServerHandlers .
func (h *Handler) RegistrationServerHandlers(_ *http.ServeMux) {
}

// RegisterServiceHandlerFromEndpoint .
func (h *Handler) RegisterServiceHandlerFromEndpoint() server.HandlerFromEndpoint {
	return desc.RegisterVideoContentServiceHandlerFromEndpoint
}

// RegisterServiceServer регистрация
func (h *Handler) RegisterServiceServer(server *grpc.Server) {
	desc.RegisterVideoContentServiceServer(server, h)
}
