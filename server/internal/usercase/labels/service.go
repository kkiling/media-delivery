package labels

import (
	"context"
	"errors"
	"fmt"

	"github.com/kkiling/goplatform/storagebase"
	"github.com/kkiling/media-delivery/internal/common"
	ucerr "github.com/kkiling/media-delivery/internal/usercase/err"
)

type Service struct {
	storage Storage
}

func NewService(storage Storage) *Service {
	return &Service{
		storage: storage,
	}
}

func (s *Service) AddLabel(ctx context.Context, label Label) error {
	if err := s.storage.SaveLabel(ctx, label); err != nil {
		if errors.Is(err, storagebase.ErrAlreadyExists) {
			return fmt.Errorf("storage.SaveLabel: %w", ucerr.AlreadyExists)
		}
		return fmt.Errorf("storage.SaveLabel: %w", err)
	}

	return nil
}

func (s *Service) GetLabels(ctx context.Context, contentID common.ContentID) ([]Label, error) {
	return s.storage.GetLabels(ctx, contentID)
}

func (s *Service) DeleteLabel(ctx context.Context, contentID common.ContentID, typeLabel TypeLabel) error {
	if err := s.storage.DeleteLabel(ctx, contentID, typeLabel); err != nil {
		if errors.Is(err, storagebase.ErrNotFound) {
			return fmt.Errorf("storage.DeleteLabel: %w", ucerr.NotFound)
		}
		return fmt.Errorf("storage.DeleteLabel: %w", err)
	}

	return nil
}
