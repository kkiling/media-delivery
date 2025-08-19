package videocontent

import (
	"github.com/samber/lo"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/kkiling/media-delivery/internal/usercase/videocontent"
	desc "github.com/kkiling/media-delivery/pkg/gen/media-delivery"
)

func mapContentIDReq(id *desc.ContentID) videocontent.ContentID {
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
func mapContentID(id *videocontent.ContentID) *desc.ContentID {
	var result desc.ContentID
	if id.MovieID != nil {
		result.MovieId = id.MovieID
	}

	if id.TVShow != nil {
		result.TvShow = &desc.TVShowID{
			Id:           id.TVShow.ID,
			SeasonNumber: uint32(id.TVShow.SeasonNumber),
		}
	}

	return &result
}

func mapVideoContent(result videocontent.VideoContent) *desc.VideoContent {
	return &desc.VideoContent{
		Id:             result.ID.String(),
		CreatedAt:      timestamppb.New(result.CreatedAt),
		DeliveryStatus: mapDeliveryStatus(result.DeliveryStatus),
	}
}

func mapDeliveryStatus(deliveryStatus videocontent.DeliveryStatus) desc.DeliveryStatus {
	switch deliveryStatus {
	case videocontent.DeliveryStatusFailed:
		return desc.DeliveryStatus_DeliveryStatusFailed
	case videocontent.DeliveryStatusInProgress:
		return desc.DeliveryStatus_DeliveryStatusInProgress
	case videocontent.DeliveryStatusDelivered:
		return desc.DeliveryStatus_DeliveryStatusDelivered
	case videocontent.DeliveryStatusUpdating:
		return desc.DeliveryStatus_DeliveryStatusUnknown
	case videocontent.DeliveryStatusDeleting:
		return desc.DeliveryStatus_DeliveryStatusUnknown
	case videocontent.DeliveryStatusDeleted:
		return desc.DeliveryStatus_DeliveryStatusUnknown
	default:
		return desc.DeliveryStatus_DeliveryStatusUnknown
	}
}

func mapDeliveryStep(step videocontent.StepDelivery) desc.TVShowDeliveryStatus {
	switch step {
	case videocontent.GenerateSearchQuery:
		return desc.TVShowDeliveryStatus_GenerateSearchQuery
	case videocontent.SearchTorrents:
		return desc.TVShowDeliveryStatus_SearchTorrents
	case videocontent.WaitingUserChoseTorrent:
		return desc.TVShowDeliveryStatus_WaitingUserChoseTorrent
	case videocontent.GetMagnetLink:
		return desc.TVShowDeliveryStatus_GetMagnetLink
	case videocontent.AddTorrentToTorrentClient:
		return desc.TVShowDeliveryStatus_AddTorrentToTorrentClient
	case videocontent.PrepareFileMatches:
		return desc.TVShowDeliveryStatus_PrepareFileMatches
	case videocontent.WaitingChoseFileMatches:
		return desc.TVShowDeliveryStatus_WaitingChoseFileMatches
	case videocontent.WaitingTorrentDownloadComplete:
		return desc.TVShowDeliveryStatus_WaitingTorrentDownloadComplete
	case videocontent.CreateVideoContentCatalogs:
		return desc.TVShowDeliveryStatus_CreateVideoContentCatalogs
	case videocontent.DeterminingNeedConvertFiles:
		return desc.TVShowDeliveryStatus_DeterminingNeedConvertFiles
	case videocontent.StartMergeVideoFiles:
		return desc.TVShowDeliveryStatus_StartMergeVideoFiles
	case videocontent.WaitingMergeVideoFiles:
		return desc.TVShowDeliveryStatus_WaitingMergeVideoFiles
	case videocontent.CreateHardLinkCopy:
		return desc.TVShowDeliveryStatus_CreateHardLinkCopy
	case videocontent.GetCatalogsSize:
		return desc.TVShowDeliveryStatus_GetCatalogsSize
	case videocontent.SetMediaMetaData:
		return desc.TVShowDeliveryStatus_SetMediaMetaData
	case videocontent.SendDeliveryNotification:
		return desc.TVShowDeliveryStatus_SendDeliveryNotification
	case videocontent.WaitingTorrentFiles:
		return desc.TVShowDeliveryStatus_WaitingTorrentFiles
	case videocontent.GetEpisodesData:
		return desc.TVShowDeliveryStatus_GetEpisodesData

	default:
		return desc.TVShowDeliveryStatus_TVShowDeliveryStatusUnknown
	}
}

func maTorrentState(state videocontent.TorrentState) desc.TorrentDownloadStatus_TorrentState {
	switch state {
	case videocontent.TorrentStateError:
		return desc.TorrentDownloadStatus_TORRENT_STATE_ERROR
	case videocontent.TorrentStateUploading:
		return desc.TorrentDownloadStatus_TORRENT_STATE_UPLOADING
	case videocontent.TorrentStateDownloading:
		return desc.TorrentDownloadStatus_TORRENT_STATE_DOWNLOADING
	case videocontent.TorrentStateStopped:
		return desc.TorrentDownloadStatus_TORRENT_STATE_STOPPED
	case videocontent.TorrentStateQueued:
		return desc.TorrentDownloadStatus_TORRENT_STATE_QUEUED
	default:
		return desc.TorrentDownloadStatus_TORRENT_STATE_UNKNOWN
	}
}

func mapFile(file videocontent.FileInfo) *desc.FileInfo {
	return &desc.FileInfo{
		RelativePath: file.RelativePath,
		FullPath:     file.FullPath,
		Size:         file.Size,
		Extension:    file.Extension,
	}
}

func mapTracks(tracks []videocontent.Track) []*desc.Track {
	return lo.Map(tracks, func(item videocontent.Track, _ int) *desc.Track {
		return &desc.Track{
			File:     mapFile(item.File),
			Name:     item.Name,
			Language: item.Language,
		}
	})
}

func mapTVShowDeliveryData(step videocontent.StepDelivery, data *videocontent.TVShowDeliveryData) *desc.TVShowDeliveryData {

	result := &desc.TVShowDeliveryData{
		SearchQuery: &desc.SearchQuery{
			Query: data.SearchQuery.Query,
		},
	}

	if step == videocontent.WaitingUserChoseTorrent {
		result.TorrentSearch = lo.Map(data.TorrentSearch, func(item videocontent.TorrentSearch, _ int) *desc.TorrentSearch {
			return &desc.TorrentSearch{
				Title:     item.Title,
				Href:      item.Href,
				Size:      item.Size,
				Seeds:     item.Seeds,
				Leeches:   item.Leeches,
				Downloads: item.Downloads,
				AddedDate: item.AddedDate,
			}
		})
	}
	if step == videocontent.WaitingChoseFileMatches {
		result.ContentMatches = lo.Map(data.ContentMatches, func(item videocontent.ContentMatches, _ int) *desc.ContentMatches {
			return &desc.ContentMatches{
				Episode: &desc.EpisodeInfo{
					SeasonNumber:  uint32(item.Episode.SeasonNumber),
					EpisodeName:   item.Episode.EpisodeName,
					EpisodeNumber: uint32(item.Episode.EpisodeNumber),
				},
				Video: &desc.VideoFile{
					File: mapFile(item.Video.File),
				},
				AudioFiles: mapTracks(item.AudioFiles),
				Subtitles:  mapTracks(item.Subtitles),
			}
		})
	}

	if step == videocontent.WaitingTorrentDownloadComplete {
		st := data.TorrentDownloadStatus
		result.TorrentDownloadStatus = &desc.TorrentDownloadStatus{
			State:      maTorrentState(st.State),
			Progress:   float32(st.Progress),
			IsComplete: st.IsComplete,
		}
	}

	if step == videocontent.WaitingMergeVideoFiles {
		st := data.MergeVideoStatus

		result.MergeVideoStatus = &desc.MergeVideoStatus{
			Progress:   float32(st.Progress),
			IsComplete: st.IsComplete,
		}
	}

	return result
}

func mapTVShowDeliveryState(state *videocontent.TVShowDeliveryState) *desc.TVShowDeliveryState {
	return &desc.TVShowDeliveryState{
		Data: mapTVShowDeliveryData(state.Step, &state.Data),
		Step: mapDeliveryStep(state.Step),
	}
}
