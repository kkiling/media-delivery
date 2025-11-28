package tvshowdelete

import (
	"context"

	"github.com/kkiling/media-delivery/internal/adapter/emby"
	"github.com/kkiling/media-delivery/internal/common"
	"github.com/kkiling/media-delivery/internal/usercase/labels"
)

type TorrentClient interface {
	DeleteTorrent(hash string, deleteFiles bool) error
}

type EmbyApi interface {
	Refresh() error
	GetCatalogInfo(path string) (*emby.CatalogInfo, error)
}

type Labels interface {
	DeleteLabel(ctx context.Context, contentID common.ContentID, typeLabel labels.TypeLabel) error
}
