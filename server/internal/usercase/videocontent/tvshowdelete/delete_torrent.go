package tvshowdelete

import (
	"context"
	"fmt"
	"os"
)

func (s *Service) DeleteTorrentFromTorrentClient(ctx context.Context, magnetHash string) error {
	return s.torrentClient.DeleteTorrent(magnetHash, false)
}

func (s *Service) DeleteTorrentFiles(ctx context.Context, torrentPath string) error {
	// Проверяем, существует ли путь
	if _, err := os.Stat(torrentPath); os.IsNotExist(err) {
		return fmt.Errorf("not found path")
	}

	// Удаляем каталог со всем содержимым
	err := os.RemoveAll(torrentPath)
	if err != nil {
		return fmt.Errorf("failed to delete folder: %w", err)
	}

	return nil
}
