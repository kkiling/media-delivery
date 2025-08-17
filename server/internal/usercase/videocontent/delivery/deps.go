package delivery

import (
	"context"

	"github.com/google/uuid"

	"github.com/kkiling/media-delivery/internal/adapter/emby"
	"github.com/kkiling/media-delivery/internal/adapter/matchtvshow"
	"github.com/kkiling/media-delivery/internal/adapter/mkvmerge"
	"github.com/kkiling/media-delivery/internal/adapter/qbittorrent"
	"github.com/kkiling/media-delivery/internal/adapter/rutracker"
	"github.com/kkiling/media-delivery/internal/usercase/tvshowlibrary"
)

type TVShowLibrary interface {
	GetTVShowInfo(ctx context.Context, params tvshowlibrary.GetTVShowParams) (*tvshowlibrary.GetTVShowResult, error)
	GetSeasonEpisodes(ctx context.Context, params tvshowlibrary.GetSeasonEpisodesParams) (*tvshowlibrary.GetSeasonEpisodesResult, error)
}

type TorrentSite interface {
	SearchTorrents(query string) (*rutracker.TorrentResponse, error)
	GetMagnetLink(torrentUrl string) (*rutracker.MagnetInfo, error)
}

type TorrentClient interface {
	AddTorrent(opts qbittorrent.TorrentAddOptions) error
	GetTorrentInfo(hash string) (*qbittorrent.TorrentInfo, error)
	GetTorrentFiles(hash string) ([]qbittorrent.TorrentFile, error)
	ResumeTorrent(hash string) error
}

type PrepareTVShow interface {
	PrepareTvShowSeason(params *matchtvshow.PrepareTvShowPrams) (*matchtvshow.PrepareTVShowSeason, error)
}

type MkvMergePipeline interface {
	AddToMerge(ctx context.Context, idempotencyKey string, params mkvmerge.MergeParams) (*mkvmerge.MergeResult, error)
	GetMergeResult(ctx context.Context, id uuid.UUID) (*mkvmerge.MergeResult, error)
}

type EmbyApi interface {
	Refresh() error
	ResetMetadata(embyID uint64) error
	RemoteSearchApply(embyID, theMovieDBID uint64) error
	GetCatalogInfo(path string) (*emby.CatalogInfo, error)
}
