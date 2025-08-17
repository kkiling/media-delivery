package delivery

import (
	"path/filepath"

	"github.com/kkiling/media-delivery/internal/adapter/qbittorrent"
)

// SearchQuery поисковый запрос для поиска раздачи
type SearchQuery struct {
	// Поисковый запрос с которым ищем раздачи на торрент сайте
	Query string
	// Предложенные вариации поискового запроса
	OptionalQuery []string
}

// TorrentSearch результат поиска торрент раздачи
type TorrentSearch struct {
	Title     string
	Href      string
	Size      string // TODO: переделать на int
	Seeds     string // TODO: переделать на int
	Leeches   string // TODO: переделать на int
	Downloads string // TODO: переделать на int
	AddedDate string // TODO: переделать на time.Time
}

type MagnetLink struct {
	// Магнет ссылка
	Magnet string
	// Хеш раздачи
	Hash string
}

// Torrent Данные раздачи
type Torrent struct {
	// Ссылка на раздачу
	Href string
	//
	MagnetLink *MagnetLink
}

type TorrentFilesData struct {
	// ContentFullPath Полный путь до каталога скачивания
	ContentFullPath string
	// Files Файлы раздачи
	Files []FileInfo
}

// TVShowCatalogPath пути каталогов сериала и сезона на медиа сервере
type TVShowCatalogPath struct {
	// Путь до каталога сериала
	TVShowPath string
	// Путь до каталога сезона (относительно каталога сериала)
	SeasonPath string
}

func (m TVShowCatalogPath) FullSeasonPath() string {
	return filepath.Join(m.TVShowPath, m.SeasonPath)
}

type EpisodesData struct {
	TVShowCatalogPath TVShowCatalogPath
	Episodes          []EpisodeInfo
}

type EpisodeInfo struct {
	// Номер сезона
	SeasonNumber uint8
	// Номер эпизода
	EpisodeNumber int
	// Наименования эпизода
	EpisodeName string
	// Наименование файла эпизода который будет лежать на медиасервере (без расширения файла)
	FileName string
}

type FileInfo struct {
	// Относительный путь до файла торрента (относительно каталога скачивания)
	RelativePath string
	// Полный путь до файла в системе
	FullPath string
	// Размер файла в байтах
	Size int64
	// Расширение файла
	Extension string
}

type VideoFile struct {
	File FileInfo
}

type Track struct {
	Name     string
	Language string
	File     FileInfo
}

// ContentMatches сопоставление видео файла с торрент файлом
type ContentMatches struct {
	Episode    EpisodeInfo
	Video      VideoFile
	AudioFiles []Track
	Subtitles  []Track
}

type TorrentState string

const (
	TorrentStateError       TorrentState = "error"       // Ошибка
	TorrentStateUploading   TorrentState = "uploading"   // Раздача (сидирование)
	TorrentStateDownloading TorrentState = "downloading" // Загрузка
	TorrentStateStopped     TorrentState = "stopped"     // Приостановлен (все виды паузы)
	TorrentStateQueued      TorrentState = "queued"      // В очереди
	TorrentStateUnknown     TorrentState = "unknown"     // Неизвестный статус
)

func mapTorrentState(qbState qbittorrent.TorrentState) TorrentState {
	switch qbState {
	// Ошибки
	case qbittorrent.TorrentStateError,
		qbittorrent.TorrentStateMissingFiles:
		return TorrentStateError

	// Раздача (сидирование)
	case qbittorrent.TorrentStateUploading,
		qbittorrent.TorrentStatePausedUP,
		qbittorrent.TorrentStateStalledUP,
		qbittorrent.TorrentStateCheckingUP,
		qbittorrent.TorrentStateQueuedUP,
		qbittorrent.TorrentStateForcedUP:
		return TorrentStateUploading

	// Загрузка
	case qbittorrent.TorrentStateDownloading,
		qbittorrent.TorrentStateMetaDL,
		qbittorrent.TorrentStateAllocating,
		qbittorrent.TorrentStateForcedDL,
		qbittorrent.TorrentStateMoving:
		return TorrentStateDownloading

	// Приостановлен
	case qbittorrent.TorrentStatePausedDL,
		qbittorrent.TorrentStateStoppedDL:
		return TorrentStateStopped

	// В очереди
	case qbittorrent.TorrentStateQueuedDL,
		qbittorrent.TorrentStateCheckingDL,
		qbittorrent.TorrentStateStalledDL,
		qbittorrent.TorrentStateCheckingResumeData:
		return TorrentStateQueued

	default:
		return TorrentStateUnknown
	}
}

type TorrentDownloadStatus struct {
	State      TorrentState
	Progress   float64
	IsComplete bool
}

// TVShowCatalog сводная информация о каталогах и размерах
type TVShowCatalog struct {
	// Путь до раздачи сезона сериала
	TorrentPath string
	// Размер файлов раздачи сезона сериала (байты)
	TorrentSize uint64
	// Путь до сезона сериала на медиасервере
	MediaServerPath TVShowCatalogPath
	// Размер файлов раздачи сезона сериала (байты)
	MediaServerSize uint64
	// Файлы скопированы с раздачи или созданы ссылочная связь
	// True - файлы скопированы
	// False - файлы созданы через линки
	IsCopyFilesInMediaServer bool
}
