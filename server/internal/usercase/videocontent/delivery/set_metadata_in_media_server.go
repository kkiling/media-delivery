package delivery

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/kkiling/torrent-to-media-server/internal/adapter/emby"
	"github.com/kkiling/torrent-to-media-server/internal/usercase/videocontent/common"
)

type SetMediaMetaDataParams struct {
	TVShowPath string
	TVShowID   common.TVShowID
}

// SetMediaMetaData установка методанных
func (s *Service) SetMediaMetaData(ctx context.Context, params SetMediaMetaDataParams) error {
	tvShowPath, err := filepath.Rel(s.config.BasePath, params.TVShowPath)
	if err != nil {
		return fmt.Errorf("failed to get relative path: %w", err)
	}

	if err := s.embyApi.Refresh(); err != nil {
		return fmt.Errorf("failed to refresh emby api: %w", err)
	}

	info, err := s.embyApi.GetCatalogInfo("/" + tvShowPath)
	if err != nil {
		return fmt.Errorf("embyApi.GetCatalogInfo: %w", err)
	}
	if info == nil {
		return fmt.Errorf("catalogInfo: info is nil")
	}

	if !info.IsFolder {
		return fmt.Errorf("catalogInfo: info is not folder")
	}
	if info.Type != emby.SeriesTypeCatalog {
		return fmt.Errorf("catalogInfo: type is not series")
	}

	err = s.embyApi.ResetMetadata(info.ID)
	if err != nil {
		return fmt.Errorf("embyApi.ResetMetadata: %w", err)
	}

	err = s.embyApi.RemoteSearchApply(info.ID, params.TVShowID.ID)
	if err != nil {
		return fmt.Errorf("embyApi.RemoteSearchApply: %w", err)
	}

	return nil
}
