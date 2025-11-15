package postgresql

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/samber/lo"

	"github.com/kkiling/media-delivery/internal/common"
	"github.com/kkiling/media-delivery/internal/usercase/videocontent/content"
	"github.com/kkiling/media-delivery/internal/usercase/videocontent/content/storage/db"
)

func (s *Storage) SaveVideoContent(ctx context.Context, videoContent *content.VideoContent) error {
	queries := s.getQueries(ctx)

	saveParams := db.SaveVideoContentParams{
		ID:             videoContent.ID,
		CreatedAt:      videoContent.CreatedAt,
		DeliveryStatus: int(videoContent.DeliveryStatus),
		States:         nil,
	}

	if videoContent.ContentID.MovieID != nil {
		saveParams.MovieID = lo.ToPtr(int64(*videoContent.ContentID.MovieID))
	}
	if videoContent.ContentID.TVShow != nil {
		saveParams.TvshowID = lo.ToPtr(int64(videoContent.ContentID.TVShow.ID))
		saveParams.SeasonNumber = lo.ToPtr(int32(videoContent.ContentID.TVShow.SeasonNumber))
	}

	states, err := json.Marshal(videoContent.States)
	if err != nil {
		return fmt.Errorf("failed to marshal States: %w", err)
	}
	saveParams.States = states

	err = queries.SaveVideoContent(ctx, saveParams)

	if err != nil {
		return s.base.HandleError(err)
	}

	return nil
}

func (s *Storage) GetVideoContents(ctx context.Context, contentID common.ContentID) ([]content.VideoContent, error) {
	queries := s.getQueries(ctx)

	if contentID.MovieID != nil {
		id := int64(*contentID.MovieID)
		res, err := queries.GetVideoContentsMovieID(ctx, &id)
		if err != nil {
			return nil, s.base.HandleError(err)
		}

		results := make([]content.VideoContent, 0, len(res))
		for _, item := range res {
			var state []content.State
			if err = json.Unmarshal(item.States, &state); err != nil {
				return nil, fmt.Errorf("failed to unmarshal States: %w", err)
			}
			results = append(results, content.VideoContent{
				ID:             item.ID,
				ContentID:      contentID,
				CreatedAt:      item.CreatedAt,
				DeliveryStatus: content.DeliveryStatus(item.DeliveryStatus),
				States:         state,
			})
		}

		return results, nil
	} else if contentID.TVShow != nil {
		res, err := queries.GetVideoContentTVShow(ctx, db.GetVideoContentTVShowParams{
			TvshowID:     lo.ToPtr(int64(contentID.TVShow.ID)),
			SeasonNumber: lo.ToPtr(int32(contentID.TVShow.SeasonNumber)),
		})
		if err != nil {
			return nil, s.base.HandleError(err)
		}

		results := make([]content.VideoContent, 0, len(res))
		for _, item := range res {
			var state []content.State
			if err = json.Unmarshal(item.States, &state); err != nil {
				return nil, fmt.Errorf("failed to unmarshal States: %w", err)
			}
			results = append(results, content.VideoContent{
				ID:             item.ID,
				ContentID:      contentID,
				CreatedAt:      item.CreatedAt,
				DeliveryStatus: content.DeliveryStatus(item.DeliveryStatus),
				States:         state,
			})
		}
		return results, nil
	}

	return nil, fmt.Errorf("contentID is not valid")
}

func (s *Storage) UpdateVideoContent(ctx context.Context, id uuid.UUID, videoContent *content.UpdateVideoContent) error {
	queries := s.getQueries(ctx)

	stateData, err := json.Marshal(videoContent.States)
	if err != nil {
		return fmt.Errorf("failed to marshal States: %w", err)
	}

	_, err = queries.UpdateVideoContent(ctx, db.UpdateVideoContentParams{
		DeliveryStatus: int(videoContent.DeliveryStatus),
		ID:             id,
		States:         stateData,
	})

	if err != nil {
		return s.base.HandleError(err)
	}

	return nil
}

func (s *Storage) GetVideoContentsByDeliveryStatus(ctx context.Context, deliveryStatus content.DeliveryStatus, limit int) ([]content.VideoContent, error) {
	queries := s.getQueries(ctx)

	res, err := queries.GetVideoContentsByDeliveryStatus(ctx, db.GetVideoContentsByDeliveryStatusParams{
		DeliveryStatus: int(deliveryStatus),
		Limit:          int32(limit),
	})
	if err != nil {
		return nil, s.base.HandleError(err)
	}

	results := make([]content.VideoContent, 0, len(res))
	for _, item := range res {
		var state []content.State
		if err = json.Unmarshal(item.States, &state); err != nil {
			return nil, fmt.Errorf("failed to unmarshal States: %w", err)
		}

		var contentID common.ContentID
		if item.MovieID != nil {
			contentID.MovieID = lo.ToPtr(uint64(*item.MovieID))
		} else if item.TvshowID != nil && item.SeasonNumber != nil {
			contentID.TVShow = &common.TVShowID{
				ID:           uint64(*item.TvshowID),
				SeasonNumber: uint8(*item.SeasonNumber),
			}
		} else {
			return nil, s.base.HandleError(fmt.Errorf("invalid contentID"))
		}

		results = append(results, content.VideoContent{
			ID:             item.ID,
			ContentID:      contentID,
			CreatedAt:      item.CreatedAt,
			DeliveryStatus: content.DeliveryStatus(item.DeliveryStatus),
			States:         state,
		})
	}

	return results, nil
}
