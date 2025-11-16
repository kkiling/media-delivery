package tvshowdeletestate

import (
	"fmt"

	"github.com/kkiling/media-delivery/internal/common"
	"github.com/kkiling/media-delivery/internal/usercase/videocontent/tvshowdelete"
)

// StepDelete статус удаления файлов видео контента
type StepDelete string

const (
	// StartDeleteTVShowSeason начальный шаг для удаленис сезона сериала
	StartDeleteTVShowSeason StepDelete = "start_delete_tv_show_season"
	// DeleteTorrentFromTorrentClient - удаление торрент раздачи из торрент клиента
	DeleteTorrentFromTorrentClient StepDelete = "delete_torrent_from_torrent_client"
	// DeleteTorrentFiles удаление файлов раздачи с диска
	DeleteTorrentFiles StepDelete = "delete_torrent_files"
	// DeleteSeasonFiles удаление файлов сезона с медиасервера
	DeleteSeasonFiles StepDelete = "delete_season_files_from_media_server"
	// DeleteSeasonFromMediaServer удаление сезона сериала с медиасервера
	DeleteSeasonFromMediaServer StepDelete = "delete_season_from_media_server"
	// DeleteLabel удаление лейбла
	DeleteLabel StepDelete = "delete_label"
)

// TVShowDeleteData модель содержащая информацию о процессе удаления файлов видео контента
type TVShowDeleteData struct {
	// Хеш торрент раздачи (что бы удалить раздачу в торрент клиенте)
	MagnetHash string
	// Путь до раздачи сезона сериала
	TorrentPath string
	// Путь каталога сериала и сезона на медиа сервере
	TVShowCatalogPath tvshowdelete.TVShowCatalogPath
}

type CreateOptions struct {
	TVShowID common.TVShowID
	// Хеш торрент раздачи (что бы удалить раздачу в торрент клиенте)
	MagnetHash string
	// Путь до раздачи сезона сериала
	TorrentPath string
	// Путь каталога сериала и сезона на медиа сервере
	TVShowCatalogPath tvshowdelete.TVShowCatalogPath
}

func (c CreateOptions) GetIdempotencyKey() string {
	return fmt.Sprintf("delivery_tv_%d_season_%d", c.TVShowID.ID, c.TVShowID.SeasonNumber)
}
