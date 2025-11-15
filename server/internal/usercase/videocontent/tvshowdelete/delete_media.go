package tvshowdelete

import (
	"context"

	"github.com/kkiling/media-delivery/internal/common"
)

func (s *Service) DeleteSeasonFromMediaServer(ctx context.Context, tvShowID common.TVShowID) error {
	panic("implement me")
}

func (s *Service) DeleteSeasonFiles(ctx context.Context, tvShowPath TVShowCatalogPath) error {
	panic("implement me")
}
