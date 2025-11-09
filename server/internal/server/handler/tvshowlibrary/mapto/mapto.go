package mapto

import (
	"github.com/samber/lo"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/kkiling/media-delivery/internal/usercase/tvshowlibrary"
	desc "github.com/kkiling/media-delivery/pkg/gen/media-delivery"
)

func Image(res *tvshowlibrary.Image) *desc.Image {
	if res == nil {
		return nil
	}
	return &desc.Image{
		Id:  res.ID,
		W92: *res.W92,
		//W154:     res.W154,
		W185: *res.W185,
		W342: *res.W342,
		//W500:     res.W500,
		//W780:     res.W780,
		Original: res.Original,
	}
}

func TVShowShort(res tvshowlibrary.TVShowShort) *desc.TVShowShort {
	return &desc.TVShowShort{
		Id:           res.ID,
		Name:         res.Name,
		OriginalName: res.OriginalName,
		Overview:     res.Overview,
		Poster:       Image(res.Poster),
		FirstAirDate: timestamppb.New(res.FirstAirDate),
		VoteAverage:  res.VoteAverage,
		VoteCount:    res.VoteCount,
		Popularity:   res.Popularity,
	}
}

func TVShowShorts(res []tvshowlibrary.TVShowShort) []*desc.TVShowShort {
	return lo.Map(res, func(item tvshowlibrary.TVShowShort, _ int) *desc.TVShowShort {
		return TVShowShort(item)
	})
}

func Season(res tvshowlibrary.Season) *desc.Season {
	return &desc.Season{
		// Id:           res.ID,
		AirDate:      timestamppb.New(res.AirDate),
		EpisodeCount: res.EpisodeCount,
		Name:         res.Name,
		Overview:     res.Overview,
		Poster:       Image(res.Poster),
		SeasonNumber: uint32(res.SeasonNumber),
		VoteAverage:  res.VoteAverage,
	}
}

func TVShow(res *tvshowlibrary.TVShow) *desc.TVShow {
	return &desc.TVShow{
		Id:               res.ID,
		Name:             res.Name,
		OriginalName:     res.OriginalName,
		Overview:         res.Overview,
		Poster:           Image(res.Poster),
		FirstAirDate:     timestamppb.New(res.FirstAirDate),
		VoteAverage:      res.VoteAverage,
		VoteCount:        res.VoteCount,
		Popularity:       res.Popularity,
		Backdrop:         Image(res.Backdrop),
		Genres:           res.Genres,
		LastAirDate:      timestamppb.New(res.LastAirDate),
		NumberOfEpisodes: res.NumberOfEpisodes,
		NumberOfSeasons:  res.NumberOfSeasons,
		OriginCountry:    res.OriginCountry,
		//Status:           res.Status,
		Tagline: res.Tagline,
		//Type:             res.Type,
		Seasons: lo.Map(res.Seasons, func(item tvshowlibrary.Season, _ int) *desc.Season {
			return Season(item)
		}),
	}
}

func Episode(res tvshowlibrary.Episode) *desc.Episode {
	return &desc.Episode{
		//Id:            res.ID,
		AirDate:       timestamppb.New(res.AirDate),
		EpisodeNumber: uint32(res.EpisodeNumber),
		//EpisodeType:   res.EpisodeType,
		Name:        res.Name,
		Overview:    res.Overview,
		Runtime:     uint32(res.Runtime),
		Still:       Image(res.Still),
		VoteAverage: res.VoteAverage,
		VoteCount:   uint32(res.VoteCount),
	}
}

func Episodes(res []tvshowlibrary.Episode) []*desc.Episode {
	return lo.Map(res, func(item tvshowlibrary.Episode, _ int) *desc.Episode {
		return Episode(item)
	})
}
