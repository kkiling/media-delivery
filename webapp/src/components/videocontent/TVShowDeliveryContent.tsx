import { ContentID, ErrorType, MediadeliveryStatus, TVShowDeliveryStatus } from '@/api/api';
import { observer } from 'mobx-react-lite';
import { useEffect } from 'react';
import { tvShowDeliveryStore } from '@/stores/tvShowDeliveryStore';
import { TorrentSearchList } from './TorrentSearchList';
import Loading from '../Loading';
import { ContentMatche } from './ContentMatches';
import { TorrentDownloadProgress, TorrentWaitingFiles } from './TorrentDownloadProgress';
import { MergeVideoProgress } from './MergeVideoProgress';
import { mapContentMatches } from './mapContentMatches';

const AUTO_RELOAD_TIMEOUT = 3000 as const;

interface TVShowDeliveryContentProps {
  contentId: ContentID;
  onNeedReload: () => void; // добавляем новый callback
}

export const TVShowDeliveryContent = observer(
  ({ contentId, onNeedReload }: TVShowDeliveryContentProps) => {
    useEffect(() => {
      tvShowDeliveryStore.fetchDeliveryData(contentId);
    }, [contentId]);

    useEffect(() => {
      let interval: NodeJS.Timeout | null = null;

      const nonPollingStates = [
        TVShowDeliveryStatus.WaitingChoseFileMatches,
        TVShowDeliveryStatus.WaitingUserChoseTorrent,
      ];

      if (
        tvShowDeliveryStore.deliveryState?.step &&
        !nonPollingStates.includes(tvShowDeliveryStore.deliveryState.step)
      ) {
        interval = setInterval(async () => {
          await tvShowDeliveryStore.fetchDeliveryData(contentId, true);

          // Check if step is Unknown and reload parent component data
          switch (tvShowDeliveryStore.deliveryState?.status) {
            case MediadeliveryStatus.FailedStatus:
            case MediadeliveryStatus.CompletedStatus:
              onNeedReload();
              break;
          }
        }, AUTO_RELOAD_TIMEOUT);
      }

      return () => {
        if (interval) {
          clearInterval(interval);
        }
      };
      // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [contentId, tvShowDeliveryStore.deliveryState?.step]);

    const onSearchSubmit = (query: string) => {
      tvShowDeliveryStore.selectTorrent(contentId, undefined, query);
    };

    const onTorrentSelect = async (href: string) => {
      await tvShowDeliveryStore.selectTorrent(contentId, href);
      document.getElementById('video-content-header')?.scrollIntoView();
    };

    const onConfirmFileMatches = async () => {
      await tvShowDeliveryStore.confirmFileMatches(contentId);
      document.getElementById('video-content-header')?.scrollIntoView();
    };

    const { loading, error, deliveryState } = tvShowDeliveryStore;

    // Показываем Loading только при первичной загрузке, не при фоновом обновлении
    if (loading && !tvShowDeliveryStore.backgroundLoading) return <Loading />;
    if (error) return <p className="text-danger">{error}</p>;
    if (!deliveryState) return <p className="text-muted">No delivery data available</p>;

    if (deliveryState.error) {
      switch (deliveryState.error.error_type) {
        case ErrorType.TorrentSiteForbidden:
          return <p className="text-danger">Torrent site is unavailable</p>;
        case ErrorType.FilesAlreadyExist:
          return (
            <p className="text-danger">
              There are already files in the season directory on the media server, please delete
              them
            </p>
          );
        default:
          return <p className="text-danger">{deliveryState.error.raw_error}</p>;
      }
    }

    switch (deliveryState.step) {
      case TVShowDeliveryStatus.GenerateSearchQuery:
      case TVShowDeliveryStatus.SearchTorrents:
        // Поиск раздач
        return <Loading text="Search torrents" />;
      case TVShowDeliveryStatus.WaitingUserChoseTorrent:
        // Отображение списка найденных торрентов
        return (
          <TorrentSearchList
            torrents={deliveryState.data?.torrent_search}
            searchQuery={deliveryState.data?.search_query?.Query}
            onSearch={onSearchSubmit}
            onSelect={onTorrentSelect}
          />
        );
      case TVShowDeliveryStatus.GetMagnetLink:
      case TVShowDeliveryStatus.AddTorrentToTorrentClient:
      case TVShowDeliveryStatus.PrepareFileMatches:
        // Ожидание получения списка файлов
        return <Loading text="Waiting torrent files" />;
      case TVShowDeliveryStatus.WaitingTorrentFiles:
        // Ожидание получения списка файлов
        return <TorrentWaitingFiles status={deliveryState.data?.torrent_download_status} />;
      case TVShowDeliveryStatus.WaitingChoseFileMatches:
        // Отображение выбора совпадений файлов
        return (
          <ContentMatche
            loading={loading}
            contentMatches={mapContentMatches(deliveryState.data?.content_matches || [])}
            onConfirm={onConfirmFileMatches}
          />
        );
      case TVShowDeliveryStatus.CreateVideoContentCatalogs:
      case TVShowDeliveryStatus.DeterminingNeedConvertFiles:
      case TVShowDeliveryStatus.StartMergeVideoFiles:
        // Обрабатываем файлы
        return <Loading text="Prepare files" />;
      case TVShowDeliveryStatus.CreateHardLinkCopy:
      case TVShowDeliveryStatus.GetCatalogsSize:
      case TVShowDeliveryStatus.SetMediaMetaData:
      case TVShowDeliveryStatus.SendDeliveryNotification:
      case TVShowDeliveryStatus.GetEpisodesData:
      case TVShowDeliveryStatus.TVShowDeliveryStatusUnknown:
        // Последние приготовления
        return <Loading text="There's just a little bit left" />;
      case TVShowDeliveryStatus.WaitingTorrentDownloadComplete:
        return <TorrentDownloadProgress status={deliveryState.data?.torrent_download_status} />;
      case TVShowDeliveryStatus.WaitingMergeVideoFiles:
        return <MergeVideoProgress status={deliveryState.data?.merge_video_status} />;
      default:
        return <Loading text="Unknown status" />;
    }
  }
);
