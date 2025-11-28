package tvshowdelete

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/kkiling/media-delivery/internal/adapter/apierr"
)

func (s *Service) DeleteSeasonFromMediaServer(ctx context.Context, tvShowPath TVShowCatalogPath) error {
	path, err := filepath.Rel(s.config.BasePath, tvShowPath.TVShowPath)
	if err != nil {
		return fmt.Errorf("failed to get relative path: %w", err)
	}

	if err = s.embyApi.Refresh(); err != nil {
		return fmt.Errorf("failed to refresh emby api: %w", err)
	}

	info, err := s.embyApi.GetCatalogInfo("/" + path)
	if err != nil {
		if errors.Is(err, apierr.ContentNotFound) {
			return nil
		}
		return fmt.Errorf("embyApi.GetCatalogInfo: %w", err)
	}

	if info == nil {
		return nil
	}

	return fmt.Errorf("tvshow is not deleted")
}

// isDirEmpty проверяет, пуста ли директория
func isDirEmpty(dirPath string) (bool, error) {
	f, err := os.Open(dirPath)
	if err != nil {
		return false, err
	}
	defer f.Close()

	// Читаем первые несколько записей
	_, err = f.Readdirnames(1)
	if err == io.EOF {
		return true, nil // Директория пуста
	}
	return false, err // Либо есть файлы, либо ошибка
}

func (s *Service) DeleteSeasonFiles(ctx context.Context, tvShowPath TVShowCatalogPath) error {
	seasonPath := filepath.Join(tvShowPath.TVShowPath, tvShowPath.SeasonPath)

	// Проверяем, существует ли путь
	if _, err := os.Stat(seasonPath); os.IsNotExist(err) {
		return fmt.Errorf("not found path")
	}

	// Удаляем папку сезона со всем содержимым
	err := os.RemoveAll(seasonPath)
	if err != nil {
		return fmt.Errorf("failed to delete folder: %w", err)
	}

	// Проверяем, остались ли файлы/папки в директории шоу
	isEmpty, err := isDirEmpty(tvShowPath.TVShowPath)
	if err != nil {
		// Если не можем проверить директорию, просто возвращаем успех
		// так как основная задача (удаление сезона) выполнена
		return nil
	}

	// Если директория шоу пуста, удаляем и её
	if isEmpty {
		err = os.Remove(tvShowPath.TVShowPath)
		if err != nil {
			return fmt.Errorf("failed to delete folder: %w", err)
		}
	}

	return nil
}
