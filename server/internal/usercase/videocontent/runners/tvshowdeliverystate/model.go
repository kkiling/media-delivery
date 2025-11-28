package tvshowdeliverystate

import (
	"fmt"

	"github.com/google/uuid"

	"github.com/kkiling/media-delivery/internal/common"
	"github.com/kkiling/media-delivery/internal/usercase/videocontent/tvshowdelivery"
)

// StepDelivery статус доставки видео файлов до медиа сервера
type StepDelivery string

const (
	// GenerateSearchQuery - генерация запросса к трекеру
	GenerateSearchQuery StepDelivery = "generate_search_query"
	// SearchTorrents - ищем раздачи сезона сериала / фильма
	SearchTorrents StepDelivery = "search_torrents"
	// WaitingUserChoseTorrent - ожидание когда пользователь выберет раздачу
	WaitingUserChoseTorrent StepDelivery = "waiting_user_chose_torrent"
	// GetMagnetLink получение магнет ссылки
	GetMagnetLink StepDelivery = "get_magnet_link"
	// AddTorrentToTorrentClient Добавление раздачи для скачивания торрент клиентом
	AddTorrentToTorrentClient StepDelivery = "add_torrent_to_torrent_client"
	// WaitingTorrentFiles Ожидание когда появится информация о файлах в раздаче
	WaitingTorrentFiles StepDelivery = "waiting_torrent_files"
	// GetEpisodesData получение информации о эпизодах и каталоге сезона
	GetEpisodesData StepDelivery = "get_episodes_data"
	// PrepareFileMatches получение информации о файлах раздачи
	PrepareFileMatches StepDelivery = "prepare_file_matches"
	// WaitingChoseFileMatches ожидание подтверждения пользователем соответствий выбора файлов
	WaitingChoseFileMatches StepDelivery = "waiting_chose_file_matches"
	// WaitingTorrentDownloadComplete ожидание завершения окончания скачивания раздачи
	WaitingTorrentDownloadComplete StepDelivery = "waiting_torrent_download_complete"
	// CreateVideoContentCatalogs Формирование каталогов и иерархии файлов
	CreateVideoContentCatalogs StepDelivery = "create_video_content_catalogs"
	// DeterminingNeedConvertFiles Определение необходимости конвертации файлов
	DeterminingNeedConvertFiles StepDelivery = "determining_need_convert_files"
	// --- Ветвь если необходимо добавление аудио дорожек/субтитров

	// StartMergeVideoFiles Запуск конвертирование файлов - полученные файлы сразу сохраняются в каталог медиасервера
	StartMergeVideoFiles StepDelivery = "merge_video_files"
	// WaitingMergeVideoFiles ожидание завершения конвертации файлов
	WaitingMergeVideoFiles StepDelivery = "waiting_merge_video_files"

	// -- Ветвь если не нужно изменять исходные файлы

	// CreateHardLinkCopy Копирование файлов из раздачи в каталог медиасервера (точнее создание симлинков)
	CreateHardLinkCopy StepDelivery = "create_hardlink_copy"
	// GetCatalogsSize получение размеров каталогов сериала
	GetCatalogsSize = "get_catalogs_size"

	// SetMediaMetaData установка методаных серий сезона сериала / фильма в медиасервере
	SetMediaMetaData StepDelivery = "set_media_meta_data"
	// AddLabel Установить лейбл для видеоконтента
	AddLabel StepDelivery = "add_label"
	// SendDeliveryNotification Отправка уведомления в telegramm о успешной доставки видеофайлов до медиа сервера
	SendDeliveryNotification StepDelivery = "send_delivery_notification"
)

// TVShowDeliveryData модель содержащая информацию о видео контенте для сезона сериала / фильма
/*
	К одному сезону сериала / фильму может быть привязано несколько VideoContent,
	но для упрощения пока будем пока разрешать только 1
*/
type TVShowDeliveryData struct {
	// SearchQuery сформированный запрос на основе названия сериала
	SearchQuery *tvshowdelivery.SearchQuery
	// TorrentSearch Результат поиска торрентов
	TorrentSearch []tvshowdelivery.TorrentSearch
	// Torrent данные найденной раздачи
	Torrent *tvshowdelivery.Torrent
	// TorrentFilesData файлы раздачи
	TorrentFilesData *tvshowdelivery.TorrentFilesData
	// EpisodesData информация о эпизодах и путях сохранения
	EpisodesData *tvshowdelivery.EpisodesData
	// ContentMatches Информация о метче файлов (метч видофайлов с аудиодоржками и субтитрами)
	ContentMatches *tvshowdelivery.ContentMatches
	// TorrentDownloadStatus статус скачивания раздачи
	TorrentDownloadStatus *tvshowdelivery.TorrentDownloadStatus
	// MergeIDs информация
	MergeIDs []uuid.UUID
	// MergeVideoStatus статус сшивания файлов (если нужен)
	MergeVideoStatus *tvshowdelivery.MergeVideoStatus
	// TVShowCatalogInfo информация о каталогах сериала
	TVShowCatalogInfo *tvshowdelivery.TVShowCatalog
}

type CreateOptions struct {
	Index    int
	TVShowID common.TVShowID
}

func (c CreateOptions) GetIdempotencyKey() string {
	return fmt.Sprintf("delivery_tv_%d_season_%d_n_%d", c.TVShowID.ID, c.TVShowID.SeasonNumber, c.Index)
}
