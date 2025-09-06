package matchtvshow

type Track struct {
	Name     string
	Language *string
	File     string
}

type Episode struct {
	EpisodeNumber int
	SeasonNumber  uint8
	VideoFile     string
	AudioFiles    []Track
	Subtitles     []Track
}
