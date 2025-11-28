package tvshowdeletestate

import (
	"context"

	"github.com/kkiling/media-delivery/internal/common"
	"github.com/kkiling/media-delivery/internal/usercase/videocontent/tvshowdelete"
)

type ContentDeleted interface {
	DeleteSeasonFromMediaServer(ctx context.Context, tvShowPath tvshowdelete.TVShowCatalogPath, tvShowID common.TVShowID) error
	DeleteSeasonFiles(ctx context.Context, tvShowPath tvshowdelete.TVShowCatalogPath) error
	DeleteTorrentFiles(ctx context.Context, torrentPath string) error
	DeleteTorrentFromTorrentClient(ctx context.Context, magnetHash string) error
	DeleteLabelHasVideoContentFiles(ctx context.Context, contentID common.ContentID) error
}
