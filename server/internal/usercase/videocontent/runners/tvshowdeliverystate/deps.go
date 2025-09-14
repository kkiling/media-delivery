package tvshowdeliverystate

import (
	"context"

	"github.com/google/uuid"

	"github.com/kkiling/media-delivery/internal/usercase/videocontent/delivery"
)

type ContentDelivery interface {
	GenerateSearchQuery(ctx context.Context, params delivery.GenerateSearchQueryParams) (*delivery.SearchQuery, error)
	SearchTorrent(ctx context.Context, params delivery.SearchTorrentParams) ([]delivery.TorrentSearch, error)
	GetMagnetLink(ctx context.Context, params delivery.GetMagnetLinkParams) (*delivery.MagnetLink, error)
	AddTorrentToTorrentClient(ctx context.Context, params delivery.AddTorrentParams) error
	WaitingTorrentFiles(ctx context.Context, params delivery.WaitingTorrentFilesParams) (*delivery.TorrentFilesData, error)
	PrepareFileMatches(ctx context.Context, params delivery.PreparingFileMatchesParams) (*delivery.ContentMatches, error)
	WaitingTorrentDownloadComplete(ctx context.Context, params delivery.WaitingTorrentDownloadCompleteParams) (*delivery.TorrentDownloadStatus, error)
	GetEpisodesData(ctx context.Context, params delivery.GetEpisodesDataParams) (*delivery.EpisodesData, error)
	CreateContentCatalogs(ctx context.Context, params delivery.CreateContentCatalogsParams) error
	StartMergeVideo(ctx context.Context, params delivery.MergeVideoParams) ([]uuid.UUID, error)
	GetMergeVideoStatus(ctx context.Context, mergeIDs []uuid.UUID) (*delivery.MergeVideoStatus, error)
	GetCatalogSize(ctx context.Context, catalogPath string) (uint64, error)
	SetMediaMetaData(ctx context.Context, params delivery.SetMediaMetaDataParams) error
	NeedPrepareFileMatches(contentMatches []delivery.ContentMatch) bool
	CreateHardLinkCopyToMediaServer(ctx context.Context, params delivery.CreateHardLinkCopyParams) error
	ValidateContentMatch(oldContentMatch *delivery.ContentMatches, newContentMatch *delivery.ContentMatches) error
}
