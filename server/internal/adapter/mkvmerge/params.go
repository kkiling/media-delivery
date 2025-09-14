package mkvmerge

type Track struct {
	Path     string
	Language *string
	Name     string
	Default  bool
}

type MergeParams struct {
	VideoInputFile  string
	VideoOutputFile string
	AudioTracks     []Track
	SubtitleTracks  []Track
	// Оставлять оригинальные аудиодорожки (если они есть)
	KeepOriginalAudio bool
	// Оставлять оригинальные субтитры (если они есть)
	KeepOriginalSubtitles bool
}
