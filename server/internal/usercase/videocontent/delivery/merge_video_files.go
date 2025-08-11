package delivery

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/samber/lo"

	"github.com/kkiling/torrent-to-media-server/internal/adapter/mkvmerge"
)

type MergeVideoParams struct {
	ContentPath    string
	IdempotencyKey string
	ContentMatches []ContentMatches
}

type MergeVideoFile struct {
	MergeID         uuid.UUID
	VideoInputFile  string
	VideoOutputFile string
}

type MergeVideoStatus struct {
	Progress   float64 // 0 до 1
	IsComplete bool
	Errors     []string
}

func mapMkvMergeParams(content ContentMatches, contentPath string) mkvmerge.MergeParams {
	// Формируем выходное наименование эпизода
	episodeName := fmt.Sprintf("S%02dE%02d %s", content.Episode.SeasonNumber, content.Episode.EpisodeNumber, content.Episode.EpisodeName)

	mergeParams := mkvmerge.MergeParams{
		VideoInputFile:  content.Video.File.FullPath,
		VideoOutputFile: filepath.Join(contentPath, episodeName) + content.Video.File.Extension,
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
func (s *Service) StartMergeVideo(ctx context.Context, params MergeVideoParams) ([]MergeVideoFile, error) {
	result := make([]MergeVideoFile, 0, len(params.ContentMatches))
	for _, content := range params.ContentMatches {
		mergeParams := mapMkvMergeParams(content, params.ContentPath)
		idempotencyKey := fmt.Sprintf("%s-episode_%d", params.IdempotencyKey, content.Episode.EpisodeNumber)
		mergeResult, err := s.mkvMerge.AddToMerge(ctx, idempotencyKey, mergeParams)
		if err != nil {
			return nil, fmt.Errorf("mkvMerge.Merge: %w", err)
		}
		result = append(result, MergeVideoFile{
			MergeID:         mergeResult.ID,
			VideoInputFile:  mergeResult.Params.VideoInputFile,
			VideoOutputFile: mergeResult.Params.VideoOutputFile,
		})
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
