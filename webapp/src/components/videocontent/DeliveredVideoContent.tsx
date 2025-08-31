import { ContentID, TVShowCatalog } from '@/api/api';
import { observer } from 'mobx-react-lite';
import { useEffect } from 'react';
import { tvShowDeliveryStore } from '@/stores/tvShowDeliveryStore';
import Loading from '../Loading';
import { HardDrive, FolderOpen, Database, HardDriveDownload, FileBox, Scale } from 'lucide-react';

interface DeliveredVideoContentProps {
  contentId: ContentID;
}

interface InfoItemProps {
  icon: React.ReactNode;
  label: string;
  value: string;
  isCode?: boolean;
  isLink?: boolean;
}

const InfoItem = ({ icon, label, value, isCode = false, isLink = false }: InfoItemProps) => (
  <div className="d-flex align-items-center mb-3 info-item">
    <div className="me-2 text-primary">{icon}</div>
    <div>
      <div className="text-muted small">{label}</div>
      {isLink ? (
        <a href={value} target="_blank" rel="noopener noreferrer" className="text-decoration-none">
          <code className="ms-0 cursor-pointer">{value}</code>
        </a>
      ) : isCode ? (
        <code className="ms-0">{value}</code>
      ) : (
        <div className="fw-medium">{value}</div>
      )}
    </div>
  </div>
);

export const DeliveredVideoContent = observer(({ contentId }: DeliveredVideoContentProps) => {
  useEffect(() => {
    tvShowDeliveryStore.fetchDeliveryData(contentId);
  }, [contentId]);

  const { loading, error, deliveryState } = tvShowDeliveryStore;

  if (loading && !tvShowDeliveryStore.backgroundLoading) {
    return <Loading />;
  }

  if (error) {
    return <p className="text-danger">{error}</p>;
  }

  if (!deliveryState?.data?.tv_show_catalog_info) {
    return <p className="text-muted">No catalog information available</p>;
  }

  const catalogInfo: TVShowCatalog = deliveryState.data.tv_show_catalog_info;
  const torrent = deliveryState.data.torrent;
  return (
    <div className="delivered-content">
      <style>{`
        .info-item:hover {
          background-color: rgba(0,0,0,.03);
          border-radius: 8px;
          padding: 8px;
          margin: -8px;
          transition: all 0.2s ease;
        }
        .info-item {
          padding: 8px;
          margin: -8px;
        }
        .info-item code {
          color: inherit;
          background-color: rgba(0,0,0,.03);
          padding: 2px 4px;
          border-radius: 4px;
        }
        .info-item code.cursor-pointer:hover {
          text-decoration: underline;
        }
      `}</style>

      <div className="row">
        <div className="col-md-6">
          <h6 className="mb-4 d-flex align-items-center">
            <HardDrive className="me-2" size={20} />
            Torrent Files
          </h6>
          <InfoItem
            icon={<FileBox size={18} />}
            label="Torrent Link"
            value={torrent?.href ?? 'No link available'}
            isCode={true}
            isLink={Boolean(torrent?.href)}
          />
          <InfoItem
            icon={<FolderOpen size={18} />}
            label="Torrent Location"
            value={catalogInfo.torrent_path ?? ''}
            isCode={true}
          />
          <InfoItem
            icon={<Database size={18} />}
            label="Size of Torrent Files"
            value={catalogInfo.torrent_size_pretty ?? ''}
          />
        </div>

        <div className="col-md-6">
          <h6 className="mb-4 d-flex align-items-center">
            <HardDriveDownload className="me-2" size={20} />
            Media Server Files
          </h6>
          <InfoItem
            icon={<FileBox size={18} />}
            label="Season Location"
            value={
              [
                catalogInfo.media_server_path?.tv_show_path,
                catalogInfo.media_server_path?.season_path,
              ]
                .filter(Boolean)
                .join('/') || ''
            }
            isCode={true}
          />
          <InfoItem
            icon={<Database size={18} />}
            label="Size of Media Files"
            value={catalogInfo.media_server_size_pretty ?? ''}
          />
          <InfoItem
            icon={<Scale size={18} />}
            label="Storage Type"
            value={catalogInfo.is_copy_files_in_media_server ? 'Copied' : 'Hardlinked'}
          />
        </div>
      </div>
    </div>
  );
});
