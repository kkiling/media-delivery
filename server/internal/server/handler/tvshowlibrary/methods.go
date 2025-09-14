package tvshowlibrary

import (
	"context"

	"github.com/kkiling/media-delivery/internal/server/handler"
	"github.com/kkiling/media-delivery/internal/server/handler/tvshowlibrary/mapto"
	"github.com/kkiling/media-delivery/internal/usercase/tvshowlibrary"
	desc "github.com/kkiling/media-delivery/pkg/gen/media-delivery"
)

func (h *Handler) SearchTVShow(ctx context.Context, request *desc.SearchTVShowRequest) (*desc.SearchTVShowResponse, error) {
	result, err := h.tvShowLibrary.SearchTVShow(ctx, tvshowlibrary.TVShowSearchParams{
		Query: request.Query,
	})

	if err != nil {
		return nil, handler.HandleError(err, "tvShowLibrary.SearchTVShow")
	}

	return &desc.SearchTVShowResponse{
		Items: mapto.TVShowShorts(result.Items),
	}, nil
}

func (h *Handler) GetTVShowInfo(ctx context.Context, request *desc.GetTVShowInfoRequest) (*desc.GetTVShowInfoResponse, error) {
	result, err := h.tvShowLibrary.GetTVShowInfo(ctx, tvshowlibrary.GetTVShowParams{
		TVShowID: request.TvShowId,
	})

	if err != nil {
		return nil, handler.HandleError(err, "tvShowLibrary.GetTVShowInfo")
	}

	return &desc.GetTVShowInfoResponse{Result: mapto.TVShow(result.Result)}, nil
}

func (h *Handler) GetSeasonInfo(ctx context.Context, request *desc.GetSeasonInfoRequest) (*desc.GetSeasonInfoResponse, error) {
	res, err := h.tvShowLibrary.GetSeasonInfo(ctx, tvshowlibrary.GetSeasonInfoParams{
		TVShowID:     request.TvShowId,
		SeasonNumber: uint8(request.SeasonNumber),
	})

	if err != nil {
		return nil, handler.HandleError(err, "tvShowLibrary.GetSeasonEpisodes")
	}

	return &desc.GetSeasonInfoResponse{
		Season:   mapto.Season(res.Result.Season),
		Episodes: mapto.Episodes(res.Result.Episodes),
	}, nil
}

func (h *Handler) GetTVShowsFromLibrary(ctx context.Context, request *desc.GetTVShowsFromLibraryRequest) (*desc.GetTVShowsFromLibraryResponse, error) {
	result, err := h.tvShowLibrary.GetTVShowsFromLibrary(ctx, tvshowlibrary.GetTVShowsFromLibraryParams{})

	if err != nil {
		return nil, handler.HandleError(err, "tvShowLibrary.GetSeasonEpisodes")
	}

	return &desc.GetTVShowsFromLibraryResponse{
		Items: mapto.TVShowShorts(result.Items),
	}, nil
}
