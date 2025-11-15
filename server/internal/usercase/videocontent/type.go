package videocontent

import (
	"github.com/kkiling/media-delivery/internal/common"
	"github.com/kkiling/media-delivery/internal/usercase/videocontent/content"
	"github.com/kkiling/media-delivery/internal/usercase/videocontent/runners/tvshowdeliverystate"
	"github.com/kkiling/media-delivery/internal/usercase/videocontent/tvshowdelivery"
)

type TVShowID = common.TVShowID
type ContentID = common.ContentID
type VideoContent = content.VideoContent

type TorrentState = tvshowdelivery.TorrentState

const (
	TorrentStateError       = tvshowdelivery.TorrentStateError
	TorrentStateUploading   = tvshowdelivery.TorrentStateUploading
	TorrentStateDownloading = tvshowdelivery.TorrentStateDownloading
	TorrentStateStopped     = tvshowdelivery.TorrentStateStopped
	TorrentStateQueued      = tvshowdelivery.TorrentStateQueued
	TorrentStateUnknown     = tvshowdelivery.TorrentStateUnknown
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
type TorrentSearch = tvshowdelivery.TorrentSearch
type ContentMatches = tvshowdelivery.ContentMatches
type ContentMatch = tvshowdelivery.ContentMatch
type ContentMatchesOptions = tvshowdelivery.ContentMatchesOptions
type Track = tvshowdelivery.Track

type TrackType = tvshowdelivery.TrackType

const (
	TrackTypeVideo    = tvshowdelivery.TrackTypeVideo
	TrackTypeAudio    = tvshowdelivery.TrackTypeAudio
	TrackTypeSubtitle = tvshowdelivery.TrackTypeSubtitle
)

type FileInfo = tvshowdelivery.FileInfo
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
