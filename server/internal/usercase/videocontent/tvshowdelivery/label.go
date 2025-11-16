package tvshowdelivery

import (
	"context"
	"fmt"
	"time"

	"github.com/kkiling/media-delivery/internal/common"
	"github.com/kkiling/media-delivery/internal/usercase/labels"
)

func (s *Service) AddLabelHasVideoContentFiles(ctx context.Context, contentID common.ContentID) error {
	err := s.labels.AddLabel(ctx, labels.Label{
		ContentID: contentID,
		TypeLabel: labels.HasVideoContentFiles,
		CreatedAt: time.Now(),
	})
	if err != nil {
		return fmt.Errorf("labels.AddLabe: %w", err)
	}

	return nil
}
