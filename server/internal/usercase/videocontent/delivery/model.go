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
	// Наименование раздачи
	Title string
	// Категория
	Category string
	// Ссылка на раздачу
	Href string
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
	// Путь каталога сериала и сезона на медиа сервере
	TVShowCatalogPath TVShowCatalogPath
	// Информация о эпизодах
	Episodes []EpisodeInfo
}

type EpisodeInfo struct {
	// Номер сезона
	SeasonNumber uint8
	// Номер эпизода
	EpisodeNumber int
	// Наименование файла эпизода который будет лежать на медиасервере
	FullPath string
	// Относительный путь относительно каталога сезона
	RelativePath string
}

type FileInfo struct {
	// Относительный путь до файла торрента (относительно каталога скачивания)
	RelativePath string
	// Полный путь до файла в системе
	FullPath string
}

type VideoFile struct {
	File FileInfo
}

type TrackType string

const (
	TrackTypeVideo    TrackType = "video"
	TrackTypeAudio    TrackType = "audio"
	TrackTypeSubtitle TrackType = "subtitle"
)

type Track struct {
	Type     TrackType
	Name     *string
	Language *string
	File     FileInfo
}

// ContentMatch сопоставление видео файла с торрент файлом
type ContentMatch struct {
	Episode     EpisodeInfo
	Video       *Track
	AudioTracks []Track
	Subtitles   []Track
}

type ContentMatchesOptions struct {
	// Оставлять оригинальные аудиодорожки (если они есть)
	KeepOriginalAudio bool
	// Оставлять оригинальные субтитры (если они есть)
	KeepOriginalSubtitles bool
	// Дефолтная аудиодорожка
	DefaultAudioTrackName *string
	// Дефолтные субтитры
	DefaultSubtitleTrack *string
}

type ContentMatches struct {
	Matches     []ContentMatch
	Unallocated []Track
	Options     ContentMatchesOptions
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
	//
	TorrentSizePretty string
	// Путь до сезона сериала на медиасервере
	MediaServerPath TVShowCatalogPath
	// Размер файлов раздачи сезона сериала (байты)
	MediaServerSize       uint64
	MediaServerSizePretty string
	// Файлы скопированы с раздачи или созданы ссылочная связь
	// True - файлы скопированы
	// False - файлы созданы через линки
	IsCopyFilesInMediaServer bool
}
