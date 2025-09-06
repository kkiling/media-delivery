package videocontent

import (
	"github.com/kkiling/media-delivery/internal/usercase/videocontent/common"
	"github.com/kkiling/media-delivery/internal/usercase/videocontent/content"
	"github.com/kkiling/media-delivery/internal/usercase/videocontent/delivery"
	"github.com/kkiling/media-delivery/internal/usercase/videocontent/runners/tvshowdeliverystate"
)

type TVShowID = common.TVShowID
type ContentID = common.ContentID
type VideoContent = content.VideoContent

type TorrentState = delivery.TorrentState

const (
	TorrentStateError       = delivery.TorrentStateError
	TorrentStateUploading   = delivery.TorrentStateUploading
	TorrentStateDownloading = delivery.TorrentStateDownloading
	TorrentStateStopped     = delivery.TorrentStateStopped
	TorrentStateQueued      = delivery.TorrentStateQueued
	TorrentStateUnknown     = delivery.TorrentStateUnknown
)

type DeliveryStatus = content.DeliveryStatus

const (
	DeliveryStatusFailed     = content.DeliveryStatusFailed
	DeliveryStatusInProgress = content.DeliveryStatusInProgress
	DeliveryStatusDelivered  = content.DeliveryStatusDelivered
	DeliveryStatusUpdating   = content.DeliveryStatusUpdating
	DeliveryStatusDeleting   = content.DeliveryStatusDeleting
	DeliveryStatusDeleted    = content.DeliveryStatusDeleted
)

type CreateVideoContentParams = content.CreateVideoContentParams
type TVShowDeliveryState = tvshowdeliverystate.State
type TVShowDeliveryData = tvshowdeliverystate.TVShowDeliveryData
type TorrentSearch = delivery.TorrentSearch
type ContentMatches = delivery.ContentMatches
type ContentMatch = delivery.ContentMatch
type Track = delivery.Track

type TrackType = delivery.TrackType

const (
	TrackTypeVideo    = delivery.TrackTypeVideo
	TrackTypeAudio    = delivery.TrackTypeAudio
	TrackTypeSubtitle = delivery.TrackTypeSubtitle
)

type FileInfo = delivery.FileInfo
type ChoseTorrentOptions = tvshowdeliverystate.ChoseTorrentOptions
type ChoseFileMatchesOptions = tvshowdeliverystate.ChoseFileMatchesOptions

type StepDelivery = tvshowdeliverystate.StepDelivery

const (
	GenerateSearchQuery            = tvshowdeliverystate.GenerateSearchQuery
	SearchTorrents                 = tvshowdeliverystate.SearchTorrents
	WaitingUserChoseTorrent        = tvshowdeliverystate.WaitingUserChoseTorrent
	GetMagnetLink                  = tvshowdeliverystate.GetMagnetLink
	AddTorrentToTorrentClient      = tvshowdeliverystate.AddTorrentToTorrentClient
	PrepareFileMatches             = tvshowdeliverystate.PrepareFileMatches
	WaitingChoseFileMatches        = tvshowdeliverystate.WaitingChoseFileMatches
	WaitingTorrentDownloadComplete = tvshowdeliverystate.WaitingTorrentDownloadComplete
	CreateVideoContentCatalogs     = tvshowdeliverystate.CreateVideoContentCatalogs
	DeterminingNeedConvertFiles    = tvshowdeliverystate.DeterminingNeedConvertFiles
	StartMergeVideoFiles           = tvshowdeliverystate.StartMergeVideoFiles
	WaitingMergeVideoFiles         = tvshowdeliverystate.WaitingMergeVideoFiles
	CreateHardLinkCopy             = tvshowdeliverystate.CreateHardLinkCopy
	GetCatalogsSize                = tvshowdeliverystate.GetCatalogsSize
	SetMediaMetaData               = tvshowdeliverystate.SetMediaMetaData
	SendDeliveryNotification       = tvshowdeliverystate.SendDeliveryNotification
	WaitingTorrentFiles            = tvshowdeliverystate.WaitingTorrentFiles
	GetEpisodesData                = tvshowdeliverystate.GetEpisodesData
)
