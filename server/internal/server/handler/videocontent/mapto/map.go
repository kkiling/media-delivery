package mapto

import (
	"strings"

	"github.com/kkiling/statemachine"
	"github.com/samber/lo"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/kkiling/media-delivery/internal/usercase/videocontent"
	desc "github.com/kkiling/media-delivery/pkg/gen/media-delivery"
)

func ContentID(id *videocontent.ContentID) *desc.ContentID {
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

func VideoContent(in videocontent.VideoContent) *desc.VideoContent {
	return &desc.VideoContent{
		Id:             in.ID.String(),
		CreatedAt:      timestamppb.New(in.CreatedAt),
		DeliveryStatus: deliveryStatus(in.DeliveryStatus),
	}
}

func VideoContents(items []videocontent.VideoContent) []*desc.VideoContent {
	return lo.Map(items, func(it videocontent.VideoContent, _ int) *desc.VideoContent {
		return VideoContent(it)
	})
}

func tracks(tracks []videocontent.Track) []*desc.Track {
	return lo.Map(tracks, func(item videocontent.Track, _ int) *desc.Track {
		return &desc.Track{
			RelativePath: item.File.RelativePath,
			FullPath:     item.File.FullPath,
			Name:         item.Name,
			Language:     item.Language,
			Type:         trackType(item.Type),
		}
	})
}

func tvShowDeliveryData(state *videocontent.TVShowDeliveryState) *desc.TVShowDeliveryData {
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
					SeasonNumber:  uint32(item.Episode.SeasonNumber),
					EpisodeNumber: uint32(item.Episode.EpisodeNumber),
					FullPath:      item.Episode.FullPath,
					RelativePath:  item.Episode.RelativePath,
				},
				Video: &desc.Track{
					RelativePath: item.Video.File.RelativePath,
					FullPath:     item.Video.File.FullPath,
					Type:         trackType(item.Video.Type),
				},
				AudioTracks: tracks(item.AudioTracks),
				Subtitles:   tracks(item.Subtitles),
			}
		})

		result.ContentMatches = &desc.ContentMatches{
			Matches:     matches,
			Unallocated: tracks(cm.Unallocated),
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
				State:      torrentState(st.State),
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

func tvShowDeliveryError(state *videocontent.TVShowDeliveryState) *desc.TVShowDeliveryError {
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

func TVShowDeliveryState(state *videocontent.TVShowDeliveryState) *desc.TVShowDeliveryState {
	return &desc.TVShowDeliveryState{
		Data:   tvShowDeliveryData(state),
		Step:   tvShowDeliveryStep(state.Step),
		Status: status(state.Status),
		Error:  tvShowDeliveryError(state),
	}
}

func trackType(typeTrack videocontent.TrackType) desc.Track_TrackType {
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

func torrentState(state videocontent.TorrentState) desc.TorrentDownloadStatus_TorrentState {
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

func status(status statemachine.Status) desc.StateStatus {
	switch status {
	case statemachine.NewStatus:
		return desc.StateStatus_NewStatus
	case statemachine.InProgressStatus:
		return desc.StateStatus_InProgressStatus
	case statemachine.CompletedStatus:
		return desc.StateStatus_CompletedStatus
	case statemachine.FailedStatus:
		return desc.StateStatus_FailedStatus
	default:
		return desc.StateStatus_StatusUnknown
	}
}

func deliveryStatus(deliveryStatus videocontent.DeliveryStatus) desc.DeliveryStatus {
	switch deliveryStatus {
	case videocontent.DeliveryStatusFailed:
		return desc.DeliveryStatus_DeliveryStatusFailed
	case videocontent.DeliveryStatusNew:
		return desc.DeliveryStatus_DeliveryStatusNew
	case videocontent.DeliveryStatusInProgress:
		return desc.DeliveryStatus_DeliveryStatusInProgress
	case videocontent.DeliveryStatusDelivered:
		return desc.DeliveryStatus_DeliveryStatusDelivered
	case videocontent.DeliveryStatusUpdating:
		return desc.DeliveryStatus_DeliveryStatusUnknown
	case videocontent.DeliveryStatusDeleting:
		return desc.DeliveryStatus_DeliveryStatusDeleting
	case videocontent.DeliveryStatusDeleted:
		return desc.DeliveryStatus_DeliveryStatusDeleted
	default:
		return desc.DeliveryStatus_DeliveryStatusUnknown
	}
}

func tvShowDeliveryStep(step videocontent.StepDelivery) desc.TVShowDeliveryStep {
	switch step {
	case videocontent.GenerateSearchQuery:
		return desc.TVShowDeliveryStep_GenerateSearchQuery
	case videocontent.SearchTorrents:
		return desc.TVShowDeliveryStep_SearchTorrents
	case videocontent.WaitingUserChoseTorrent:
		return desc.TVShowDeliveryStep_WaitingUserChoseTorrent
	case videocontent.GetMagnetLink:
		return desc.TVShowDeliveryStep_GetMagnetLink
	case videocontent.AddTorrentToTorrentClient:
		return desc.TVShowDeliveryStep_AddTorrentToTorrentClient
	case videocontent.PrepareFileMatches:
		return desc.TVShowDeliveryStep_PrepareFileMatches
	case videocontent.WaitingChoseFileMatches:
		return desc.TVShowDeliveryStep_WaitingChoseFileMatches
	case videocontent.WaitingTorrentDownloadComplete:
		return desc.TVShowDeliveryStep_WaitingTorrentDownloadComplete
	case videocontent.CreateVideoContentCatalogs:
		return desc.TVShowDeliveryStep_CreateVideoContentCatalogs
	case videocontent.DeterminingNeedConvertFiles:
		return desc.TVShowDeliveryStep_DeterminingNeedConvertFiles
	case videocontent.StartMergeVideoFiles:
		return desc.TVShowDeliveryStep_StartMergeVideoFiles
	case videocontent.WaitingMergeVideoFiles:
		return desc.TVShowDeliveryStep_WaitingMergeVideoFiles
	case videocontent.CreateHardLinkCopy:
		return desc.TVShowDeliveryStep_CreateHardLinkCopy
	case videocontent.GetCatalogsSize:
		return desc.TVShowDeliveryStep_GetCatalogsSize
	case videocontent.SetMediaMetaData:
		return desc.TVShowDeliveryStep_SetMediaMetaData
	case videocontent.SendDeliveryNotification:
		return desc.TVShowDeliveryStep_SendDeliveryNotification
	case videocontent.WaitingTorrentFiles:
		return desc.TVShowDeliveryStep_WaitingTorrentFiles
	case videocontent.GetEpisodesData:
		return desc.TVShowDeliveryStep_GetEpisodesData

	default:
		return desc.TVShowDeliveryStep_TVShowDeliveryStepUnknown
	}
}

func TVShowDeleteState(state *videocontent.TVShowDeleteState) *desc.TVShowDeleteState {
	return &desc.TVShowDeleteState{
		Status: status(state.Status),
		Error:  tvShowDeleteError(state),
	}
}

func tvShowDeleteError(state *videocontent.TVShowDeleteState) *desc.TVShowDeleteError {
	if state.Error == nil {
		return nil
	}
	return &desc.TVShowDeleteError{
		RawError:  *state.Error,
		ErrorType: desc.TVShowDeleteError_TVShowDeleteError_Unknown,
	}
}
