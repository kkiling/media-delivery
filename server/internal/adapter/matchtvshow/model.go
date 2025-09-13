package matchtvshow

type TrackType string

const (
	TrackTypeUnknown  TrackType = "unknown"
	TrackTypeVideo    TrackType = "video"
	TrackTypeAudio    TrackType = "audio"
	TrackTypeSubtitle TrackType = "subtitle"
)

type Track struct {
	Name     string
	Language *string
	File     string
	Type     TrackType
}

type ContentMatch struct {
	EpisodeNumber int
	SeasonNumber  uint8
	Video         *Track
	AudioTracks   []Track
	Subtitles     []Track
}

type ContentMatches struct {
	Matches     []ContentMatch
	Unallocated []Track
}
