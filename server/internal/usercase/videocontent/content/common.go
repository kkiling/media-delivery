package content

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/kkiling/media-delivery/internal/common"
	ucerr "github.com/kkiling/media-delivery/internal/usercase/err"
	"github.com/kkiling/media-delivery/internal/usercase/videocontent/runners"
	"github.com/samber/lo"
)

func (s *Service) getStateID(ctx context.Context, contentID common.ContentID, runersType runners.Type) (uuid.UUID, error) {
	contents, err := s.storage.GetVideoContents(ctx, contentID)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("storage.GetVideoContent: %w", err)
	}

	if len(contents) != 1 {
		return uuid.UUID{}, ucerr.NotFound
	}

	content := contents[0]
	state, find := lo.Find(content.States, func(item State) bool {
		return item.Type == runersType
	})
	if !find {
		return uuid.UUID{}, ucerr.NotFound
	}

	return state.StateID, nil
}
