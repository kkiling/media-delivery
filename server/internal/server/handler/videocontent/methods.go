package videocontent

import (
	"context"

	"github.com/kkiling/media-delivery/internal/server/handler/videocontent/mapfrom"
	"github.com/kkiling/media-delivery/internal/server/handler/videocontent/mapto"

	"github.com/kkiling/media-delivery/internal/server/handler"
	"github.com/kkiling/media-delivery/internal/usercase/videocontent"
	desc "github.com/kkiling/media-delivery/pkg/gen/media-delivery"
)

func (h *Handler) CreateVideoContent(ctx context.Context, request *desc.CreateVideoContentRequest) (*desc.CreateVideoContentResponse, error) {
	contentID := mapfrom.ContentID(request.ContentId)

	result, err := h.videoContent.CreateVideoContent(ctx, videocontent.CreateVideoContentParams{
		ContentID: contentID,
	})

	if err != nil {
		return nil, handler.HandleError(err, "videoContent.CreateVideoContent")
	}

	return &desc.CreateVideoContentResponse{
		Result: mapto.VideoContent(*result),
	}, nil
}

func (h *Handler) GetVideoContent(ctx context.Context, request *desc.GetVideoContentRequest) (*desc.GetVideoContentResponse, error) {
	contentID := mapfrom.ContentID(request.ContentId)

	items, err := h.videoContent.GetVideoContent(ctx, contentID)
	if err != nil {
		return nil, handler.HandleError(err, "videoContent.GetVideoContent")
	}

	return &desc.GetVideoContentResponse{
		Items: mapto.VideoContents(items),
	}, nil
}

func (h *Handler) GetTVShowDeliveryData(ctx context.Context, request *desc.GetTVShowDeliveryDataRequest) (*desc.GetTVShowDeliveryDataResponse, error) {
	contentID := mapfrom.ContentID(request.ContentId)

	state, err := h.videoContent.GetTVShowDeliveryData(ctx, contentID)
	if err != nil {
		return nil, handler.HandleError(err, "videoContent.GetTVShowDeliveryData")
	}

	return &desc.GetTVShowDeliveryDataResponse{
		Result: mapto.TVShowDeliveryState(state),
	}, nil
}

func (h *Handler) ChoseTorrentOptions(ctx context.Context, request *desc.ChoseTorrentOptionsRequest) (*desc.ChoseTorrentOptionsResponse, error) {
	contentID := mapfrom.ContentID(request.ContentId)

	state, err := h.videoContent.ChoseTorrentOptions(ctx, contentID, videocontent.ChoseTorrentOptions{
		Href:           request.Href,
		NewSearchQuery: request.NewSearchQuery,
	})
	if err != nil {
		return nil, handler.HandleError(err, "videoContent.ChoseTorrentOptions")
	}

	return &desc.ChoseTorrentOptionsResponse{
		Result: mapto.TVShowDeliveryState(state),
	}, nil
}

func (h *Handler) ChoseFileMatchesOptions(ctx context.Context, request *desc.ChoseFileMatchesOptionsRequest) (*desc.ChoseFileMatchesOptionsResponse, error) {
	contentID := mapfrom.ContentID(request.ContentId)

	state, err := h.videoContent.ChoseFileMatchesOptions(ctx, contentID, videocontent.ChoseFileMatchesOptions{
		Approve:        request.Approve,
		ContentMatches: mapfrom.ContentMatches(request.ContentMatches),
	})
	if err != nil {
		return nil, handler.HandleError(err, "videoContent.ChoseTorrentOptions")
	}

	return &desc.ChoseFileMatchesOptionsResponse{
		Result: mapto.TVShowDeliveryState(state),
	}, nil
}
