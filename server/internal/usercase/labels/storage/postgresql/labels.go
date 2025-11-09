package postgresql

import (
	"context"
	"fmt"

	"github.com/samber/lo"

	"github.com/kkiling/media-delivery/internal/common"
	"github.com/kkiling/media-delivery/internal/usercase/labels"
	"github.com/kkiling/media-delivery/internal/usercase/labels/storage/db"
)

func (s *Storage) SaveLabel(ctx context.Context, label labels.Label) error {
	queries := s.getQueries(ctx)

	saveParams := db.SaveLabelParams{
		CreatedAt: label.CreatedAt,
		TypeLabel: int(label.TypeLabel),
	}

	if label.ContentID.MovieID != nil {
		saveParams.MovieID = lo.ToPtr(int64(*label.ContentID.MovieID))
	}
	if label.ContentID.TVShow != nil {
		saveParams.TvshowID = lo.ToPtr(int64(label.ContentID.TVShow.ID))
		saveParams.SeasonNumber = lo.ToPtr(int32(label.ContentID.TVShow.SeasonNumber))
	}

	err := queries.SaveLabel(ctx, saveParams)

	if err != nil {
		return s.base.HandleError(err)
	}

	return nil
}

func (s *Storage) GetLabels(ctx context.Context, contentID common.ContentID) ([]labels.Label, error) {
	queries := s.getQueries(ctx)

	if contentID.MovieID != nil {
		id := int64(*contentID.MovieID)
		res, err := queries.GetLabelsMovieID(ctx, &id)
		if err != nil {
			return nil, s.base.HandleError(err)
		}
		return lo.Map(res, func(item db.GetLabelsMovieIDRow, index int) labels.Label {
			return labels.Label{
				ContentID: contentID,
				TypeLabel: labels.TypeLabel(item.TypeLabel),
				CreatedAt: item.CreatedAt,
			}
		}), nil
	} else if contentID.TVShow != nil {
		res, err := queries.GetLabelsTVShow(ctx, db.GetLabelsTVShowParams{
			TvshowID:     lo.ToPtr(int64(contentID.TVShow.ID)),
			SeasonNumber: lo.ToPtr(int32(contentID.TVShow.SeasonNumber)),
		})
		if err != nil {
			return nil, s.base.HandleError(err)
		}

		return lo.Map(res, func(item db.GetLabelsTVShowRow, index int) labels.Label {
			return labels.Label{
				ContentID: contentID,
				TypeLabel: labels.TypeLabel(item.TypeLabel),
				CreatedAt: item.CreatedAt,
			}
		}), nil
	}

	return nil, fmt.Errorf("contentID is not valid")
}

func (s *Storage) DeleteLabel(ctx context.Context, contentID common.ContentID, typeLabel labels.TypeLabel) error {
	queries := s.getQueries(ctx)

	if contentID.MovieID != nil {
		id := int64(*contentID.MovieID)
		_, err := queries.DeleteLabelMovieID(ctx, &id)
		if err != nil {
			return s.base.HandleError(err)
		}
		return nil
	} else if contentID.TVShow != nil {
		_, err := queries.DeleteLabelTVShow(ctx, db.DeleteLabelTVShowParams{
			TvshowID:     lo.ToPtr(int64(contentID.TVShow.ID)),
			SeasonNumber: lo.ToPtr(int32(contentID.TVShow.SeasonNumber)),
		})
		if err != nil {
			return s.base.HandleError(err)
		}
		return nil
	}

	return fmt.Errorf("contentID is not valid")
}
