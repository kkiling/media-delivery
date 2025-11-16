package tvshowdelete

import (
	"context"

	"github.com/kkiling/media-delivery/internal/adapter/emby"
	"github.com/kkiling/media-delivery/internal/adapter/qbittorrent"
	"github.com/kkiling/media-delivery/internal/common"
	"github.com/kkiling/media-delivery/internal/usercase/labels"
)

type TorrentClient interface {
	AddTorrent(opts qbittorrent.TorrentAddOptions) error
	GetTorrentInfo(hash string) (*qbittorrent.TorrentInfo, error)
	GetTorrentFiles(hash string) ([]qbittorrent.TorrentFile, error)
	ResumeTorrent(hash string) error
}

type EmbyApi interface {
	Refresh() error
	ResetMetadata(embyID uint64) error
	RemoteSearchApply(embyID, theMovieDBID uint64) error
	GetCatalogInfo(path string) (*emby.CatalogInfo, error)
}

type Labels interface {
	DeleteLabel(ctx context.Context, contentID common.ContentID, typeLabel labels.TypeLabel) error
}
