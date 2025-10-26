package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"time"

	"github.com/kkiling/goplatform/storagebase"
	"github.com/kkiling/media-delivery/internal/common"
	"github.com/kkiling/media-delivery/internal/usercase/labels"
	"github.com/samber/lo"
)

func (s *Storage) SaveLabel(ctx context.Context, label labels.Label) error {
	var movieID *string
	var tvshowID *string
	var seasonNumber *uint8

	if label.ContentID.MovieID != nil {
		movieID = lo.ToPtr(strconv.FormatUint(*label.ContentID.MovieID, 10))
	}
	if label.ContentID.TVShow != nil {
		tvshowID = lo.ToPtr(strconv.FormatUint(label.ContentID.TVShow.ID, 10))
		seasonNumber = &label.ContentID.TVShow.SeasonNumber
	}

	query := `
        INSERT INTO content_label (
            created_at,    
            type_label,
            movie_id,
            tvshow_id,
            season_number
        ) VALUES (?, ?, ?, ?, ?)
    `

	// Генерируем UUID для новой записи

	_, err := s.base.Next(ctx).ExecContext(ctx, query,
		label.CreatedAt.Format(time.RFC3339),
		label.TypeLabel,
		movieID,
		tvshowID,
		seasonNumber,
	)
	if err != nil {
		return s.base.HandleError(err)
	}

	return nil
}

func (s *Storage) getLabels(rows *sql.Rows) ([]labels.Label, error) {
	var results []labels.Label
	for rows.Next() {
		var (
			vc           labels.Label
			movieID      sql.NullInt64
			tvshowID     sql.NullInt64
			seasonNumber sql.NullInt16
		)

		err := rows.Scan(
			&vc.CreatedAt,
			&movieID,
			&tvshowID,
			&seasonNumber,
			&vc.TypeLabel,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan content_label: %w", err)
		}

		// rebuild ContentID
		var cid common.ContentID
		if movieID.Valid {
			mid := uint64(movieID.Int64)
			cid.MovieID = &mid
		}
		if tvshowID.Valid {
			tid := uint64(tvshowID.Int64)
			sn := uint8(0)
			if seasonNumber.Valid {
				sn = uint8(seasonNumber.Int16)
			}
			cid.TVShow = &common.TVShowID{ID: tid, SeasonNumber: sn}
		}
		vc.ContentID = cid

		results = append(results, vc)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return results, nil
}

func (s *Storage) GetLabels(ctx context.Context, contentID common.ContentID) ([]labels.Label, error) {
	next := s.base.Next(ctx)
	var rows *sql.Rows
	var err error

	if contentID.MovieID != nil {
		rows, err = next.QueryContext(ctx, `
            SELECT 
                movie_id,
                tvshow_id,
                season_number,
                created_at,
                type_label
            FROM content_label
            WHERE movie_id = ?
        `, *contentID.MovieID)
	} else if contentID.TVShow != nil {
		rows, err = next.QueryContext(ctx, `
            SELECT 
                id,
                created_at,
                movie_id,
                tvshow_id,
                season_number,
                type_label
            FROM content_label
            WHERE tvshow_id = ? AND season_number = ?
        `, contentID.TVShow.ID, contentID.TVShow.SeasonNumber)
	} else {
		return nil, fmt.Errorf("invalid contentID: both MovieID and TVShow are nil")
	}

	if err != nil {
		return nil, s.base.HandleError(err)
	}
	defer rows.Close()

	return s.getLabels(rows)
}

func (s *Storage) DeleteLabel(ctx context.Context, contentID common.ContentID, typeLabel labels.TypeLabel) error {
	return storagebase.ErrNotFound
}
