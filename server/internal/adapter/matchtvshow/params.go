package matchtvshow

// PrepareTvShowPrams входыне параметры
type PrepareTvShowPrams struct {
	SeasonInfo   TVShowSeasonInfo
	Episodes     []EpisodeInfo
	TorrentFiles []string
}
