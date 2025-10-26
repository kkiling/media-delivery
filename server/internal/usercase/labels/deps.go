package labels

import (
	"context"

	"github.com/kkiling/media-delivery/internal/common"
)

type Storage interface {
	SaveLabel(ctx context.Context, label Label) error
	GetLabels(ctx context.Context, contentID common.ContentID) ([]Label, error)
	DeleteLabel(ctx context.Context, contentID common.ContentID, typeLabel TypeLabel) error
}
