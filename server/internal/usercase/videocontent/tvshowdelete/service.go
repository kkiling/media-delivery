package tvshowdelete

type Config struct {
	// BasePath Базовый путь от которого расположены все файлы торрента или медиа сервера
	// Например скачанные сериалы лежат по пути BasePath + TVShowTorrentSavePath
	BasePath string // "/nfs"
	// TVShowTorrentSavePath путь сохранения сериалов относительно торрент клиента
	TVShowTorrentSavePath string
	// TVShowMediaSaveTvShowsPath путь сохранения сериалов относительно медиа сервера
	TVShowMediaSaveTvShowsPath string
}

type Service struct {
	config        Config
	torrentClient TorrentClient
	embyApi       EmbyApi
}

func NewService(
	config Config,
	torrentClient TorrentClient,
	embyApi EmbyApi,
) *Service {
	return &Service{
		config:        config,
		torrentClient: torrentClient,
		embyApi:       embyApi,
	}
}
