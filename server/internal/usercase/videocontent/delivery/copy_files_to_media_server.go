package delivery

import (
	"context"
	"fmt"
	"os"
)

type CreateHardLinkCopyParams struct {
	ContentMatches []ContentMatches
}

// CreateHardLinkCopyToMediaServer шаг копирования файлов на медиа сервер
func (s *Service) CreateHardLinkCopyToMediaServer(_ context.Context, params CreateHardLinkCopyParams) error {
	for _, match := range params.ContentMatches {
		from := match.Video.File.FullPath
		to := match.Episode.FileName
		if from == "" || to == "" {
			continue
		}
		// Создание hardlink ссылки (to) на файл from
		// Create hard link
		err := os.Link(from, to)
		if err != nil {
			// Handle error (file might already exist, permissions issue, etc.)
			return fmt.Errorf("failed to create hard link from %s to %s: %w", from, to, err)
		}
	}

	return nil
}
