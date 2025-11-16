package content

import (
	"time"

	"github.com/google/uuid"

	"github.com/kkiling/media-delivery/internal/common"
	"github.com/kkiling/media-delivery/internal/usercase/videocontent/runners"
)

type DeliveryStatus int

const (
	// DeliveryStatusFailed доставка была зафейлена
	DeliveryStatusFailed DeliveryStatus = iota
	// DeliveryStatusNew - Новая доставка
	DeliveryStatusNew DeliveryStatus = iota
	// DeliveryStatusInProgress - В процессе доставки файлов
	DeliveryStatusInProgress DeliveryStatus = iota
	// DeliveryStatusDelivered - Файлы доставлены
	DeliveryStatusDelivered DeliveryStatus = iota
	// DeliveryStatusUpdating - Обновление раздачи
	DeliveryStatusUpdating DeliveryStatus = iota
	// DeliveryStatusDeleting - В процессе удаления файлов в диска
	DeliveryStatusDeleting DeliveryStatus = iota
	// DeliveryStatusDeleted - Файлы удалены
	DeliveryStatusDeleted DeliveryStatus = iota
)

type TorrentInfo struct {
	// href ссылки на торрент сайт раздачи
	Href *string
	// Magnet ссылка текущую раздачу
	Magnet *string
	// Хеш торрента
	Hash *string
}

// VideoContent информация о файлах
type VideoContent struct {
	// ID информации о файлах
	ID uuid.UUID
	// ID сериала/фильма
	ContentID common.ContentID
	// CreatedAt Датоа создания
	CreatedAt time.Time
	// Статус
	DeliveryStatus DeliveryStatus
	// Стейты привязанные к текущему контенту
	States []State
}

// State Таблица выпусков связанных с TVShowContent
type State struct {
	StateID uuid.UUID
	Type    runners.Type
}

type UpdateVideoContent struct {
	// Статус
	DeliveryStatus DeliveryStatus
	// Стейты привязанные к текущему контенту
	States []State
}
