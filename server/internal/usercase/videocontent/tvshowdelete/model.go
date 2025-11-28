package tvshowdelete

// TVShowCatalogPath пути каталогов сериала и сезона на медиа сервере
type TVShowCatalogPath struct {
	// Путь до каталога сериала
	TVShowPath string
	// Путь до каталога сезона (относительно каталога сериала)
	SeasonPath string
}
