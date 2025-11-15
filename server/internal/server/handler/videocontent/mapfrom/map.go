package mapfrom

import (
	"github.com/kkiling/media-delivery/internal/usercase/videocontent/tvshowdelivery"
	"github.com/samber/lo"

	"github.com/kkiling/media-delivery/internal/usercase/videocontent"
	desc "github.com/kkiling/media-delivery/pkg/gen/media-delivery"
)

func ContentID(id *desc.ContentID) videocontent.ContentID {
	if id == nil || (id.MovieId == nil && id.TvShow == nil) {
		return videocontent.ContentID{}
	}

	var result videocontent.ContentID

	if id.MovieId != nil {
		result.MovieID = id.MovieId
	}

	if id.TvShow != nil {
		result.TVShow = &videocontent.TVShowID{
			ID:           id.TvShow.Id,
			SeasonNumber: uint8(id.TvShow.SeasonNumber),
		}
	}

	return result
}

func ContentMatches(match *desc.ContentMatches) *videocontent.ContentMatches {
	if match == nil {
		return nil
	}
	toTrack := func(item *desc.Track) videocontent.Track {
		trackType := videocontent.TrackTypeVideo
		switch item.Type {
		case desc.Track_TRACK_TYPE_AUDIO:
			trackType = videocontent.TrackTypeAudio
		case desc.Track_TRACK_TYPE_VIDEO:
			trackType = videocontent.TrackTypeVideo
		case desc.Track_TRACK_TYPE_SUBTITLE:
			trackType = videocontent.TrackTypeSubtitle
		}

		return videocontent.Track{
			Type:     trackType,
			Name:     item.Name,
			Language: item.Language,
			File: tvshowdelivery.FileInfo{
				RelativePath: item.RelativePath,
				FullPath:     item.FullPath,
			},
		}
	}

	toTracks := func(items []*desc.Track) []videocontent.Track {
		return lo.Map(items, func(item *desc.Track, index int) videocontent.Track {
			return toTrack(item)
		})
	}

	return &videocontent.ContentMatches{
		Matches: lo.Map(match.Matches, func(item *desc.ContentMatch, index int) videocontent.ContentMatch {
			return videocontent.ContentMatch{
				Episode: tvshowdelivery.EpisodeInfo{
					SeasonNumber:  uint8(item.Episode.SeasonNumber),
					EpisodeNumber: int(item.Episode.EpisodeNumber),
					FullPath:      item.Episode.FullPath,
					RelativePath:  item.Episode.RelativePath,
				},
				Video: func() *videocontent.Track {
					if item.Video != nil && item.Video.RelativePath != "" {
						return lo.ToPtr(toTrack(item.Video))
					}
					return nil
				}(),
				AudioTracks: toTracks(item.AudioTracks),
				Subtitles:   toTracks(item.Subtitles),
			}
		}),
		Unallocated: toTracks(match.Unallocated),
		Options: tvshowdelivery.ContentMatchesOptions{
			KeepOriginalAudio:     match.Options.KeepOriginalAudio,
			KeepOriginalSubtitles: match.Options.KeepOriginalSubtitles,
			DefaultAudioTrackName: match.Options.DefaultAudioTrackName,
			DefaultSubtitleTrack:  match.Options.DefaultSubtitleTrack,
		},
	}
}
