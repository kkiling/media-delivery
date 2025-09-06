package videocontent

import (
	"strings"

	"github.com/kkiling/statemachine"
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

func mapStatus(status statemachine.Status) desc.Status {
	switch status {
	case statemachine.NewStatus:
		return desc.Status_NewStatus
	case statemachine.InProgressStatus:
		return desc.Status_InProgressStatus
	case statemachine.CompletedStatus:
		return desc.Status_CompletedStatus
	case statemachine.FailedStatus:
		return desc.Status_FailedStatus
	default:
		return desc.Status_StatusUnknown
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

func mapTracks(tracks []videocontent.Track) []*desc.Track {
	return lo.Map(tracks, func(item videocontent.Track, _ int) *desc.Track {
		return &desc.Track{
			RelativePath: item.File.RelativePath,
			Name:         item.Name,
			Language:     item.Language,
			Type:         mapTrackType(item.Type),
		}
	})
}

func mapTrackType(typeTrack videocontent.TrackType) desc.Track_TrackType {
	switch typeTrack {
	case videocontent.TrackTypeVideo:
		return desc.Track_TRACK_TYPE_VIDEO
	case videocontent.TrackTypeAudio:
		return desc.Track_TRACK_TYPE_AUDIO
	case videocontent.TrackTypeSubtitle:
		return desc.Track_TRACK_TYPE_SUBTITLE
	default:
		return desc.Track_TRACK_TYPE_UNKNOWN
	}
}

func mapTVShowDeliveryData(state *videocontent.TVShowDeliveryState) *desc.TVShowDeliveryData {
	result := &desc.TVShowDeliveryData{
		SearchQuery: &desc.SearchQuery{
			Query: func() string {
				if state.Data.SearchQuery == nil {
					return ""
				}
				return state.Data.SearchQuery.Query
			}(),
		},
	}

	if state.Step == videocontent.WaitingUserChoseTorrent {
		result.TorrentSearch = lo.Map(state.Data.TorrentSearch, func(item videocontent.TorrentSearch, _ int) *desc.TorrentSearch {
			return &desc.TorrentSearch{
				Title:     item.Title,
				Href:      item.Href,
				Size:      item.SizePretty,
				Seeds:     int64(item.Seeds),
				Leeches:   int64(item.Leeches),
				Downloads: int64(item.Downloads),
				AddedDate: item.AddedDate,
				Category:  item.Category,
			}
		})
	}
	if state.Step == videocontent.WaitingChoseFileMatches {
		cm := state.Data.ContentMatches
		matches := lo.Map(cm.Matches, func(item videocontent.ContentMatch, _ int) *desc.ContentMatch {
			return &desc.ContentMatch{
				Episode: &desc.EpisodeInfo{
					EpisodeNumber: uint32(item.Episode.EpisodeNumber),
					EpisodeName:   item.Episode.EpisodeName,
					RelativePath:  item.Episode.RelativePath,
				},
				Video: &desc.Track{
					RelativePath: item.Video.File.RelativePath,
					Type:         mapTrackType(item.Video.Type),
				},
				AudioFiles: mapTracks(item.AudioFiles),
				Subtitles:  mapTracks(item.Subtitles),
			}
		})

		result.ContentMatches = &desc.ContentMatches{
			Matches:     matches,
			Unallocated: mapTracks(cm.Unallocated),
			Options: &desc.ContentMatches_Options{
				KeepOriginalAudio:     cm.Options.KeepOriginalAudio,
				KeepOriginalSubtitles: cm.Options.KeepOriginalSubtitles,
				DefaultAudioTrackName: cm.Options.DefaultAudioTrackName,
				DefaultSubtitleTrack:  cm.Options.DefaultSubtitleTrack,
			},
		}
	}

	if state.Step == videocontent.WaitingTorrentDownloadComplete || state.Step == videocontent.WaitingTorrentFiles {
		st := state.Data.TorrentDownloadStatus
		if st != nil {
			result.TorrentDownloadStatus = &desc.TorrentDownloadStatus{
				State:      maTorrentState(st.State),
				Progress:   float32(st.Progress),
				IsComplete: st.IsComplete,
			}
		}
	}

	if state.Step == videocontent.WaitingMergeVideoFiles {
		st := state.Data.MergeVideoStatus

		result.MergeVideoStatus = &desc.MergeVideoStatus{
			Progress:   float32(st.Progress),
			IsComplete: st.IsComplete,
		}
	}

	if state.Status == statemachine.CompletedStatus {
		inf := state.Data.TVShowCatalogInfo
		result.TvShowCatalogInfo = &desc.TVShowCatalog{
			TorrentPath:       inf.TorrentPath,
			TorrentSizePretty: inf.TorrentSizePretty,
			MediaServerPath: &desc.TVShowCatalogPath{
				TvShowPath: inf.MediaServerPath.TVShowPath,
				SeasonPath: inf.MediaServerPath.SeasonPath,
			},
			MediaServerSizePretty:    inf.MediaServerSizePretty,
			IsCopyFilesInMediaServer: inf.IsCopyFilesInMediaServer,
		}
		result.Torrent = &desc.Torrent{
			Href: state.Data.Torrent.Href,
		}
	}

	return result
}

func mapTVShowDeliveryError(state *videocontent.TVShowDeliveryState) *desc.TVShowDeliveryError {
	if state.Error == nil {
		return nil
	}
	errorType := desc.TVShowDeliveryError_TVShowDeliveryError_Unknown
	if state.Step == videocontent.SearchTorrents {
		if strings.Contains(*state.Error, "Forbidden") {
			errorType = desc.TVShowDeliveryError_TorrentSiteForbidden
		}
	}
	if state.Step == videocontent.CreateVideoContentCatalogs {
		if strings.Contains(*state.Error, "already exists") {
			errorType = desc.TVShowDeliveryError_FilesAlreadyExist
		}
	}

	return &desc.TVShowDeliveryError{
		RawError:  *state.Error,
		ErrorType: errorType,
	}
}

func mapTVShowDeliveryState(state *videocontent.TVShowDeliveryState) *desc.TVShowDeliveryState {
	return &desc.TVShowDeliveryState{
		Data:   mapTVShowDeliveryData(state),
		Step:   mapDeliveryStep(state.Step),
		Status: mapStatus(state.Status),
		Error:  mapTVShowDeliveryError(state),
	}
}
