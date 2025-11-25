package tvshowdeliverystate

import (
	"context"

	"github.com/google/uuid"

	"github.com/kkiling/media-delivery/internal/common"
	"github.com/kkiling/media-delivery/internal/usercase/videocontent/tvshowdelivery"
)

type ContentDelivery interface {
	GenerateSearchQuery(ctx context.Context, params tvshowdelivery.GenerateSearchQueryParams) (*tvshowdelivery.SearchQuery, error)
	SearchTorrent(ctx context.Context, params tvshowdelivery.SearchTorrentParams) ([]tvshowdelivery.TorrentSearch, error)
	GetMagnetLink(ctx context.Context, params tvshowdelivery.GetMagnetLinkParams) (*tvshowdelivery.MagnetLink, error)
	AddTorrentToTorrentClient(ctx context.Context, params tvshowdelivery.AddTorrentParams) error
	WaitingTorrentFiles(ctx context.Context, params tvshowdelivery.WaitingTorrentFilesParams) (*tvshowdelivery.TorrentFilesData, error)
	PrepareFileMatches(ctx context.Context, params tvshowdelivery.PreparingFileMatchesParams) (*tvshowdelivery.ContentMatches, error)
	WaitingTorrentDownloadComplete(ctx context.Context, params tvshowdelivery.WaitingTorrentDownloadCompleteParams) (*tvshowdelivery.TorrentDownloadStatus, error)
	GetEpisodesData(ctx context.Context, params tvshowdelivery.GetEpisodesDataParams) (*tvshowdelivery.EpisodesData, error)
	CreateContentCatalogs(ctx context.Context, params tvshowdelivery.CreateContentCatalogsParams) error
	StartMergeVideo(ctx context.Context, params tvshowdelivery.MergeVideoParams) ([]uuid.UUID, error)
	GetMergeVideoStatus(ctx context.Context, mergeIDs []uuid.UUID) (*tvshowdelivery.MergeVideoStatus, error)
	GetCatalogSize(ctx context.Context, catalogPath string) (uint64, error)
	SetMediaMetaData(ctx context.Context, params tvshowdelivery.SetMediaMetaDataParams) error
	NeedPrepareFileMatches(contentMatches []tvshowdelivery.ContentMatch) bool
	CreateHardLinkCopyToMediaServer(ctx context.Context, params tvshowdelivery.CreateHardLinkCopyParams) error
	ValidateContentMatch(oldContentMatch *tvshowdelivery.ContentMatches, newContentMatch *tvshowdelivery.ContentMatches) error
	AddLabelHasVideoContentFiles(ctx context.Context, contentID common.ContentID) error
}
