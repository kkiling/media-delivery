package tvshowlibrary

import (
	"context"
	"time"

	"github.com/kkiling/media-delivery/internal/adapter/themoviedb"
)

type TheMovieDb interface {
	SearchTV(ctx context.Context, params themoviedb.SearchQuery) (*themoviedb.TVShowSearchResponse, error)
	GetTV(ctx context.Context, tvID uint64, language themoviedb.Language) (*themoviedb.TVShow, error)
	GetSeason(ctx context.Context, tvID uint64, seasonNumber uint8, language themoviedb.Language) (*themoviedb.SeasonWithEpisodes, error)
}

type Storage interface {
	SaveTVShow(ctx context.Context, tvShow *TVShow) error
	GetTVShows(ctx context.Context) ([]TVShowShort, error)
	SaveEpisodes(ctx context.Context, tvID uint64, seasonNumber uint8, episodes []Episode) error
}

// Clock интерфейс для работы со временем (реальный или мок)
type Clock interface {
	Now() time.Time
}
