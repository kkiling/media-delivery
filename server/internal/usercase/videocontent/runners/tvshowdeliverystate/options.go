package tvshowdeliverystate

import "github.com/kkiling/media-delivery/internal/usercase/videocontent/delivery"

type ChoseTorrentOptions struct {
	// Пользователь выбрал конкретный торрента файл
	Href *string
	// Пользователь поменял поисковый запрос
	NewSearchQuery *string
}

type ChoseFileMatchesOptions struct {
	// Пользователь подтверждает сметченные файлы
	Approve bool
	// Метч контента если указан
	ContentMatches *delivery.ContentMatches
}
