package delivery

import (
	"context"
	"fmt"
)

type GetMagnetLinkParams struct {
	Href string
}

// GetMagnetLink получение магнет ссылки на основе выбора раздачи пользователем
func (s *Service) GetMagnetLink(_ context.Context, params GetMagnetLinkParams) (*MagnetLink, error) {
	// Получение магнет ссылки
	magnetInfo, err := s.torrentSite.GetMagnetLink(params.Href)
	if err != nil {
		return nil, fmt.Errorf("torrentSite.GetMagnetLink: %w", err)
	}

	return &MagnetLink{
		Magnet: magnetInfo.Magnet,
		Hash:   magnetInfo.Hash,
	}, nil
}
