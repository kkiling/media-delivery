package delivery

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"syscall"

	ucerr "github.com/kkiling/torrent-to-media-server/internal/usercase/err"
)

type CreateContentCatalogsParams struct {
	TVShowCatalogPath TVShowCatalogPath
}

func (s *Service) createDirectories(seasonPath string) error {
	// Проверяем, что catalog действительно является подкаталогом base
	relPath, err := filepath.Rel(s.config.BasePath, seasonPath)
	if err != nil {
		return fmt.Errorf("catalog is not a subdirectory of base: %v", err)
	}

	// Разбиваем относительный путь на компоненты
	parts := strings.Split(relPath, string(filepath.Separator))

	// Постепенно создаём каталоги
	currentPath := s.config.BasePath
	for _, part := range parts {
		currentPath = filepath.Join(currentPath, part)
		// Проверяем существование каталога
		if _, err = os.Stat(currentPath); os.IsNotExist(err) {
			// Создаем каталог

			if mkdirErr := syscall.Mkdir(currentPath, 0775); mkdirErr != nil {
				return fmt.Errorf("syscall.Mkdir: %w", mkdirErr)
			}
			if err := os.Chmod(currentPath, 0775); err != nil {
				return err
			}
			// Меняем группу пользователей
			if s.config.UserGroup != "" {
				if err = setGroup(currentPath, s.config.UserGroup); err != nil {
					return fmt.Errorf("syscall.Chown: %w", err)
				}
			}
		} else if err != nil {
			return fmt.Errorf("error checking directory %s: %v", currentPath, err)
		}
	}

	return nil
}

func isEmpty(dirPath string) (bool, error) {
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return false, err
	}
	return len(entries) == 0, nil
}

// CreateContentCatalogs формирование каталога куда будет сохранен контент
func (s *Service) CreateContentCatalogs(ctx context.Context, params CreateContentCatalogsParams) error {
	seasonPath := params.TVShowCatalogPath.FullSeasonPath()

	if createErr := s.createDirectories(seasonPath); createErr != nil {
		return fmt.Errorf("createDirectories: %w", createErr)
	}

	// Если каталог сезона не пустой, то выдаем ошибку, что бы пользователь сам устранил ошибку
	if ok, err := isEmpty(seasonPath); err != nil {
		return fmt.Errorf("isEmpty: %w", err)
	} else if !ok {
		return fmt.Errorf("catalog is not empty: %w", ucerr.AlreadyExists)
	}

	return nil
}
