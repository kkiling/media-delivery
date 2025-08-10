package videocontent

import (
	"context"

	"github.com/samber/lo"

	"github.com/kkiling/torrent-to-media-server/internal/server/handler"
	"github.com/kkiling/torrent-to-media-server/internal/usercase/videocontent"
	desc "github.com/kkiling/torrent-to-media-server/pkg/gen/torrent-to-media-server"
)

func (h *Handler) CreateVideoContent(ctx context.Context, request *desc.CreateVideoContentRequest) (*desc.CreateVideoContentResponse, error) {
	contentID := mapContentIDReq(request.ContentId)

	result, err := h.videoContent.CreateVideoContent(ctx, videocontent.CreateVideoContentParams{
		ContentID: contentID,
	})

	if err != nil {
		return nil, handler.HandleError(err, "videoContent.CreateVideoContent")
	}

	return &desc.CreateVideoContentResponse{
		Result: mapVideoContent(*result),
	}, nil
}

func (h *Handler) GetVideoContent(ctx context.Context, request *desc.GetVideoContentRequest) (*desc.GetVideoContentResponse, error) {
	contentID := mapContentIDReq(request.ContentId)

	items, err := h.videoContent.GetVideoContent(ctx, contentID)
	if err != nil {
		return nil, handler.HandleError(err, "videoContent.GetVideoContent")
	}

	return &desc.GetVideoContentResponse{
		Items: lo.Map(items, func(it videocontent.VideoContent, _ int) *desc.VideoContent {
			return mapVideoContent(it)
		}),
	}, nil
}

func (h *Handler) GetTVShowDeliveryData(ctx context.Context, request *desc.GetTVShowDeliveryDataRequest) (*desc.GetTVShowDeliveryDataResponse, error) {
	contentID := mapContentIDReq(request.ContentId)

	state, err := h.videoContent.GetTVShowDeliveryData(ctx, contentID)
	if err != nil {
		return nil, handler.HandleError(err, "videoContent.GetTVShowDeliveryData")
	}

	return &desc.GetTVShowDeliveryDataResponse{
		Result: mapTVShowDeliveryState(state),
	}, nil
}

func (h *Handler) ChoseTorrentOptions(ctx context.Context, request *desc.ChoseTorrentOptionsRequest) (*desc.ChoseTorrentOptionsResponse, error) {
	contentID := mapContentIDReq(request.ContentId)

	state, err := h.videoContent.ChoseTorrentOptions(ctx, contentID, videocontent.ChoseTorrentOptions{
		Href:           request.Href,
		NewSearchQuery: request.NewSearchQuery,
	})
	if err != nil {
		return nil, handler.HandleError(err, "videoContent.ChoseTorrentOptions")
	}

	return &desc.ChoseTorrentOptionsResponse{
		Result: mapTVShowDeliveryState(state),
	}, nil
}

func (h *Handler) ChoseFileMatchesOptions(ctx context.Context, request *desc.ChoseFileMatchesOptionsRequest) (*desc.ChoseFileMatchesOptionsResponse, error) {
	contentID := mapContentIDReq(request.ContentId)

	state, err := h.videoContent.ChoseFileMatchesOptions(ctx, contentID, videocontent.ChoseFileMatchesOptions{
		Approve: request.Approve,
	})
	if err != nil {
		return nil, handler.HandleError(err, "videoContent.ChoseTorrentOptions")
	}

	return &desc.ChoseFileMatchesOptionsResponse{
		Result: mapTVShowDeliveryState(state),
	}, nil
}
