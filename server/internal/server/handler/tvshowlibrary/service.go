package tvshowlibrary

import (
	"context"
	"net/http"

	"github.com/kkiling/goplatform/log"
	"github.com/kkiling/goplatform/server"
	"google.golang.org/grpc"

	"github.com/kkiling/media-delivery/internal/usercase/tvshowlibrary"
	desc "github.com/kkiling/media-delivery/pkg/gen/media-delivery"
)

// TVShowLibrary юзеркейс работы с библиотекой сериалов
type TVShowLibrary interface {
	SearchTVShow(ctx context.Context, params tvshowlibrary.TVShowSearchParams) (*tvshowlibrary.TVShowSearchResult, error)
	GetTVShowInfo(ctx context.Context, params tvshowlibrary.GetTVShowParams) (*tvshowlibrary.GetTVShowResult, error)
	GetSeasonInfo(ctx context.Context, params tvshowlibrary.GetSeasonInfoParams) (*tvshowlibrary.GetSeasonInfoResult, error)
	GetTVShowsFromLibrary(ctx context.Context, params tvshowlibrary.GetTVShowsFromLibraryParams) (*tvshowlibrary.GetTVShowsFromLibraryResult, error)
}

type Handler struct {
	desc.TVShowLibraryServiceServer
	logger        log.Logger
	tvShowLibrary TVShowLibrary
}

// NewHandler новый хендлер
func NewHandler(logger log.Logger,
	tvShowLibrary TVShowLibrary,
) *Handler {
	return &Handler{
		logger:        logger.Named("tv_show_library"),
		tvShowLibrary: tvShowLibrary,
	}
}

// RegistrationServerHandlers .
func (h *Handler) RegistrationServerHandlers(_ *http.ServeMux) {
}

// RegisterServiceHandlerFromEndpoint .
func (h *Handler) RegisterServiceHandlerFromEndpoint() server.HandlerFromEndpoint {
	return desc.RegisterTVShowLibraryServiceHandlerFromEndpoint
}

// RegisterServiceServer регистрация
func (h *Handler) RegisterServiceServer(server *grpc.Server) {
	desc.RegisterTVShowLibraryServiceServer(server, h)
}
