import { Modal, Button } from 'react-bootstrap';

interface ConfirmModalProps {
  show: boolean;
  onClose: () => void;
  onSubmit: () => void;
}

export const ConfirmModal: React.FC<ConfirmModalProps> = ({ show, onClose, onSubmit }) => (
  <Modal show={show} onHide={onClose} centered>
    <Modal.Header closeButton>
      <Modal.Title>Confirm File Matching</Modal.Title>
    </Modal.Header>
    <Modal.Body>Are you sure you want to confirm the file matching?</Modal.Body>
    <Modal.Footer>
      <Button variant="secondary" onClick={onClose}>
        Cancel
      </Button>
      <Button variant="primary" onClick={onSubmit}>
        Confirm
      </Button>
    </Modal.Footer>
  </Modal>
);
