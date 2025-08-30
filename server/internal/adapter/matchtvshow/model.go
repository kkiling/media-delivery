package matchtvshow

type TVShowSeasonInfo struct {
	TVShowName    string
	FirstAirYear  string
	SeasonName    string
	SeasonNumber  uint8
	SeasonAirYear string
}

type EpisodeInfo struct {
	// Номер эпизода в сезоне
	EpisodeNumber int
	// Наименования эпизода
	EpisodeName string
}

type PrepareTrack struct {
	Name     string
	Language string
	File     string
}

type PrepareVideo struct {
	File string
}

type PrepareEpisode struct {
	Episode    EpisodeInfo
	VideoFile  *PrepareVideo
	AudioFiles []PrepareTrack
	Subtitles  []PrepareTrack
}

type PrepareTVShowSeason struct {
	Episodes []PrepareEpisode
}
