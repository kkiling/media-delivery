package content

import (
	"context"
	"fmt"
	"sort"

	"github.com/google/uuid"
	"github.com/samber/lo"

	"github.com/kkiling/media-delivery/internal/common"
	ucerr "github.com/kkiling/media-delivery/internal/usercase/err"
	"github.com/kkiling/media-delivery/internal/usercase/videocontent/runners"
)

func getLastState(content VideoContent, runersType runners.Type) *uuid.UUID {
	sort.Slice(content.States, func(i, j int) bool {
		return content.States[i].CreatedAt.After(content.States[j].CreatedAt)
	})
	res, find := lo.Find(content.States, func(item State) bool {
		return item.Type == runersType
	})
	if !find {
		return nil
	}

	return &res.StateID
}

func (s *Service) getVideoContent(ctx context.Context, contentID common.ContentID) (VideoContent, error) {
	contents, err := s.storage.GetVideoContents(ctx, contentID)
	if err != nil {
		return VideoContent{}, fmt.Errorf("storage.GetVideoContent: %w", err)
	}

	if len(contents) != 1 {
		return VideoContent{}, ucerr.NotFound
	}

	return contents[0], nil
}
