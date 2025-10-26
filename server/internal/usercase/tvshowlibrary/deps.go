package tvshowlibrary

import (
	"context"

	"github.com/kkiling/media-delivery/internal/adapter/themoviedb"
)

type TheMovieDb interface {
	SearchTV(ctx context.Context, params themoviedb.SearchQuery) (*themoviedb.TVShowSearchResponse, error)
	GetTV(ctx context.Context, tvID uint64, language themoviedb.Language) (*themoviedb.TVShow, error)
	GetSeasonInfo(ctx context.Context, tvID uint64, seasonNumber uint8, language themoviedb.Language) (*themoviedb.SeasonWithEpisodes, error)
}

type Storage interface {
	SaveOrUpdateTVShow(ctx context.Context, tvShow *TVShow) error
	GetTVShows(ctx context.Context) ([]TVShowShort, error)
	SaveOrUpdateSeasonEpisode(ctx context.Context, tvID uint64, seasonNumber uint8, episodes []Episode) error
}
