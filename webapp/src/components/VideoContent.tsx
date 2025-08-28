import { ContentID, DeliveryStatus, TVShowDeliveryStatus } from '@/api/api';
import { Card, Button, Modal } from 'react-bootstrap';
import { PlusCircle } from 'react-bootstrap-icons';
import { observer } from 'mobx-react-lite';
import { useEffect, useState } from 'react';
import { videoContentStore } from '@/stores/videoContentStore';
import { tvShowDeliveryStore } from '@/stores/tvShowDeliveryStore';
import { SearchInput } from './SearchInput';
import { Table } from 'react-bootstrap';

interface CreateVideoContentCardProps {
  onCreateDelivery: () => void;
}

const CreateVideoContentCard = ({ onCreateDelivery }: CreateVideoContentCardProps) => {
  return (
    <div className="text-center">
      <p className="text-muted mb-3">Not found</p>
      <Button
        variant="primary"
        onClick={onCreateDelivery}
        className="d-inline-flex align-items-center gap-3"
      >
        <PlusCircle size={24} />
        Create video delivery
      </Button>
    </div>
  );
};

interface TVShowDeliveryContentProps {
  contentId: ContentID;
}

const TVShowDeliveryContent = observer(({ contentId }: TVShowDeliveryContentProps) => {
  useEffect(() => {
    tvShowDeliveryStore.fetchDeliveryData(contentId);
  }, [contentId]);

  const onSearchSubmit = (query: string) => {
    tvShowDeliveryStore.selectTorrent(contentId, undefined, query);
  };

  const onTorrentSelect = (href: string) => {
    tvShowDeliveryStore.selectTorrent(contentId, href);
  };

  const { loading, error, deliveryState } = tvShowDeliveryStore;

  if (loading) {
    return <p>Loading...</p>;
  }

  if (error) {
    return <p className="text-danger">{error}</p>;
  }

  if (!deliveryState) {
    return <p className="text-muted">No delivery data available</p>;
  }

  if (deliveryState.step === TVShowDeliveryStatus.WaitingUserChoseTorrent) {
    return (
      <div className="d-flex flex-column gap-3">
        <SearchInput
          placeholder="Search torrents..."
          initialQuery={deliveryState.data?.search_query?.Query}
          onSubmit={onSearchSubmit}
        />

        <div className="table-responsive">
          <Table hover>
            <thead>
              <tr>
                <th>Title</th>
                <th>Size</th>
                <th>Seeds</th>
                <th>Leeches</th>
                <th>Downloads</th>
                <th>Added Date</th>
                <th>Action</th>
              </tr>
            </thead>
            <tbody>
              {deliveryState.data?.torrent_search?.map((torrent, index) => (
                <tr key={index}>
                  <td>{torrent.title}</td>
                  <td>{torrent.size}</td>
                  <td className="text-success">{torrent.seeds}</td>
                  <td className="text-danger">{torrent.leeches}</td>
                  <td>{torrent.downloads}</td>
                  <td>{torrent.added_date}</td>
                  <td>
                    <Button
                      size="sm"
                      variant="primary"
                      onClick={() => torrent.href && onTorrentSelect(torrent.href)}
                    >
                      Select
                    </Button>
                  </td>
                </tr>
              ))}
            </tbody>
          </Table>
        </div>
      </div>
    );
  }

  return <p>{deliveryState.step}</p>;
});

interface VideoContentProps {
  contentId: ContentID;
}

export const VideoContent = observer(({ contentId }: VideoContentProps) => {
  useEffect(() => {
    videoContentStore.fetchVideoContent(contentId);
  }, [contentId]);

  const onCreateDelivery = async () => {
    try {
      await videoContentStore.createVideoContent(contentId);
    } catch (error) {
      console.error('Failed to create delivery:', error);
    }
  };

  const { loading, error, content } = videoContentStore;

  const renderContent = () => {
    if (loading) {
      return <p>Loading...</p>;
    }

    if (error) {
      return <p className="text-danger">{error}</p>;
    }

    if (!content) {
      return <CreateVideoContentCard onCreateDelivery={onCreateDelivery} />;
    }

    return (
      <>
        <div>{content.delivery_status}</div>
        {content.delivery_status === DeliveryStatus.DeliveryStatusInProgress && (
          <TVShowDeliveryContent contentId={contentId} />
        )}
      </>
    );
  };

  return (
    <Card className="mb-4">
      <Card.Header as="h5">Video content</Card.Header>
      <Card.Body>{renderContent()}</Card.Body>
    </Card>
  );
});
