package tvshowdelete

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/kkiling/media-delivery/internal/common"
)

func (s *Service) DeleteSeasonFromMediaServer(ctx context.Context, tvShowPath TVShowCatalogPath, tvShowID common.TVShowID) error {
	return nil
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
