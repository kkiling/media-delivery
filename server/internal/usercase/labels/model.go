package labels

import (
	"time"

	"github.com/kkiling/media-delivery/internal/common"
)

type TypeLabel int

const (
	// ContentInLibrary Сериал или фильм добавлен в библиотеку,
	// UpdatedAt дата последнего обновления информации
	ContentInLibrary TypeLabel = iota
	// HasVideoContentFiles у Сериала или фильма есть видео контент (video_content) то есть имеется доставка файлов
	// UpdatedAt дата создания сущности video_content (доставки)
	HasVideoContentFiles TypeLabel = iota
)

type Label struct {
	ContentID common.ContentID
	TypeLabel TypeLabel
	CreatedAt time.Time
}
