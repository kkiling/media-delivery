import { TorrentSearch } from '@/api/api';
import { SearchInput, SearchInputSize } from '../SearchInput';
import { TorrentCard } from './TorrentCard';
import { useState } from 'react';
import { Modal } from 'react-bootstrap';

interface TorrentSearchListProps {
  torrents?: TorrentSearch[];
  searchQuery?: string;
  onSearch: (query: string) => void;
  onSelect: (href: string) => void;
}

export const TorrentSearchList = ({
  torrents,
  searchQuery,
  onSearch,
  onSelect,
}: TorrentSearchListProps) => {
  const [selectedTorrent, setSelectedTorrent] = useState<TorrentSearch | null>(null);
  const [showModal, setShowModal] = useState(false);

  const handleTorrentClick = (torrent: TorrentSearch) => {
    setSelectedTorrent(torrent);
    setShowModal(true);
  };

  const handleConfirm = () => {
    if (selectedTorrent?.href) {
      onSelect(selectedTorrent.href);
      setShowModal(false);
    }
  };

  const handleClose = () => {
    setShowModal(false);
    setSelectedTorrent(null);
  };

  return (
    <>
      <div className="d-flex flex-column gap-3">
        <div className="text-center">
          <h4 className="mb-1">Select torrent</h4>
          <p className="text-muted small mb-0">Please select torrent or change search parameters</p>
        </div>
        <SearchInput
          placeholder="Search torrents..."
          initialQuery={searchQuery}
          onSubmit={onSearch}
          size={SearchInputSize.Medium}
        />

        <div className="d-flex flex-column">
          {torrents?.map((torrent, index) => (
            <TorrentCard
              key={index}
              torrent={torrent}
              onSelect={() => handleTorrentClick(torrent)}
              clickable={true}
              isLast={index === torrents.length - 1}
            />
          ))}
        </div>
      </div>

      <Modal show={showModal} onHide={handleClose} centered size="xl">
        <Modal.Header closeButton>
          <Modal.Title>Confirm torrent selection</Modal.Title>
        </Modal.Header>
        <Modal.Body>
          {selectedTorrent && (
            <TorrentCard torrent={selectedTorrent} clickable={false} isLast={true} />
          )}
        </Modal.Body>
        <Modal.Footer>
          <button className="btn btn-secondary" onClick={handleClose}>
            Cancel
          </button>
          <button className="btn btn-primary" onClick={handleConfirm}>
            Confirm
          </button>
        </Modal.Footer>
      </Modal>
    </>
  );
};
