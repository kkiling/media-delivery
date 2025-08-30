import { ContentID, DeliveryStatus } from '@/api/api';
import { Card, Button, Spinner } from 'react-bootstrap';
import { PlusCircle } from 'react-bootstrap-icons';
import { observer } from 'mobx-react-lite';
import { useEffect } from 'react';
import { videoContentStore } from '@/stores/videoContentStore';
import { TVShowDeliveryContent } from './TVShowDeliveryContent';
import { DeliveredVideoContent } from './DeliveredVideoContent';
import Loading from '../Loading';

interface CreateVideoContentCardProps {
  loading: boolean;
  onCreateDelivery: () => void;
}

const CreateVideoContentCard = ({ loading, onCreateDelivery }: CreateVideoContentCardProps) => {
  return (
    <div className="text-center">
      <p className="text-muted mb-3">Not found</p>
      <Button
        variant="primary"
        onClick={onCreateDelivery}
        disabled={loading}
        className="d-inline-flex align-items-center gap-3"
      >
        {loading ? <Spinner animation="border" size="sm" /> : <PlusCircle size={24} />}
        <span>Create video delivery </span>
      </Button>
    </div>
  );
};

interface VideoContentProps {
  contentId: ContentID;
}

export const VideoContent = observer(({ contentId }: VideoContentProps) => {
  useEffect(() => {
    videoContentStore.fetchVideoContent(contentId);
  }, [contentId]);

  const onCreateDelivery = async () => {
    await videoContentStore.createVideoContent(contentId);
  };

  const onNeedReload = async () => {
    await videoContentStore.fetchVideoContent(contentId);
  };

  const { loading, error, content } = videoContentStore;

  const renderContent = () => {
    if (loading) {
      return <Loading />;
    }

    if (error) {
      return <p className="text-danger">{error}</p>;
    }

    if (!content) {
      return (
        <CreateVideoContentCard
          onCreateDelivery={onCreateDelivery}
          loading={videoContentStore.loading}
        />
      );
    }

    switch (content.delivery_status) {
      case DeliveryStatus.DeliveryStatusInProgress:
        return <TVShowDeliveryContent contentId={contentId} onNeedReload={onNeedReload} />;
      case DeliveryStatus.DeliveryStatusDelivered:
        return <DeliveredVideoContent contentId={contentId} />;
      default:
        return <p className="text-muted">No active delivery process</p>;
    }
  };

  return (
    <Card className="mb-4">
      <Card.Header as="h5" id="video-content-header">
        Video content
      </Card.Header>
      <Card.Body>{renderContent()}</Card.Body>
    </Card>
  );
});
