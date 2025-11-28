package tvshowdelivery

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/kkiling/media-delivery/internal/adapter/emby"
	"github.com/kkiling/media-delivery/internal/common"
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

	if err = s.embyApi.Refresh(); err != nil {
		return fmt.Errorf("failed to refresh emby api: %w", err)
	}

	info, err := s.embyApi.GetCatalogInfo("/" + tvShowPath)
	if err != nil {
		return fmt.Errorf("embyApi.GetCatalogInfo: %w", err)
	}

	if info == nil {
		return fmt.Errorf("catalogInfo: info is nil")
	}
	if info.TheMovieDbID == params.TVShowID.ID {
		// Сериал уже правильно идентифицирован
		return nil
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

	// Запрашиваем еще раз инфу о каталоге и сравниваем TheMovieDbId что бы удостоверитсья что мы установили метадату
	info, err = s.embyApi.GetCatalogInfo("/" + tvShowPath)
	if err != nil {
		return fmt.Errorf("embyApi.GetCatalogInfo: %w", err)
	}
	if info == nil {
		return fmt.Errorf("catalogInfo: info is nil")
	}
	if info.TheMovieDbID != params.TVShowID.ID {
		return fmt.Errorf("TheMovieDbID does not match")
	}

	return nil
}
