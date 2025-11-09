package postgresql

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/kkiling/media-delivery/internal/adapter/mkvmerge"
	"github.com/kkiling/media-delivery/internal/adapter/mkvmerge/storage/db"
)

func (s *Storage) Create(ctx context.Context, create *mkvmerge.CreateMergeResult) error {
	queries := s.getQueries(ctx)

	params, err := json.Marshal(create.Params)
	if err != nil {
		return fmt.Errorf("failed to marshal: %w", err)
	}

	err = queries.CreateMkvMerge(ctx, db.CreateMkvMergeParams{
		ID:             create.ID,
		IdempotencyKey: create.IdempotencyKey,
		Params:         params,
		Status:         int(create.Status),
		CreatedAt:      create.CreatedAt,
	})

	if err != nil {
		return s.base.HandleError(err)
	}

	return nil
}

func (s *Storage) Update(ctx context.Context, id uuid.UUID, update *mkvmerge.UpdateMergeResult) error {
	queries := s.getQueries(ctx)

	_, err := queries.UpdateMkvMerge(ctx, db.UpdateMkvMergeParams{
		Status: int(update.Status),
		Error:  update.Error,
		CompletedAt: pgtype.Timestamptz{
			Time: func() time.Time {
				if update.Completed != nil {
					return *update.Completed
				}
				return time.Time{}
			}(),
			Valid: update.Completed != nil,
		},
		ID: id,
	})

	if err != nil {
		return s.base.HandleError(err)
	}

	return nil
}

func (s *Storage) UpdateProgress(ctx context.Context, id uuid.UUID, progress float32) error {
	queries := s.getQueries(ctx)
	_, err := queries.UpdateProgress(ctx, db.UpdateProgressParams{
		Progress: &progress,
		ID:       id,
	})
	if err != nil {
		return s.base.HandleError(err)
	}
	return nil
}

func mapMkvMerge(res db.MkvMerge) (*mkvmerge.MergeResult, error) {
	var params mkvmerge.MergeParams
	if err := json.Unmarshal(res.Params, &params); err != nil {
		return nil, fmt.Errorf("failed to unmarshal: %w", err)
	}

	return &mkvmerge.MergeResult{
		ID:             res.ID,
		IdempotencyKey: res.IdempotencyKey,
		Params:         params,
		Status:         mkvmerge.Status(res.Status),
		Error:          res.Error,
		CreatedAt:      res.CreatedAt,
		CompletedAt: func() *time.Time {
			if res.CompletedAt.Valid {
				return &res.CompletedAt.Time
			}
			return nil
		}(),
		Progress: res.Progress,
	}, nil
}

func (s *Storage) GetByID(ctx context.Context, id uuid.UUID) (*mkvmerge.MergeResult, error) {
	queries := s.getQueries(ctx)
	res, err := queries.GetByID(ctx, id)
	if err != nil {
		return nil, s.base.HandleError(err)
	}
	return mapMkvMerge(res)
}

func (s *Storage) GetByIdempotencyKey(ctx context.Context, idempotencyKey string) (*mkvmerge.MergeResult, error) {
	queries := s.getQueries(ctx)
	res, err := queries.GetByIdempotencyKey(ctx, idempotencyKey)
	if err != nil {
		return nil, s.base.HandleError(err)
	}

	return mapMkvMerge(res)
}

func (s *Storage) GetOldestUncompleted(ctx context.Context) (*mkvmerge.MergeResult, error) {
	queries := s.getQueries(ctx)
	res, err := queries.GetOldestUncompleted(ctx)
	if err != nil {
		return nil, s.base.HandleError(err)
	}

	return mapMkvMerge(res)
}

func (s *Storage) AddMergeLogs(ctx context.Context, id uuid.UUID, log mkvmerge.MergeLogs) error {
	queries := s.getQueries(ctx)

	err := queries.AddMergeLogs(ctx, db.AddMergeLogsParams{
		MergeID:   id,
		CreatedAt: log.CreatedAt,
		Type:      int(log.Type),
		Content:   log.Content,
	})

	if err != nil {
		return s.base.HandleError(err)
	}

	return nil
}

func (s *Storage) DeleteLogs(ctx context.Context, mergeID uuid.UUID) error {
	queries := s.getQueries(ctx)

	err := queries.DeleteLogs(ctx, mergeID)
	if err != nil {
		return s.base.HandleError(err)
	}

	return nil
}
