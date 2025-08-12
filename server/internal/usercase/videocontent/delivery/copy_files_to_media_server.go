package delivery

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
)

type CreateSymLinkCopyParams struct {
	SeasonPath     string
	ContentMatches []ContentMatches
}

// CreateSymLinkCopyToMediaServer шаг копирования файлов на медиа сервер
func (s *Service) CreateSymLinkCopyToMediaServer(_ context.Context, params CreateSymLinkCopyParams) error {
	for _, match := range params.ContentMatches {
		from := match.Video.File.FullPath
		to := filepath.Join(params.SeasonPath, match.Episode.EpisodeFileName)
		// Создание hardlink ссылки (to) на файл from
		// Create hard link
		err := os.Symlink(from, to)
		if err != nil {
			// Handle error (file might already exist, permissions issue, etc.)
			return fmt.Errorf("failed to create hard link from %s to %s: %w", from, to, err)
		}
	}

	return nil
}
