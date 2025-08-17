package delivery

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/samber/lo"

	"github.com/kkiling/media-delivery/internal/adapter/mkvmerge"
)

type MergeVideoParams struct {
	IdempotencyKey string
	ContentMatches []ContentMatches
}

type MergeVideoStatus struct {
	Progress   float64 // 0 до 1
	IsComplete bool
	Errors     []string
}

func mapMkvMergeParams(content ContentMatches) mkvmerge.MergeParams {
	mergeParams := mkvmerge.MergeParams{
		VideoInputFile:  content.Video.File.FullPath,
		VideoOutputFile: content.Episode.FileName,
		AudioTracks: lo.Map(content.AudioFiles, func(item Track, index int) mkvmerge.Track {
			return mkvmerge.Track{
				Path:     item.File.FullPath,
				Language: item.Language,
				Name:     item.Name,
				Default:  index == 0,
			}
		}),
		SubtitleTracks: lo.Map(content.Subtitles, func(item Track, index int) mkvmerge.Track {
			return mkvmerge.Track{
				Path:     item.File.FullPath,
				Language: item.Language,
				Name:     item.Name,
				Default:  false,
			}
		}),
	}
	return mergeParams
}

// StartMergeVideo запуск обработки видеофайлов
func (s *Service) StartMergeVideo(ctx context.Context, params MergeVideoParams) ([]uuid.UUID, error) {
	result := make([]uuid.UUID, 0, len(params.ContentMatches))

	for _, content := range params.ContentMatches {
		mergeParams := mapMkvMergeParams(content)
		idempotencyKey := fmt.Sprintf("%s-episode_%d", params.IdempotencyKey, content.Episode.EpisodeNumber)

		mergeResult, err := s.mkvMerge.AddToMerge(ctx, idempotencyKey, mergeParams)
		if err != nil {
			return nil, fmt.Errorf("mkvMerge.Merge: %w", err)
		}
		result = append(result, mergeResult.ID)
	}
	return result, nil
}

func (s *Service) GetMergeVideoStatus(ctx context.Context, mergeIDs []uuid.UUID) (*MergeVideoStatus, error) {
	var status MergeVideoStatus

	delta := 1.0 / float64(len(mergeIDs))
	status.Progress = 0.0
	status.IsComplete = true

	for _, id := range mergeIDs {
		result, err := s.mkvMerge.GetMergeResult(ctx, id)
		if err != nil {
			return nil, fmt.Errorf("mkvMerge.GetMergeResult: %w", err)
		}

		if result.Status == mkvmerge.ErrorStatus || result.Status == mkvmerge.CompleteStatus {
			if result.Error != nil {
				status.Errors = append(status.Errors, *result.Error)
			}
			status.Progress += delta
			continue
		}

		status.IsComplete = false
		if result.Progress != nil {
			status.Progress += delta * *result.Progress
		}
	}

	return &status, nil
}
