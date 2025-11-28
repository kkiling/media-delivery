package tvshowdelete

import (
	"context"
	"fmt"

	"github.com/kkiling/media-delivery/internal/common"
	"github.com/kkiling/media-delivery/internal/usercase/labels"
)

func (s *Service) DeleteLabelHasVideoContentFiles(ctx context.Context, contentID common.ContentID) error {
	err := s.labels.DeleteLabel(ctx, contentID, labels.HasVideoContentFiles)
	if err != nil {
		return fmt.Errorf("labels.AddLabe: %w", err)
	}

	return nil
}
