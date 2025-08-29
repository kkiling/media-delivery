package rutracker

type Torrent struct {
	// Title наименование
	Title string
	// Href ссылка на раздачу
	Href string
	// Author автор раздачи
	Author string
	// Size Размер раздачи (байты)
	SizeBytes uint64
	// Размер виде строки (32Gb)
	SizePretty string
	// Seeds Информация о сидах
	Seeds uint32
	// Leeches Информация о личах
	Leeches uint32
	// Downloads количество скачиваний
	Downloads uint32
	// AddedDate Дата добавления
	AddedDate string
	// Категория
	Category string
}

type TorrentResponse struct {
	Page         int
	TotalResults int
	Results      []Torrent
}

type MagnetInfo struct {
	Magnet string
	Hash   string
}
