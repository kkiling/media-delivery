package content

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/kkiling/media-delivery/internal/common"
	"github.com/kkiling/media-delivery/internal/usercase/labels"

	"github.com/kkiling/media-delivery/internal/usercase/tvshowlibrary"
	"github.com/kkiling/media-delivery/internal/usercase/videocontent/runners/tvshowdeliverystate"
)

type Storage interface {
	CreateVideoContent(ctx context.Context, videoContent *VideoContent) error
	GetVideoContents(ctx context.Context, contentID common.ContentID) ([]VideoContent, error)
	UpdateVideoContent(ctx context.Context, id uuid.UUID, videoContent *UpdateVideoContent) error
	GetVideoContentsByStatus(ctx context.Context, status DeliveryStatus, limit int) ([]VideoContent, error)
}

type TVShowLibrary interface {
	GetTVShowInfo(ctx context.Context, params tvshowlibrary.GetTVShowParams) (*tvshowlibrary.GetTVShowResult, error)
	GetSeasonInfo(ctx context.Context, params tvshowlibrary.GetSeasonInfoParams) (*tvshowlibrary.GetSeasonInfoResult, error)
	AddTVShowInLibrary(ctx context.Context, params tvshowlibrary.AddTVShowInLibraryParams) error
}

type TVShowDeliveryState interface {
	GetStateByID(ctx context.Context, stateID uuid.UUID) (*tvshowdeliverystate.State, error)
	Create(ctx context.Context, opt tvshowdeliverystate.CreateOptions) (*tvshowdeliverystate.State, error)
	Complete(ctx context.Context, stateID uuid.UUID, options ...any) (st *tvshowdeliverystate.State, executeErr error, err error)
}

type Labels interface {
	AddLabel(ctx context.Context, label labels.Label) error
}

// UUIDGenerator интерфейс для генерации UUID (реальный или мок)
type UUIDGenerator interface {
	New() uuid.UUID
}

// Clock интерфейс для работы со временем (реальный или мок)
type Clock interface {
	Now() time.Time
}

type uuidGenerator struct{}

func (uuidGenerator) New() uuid.UUID {
	return uuid.New()
}

type realClock struct{}

func (realClock) Now() time.Time {
	return time.Now()
}
