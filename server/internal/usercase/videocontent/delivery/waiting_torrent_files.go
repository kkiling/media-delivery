package delivery

import (
	"context"
	"fmt"
	"path/filepath"
	"sort"

	"github.com/kkiling/media-delivery/internal/adapter/qbittorrent"
	ucerr "github.com/kkiling/media-delivery/internal/usercase/err"
)

type WaitingTorrentFilesParams struct {
	Hash string
}

// WaitingTorrentFiles ожидание когда появится информация о файлах в раздаче
func (s *Service) WaitingTorrentFiles(_ context.Context, params WaitingTorrentFilesParams) (*TorrentFilesData, error) {
	// Достаем инфу о торрент раздаче
	torrentInfo, err := s.torrentClient.GetTorrentInfo(params.Hash)
	if err != nil {
		return nil, fmt.Errorf("torrentClient.GetTorrentInfo: %w", err)
	}

	if torrentInfo == nil {
		return nil, fmt.Errorf("torrentInfo not found: %w", ucerr.NotFound)
	}

	switch torrentInfo.State {
	case qbittorrent.TorrentStatePausedDL, qbittorrent.TorrentStateStoppedDL:
		if err = s.torrentClient.ResumeTorrent(params.Hash); err != nil {
			return nil, fmt.Errorf("torrentClient.ResumeTorrent: %w", err)
		}
	case
		qbittorrent.TorrentStateQueuedUP,
		qbittorrent.TorrentStateDownloading,
		qbittorrent.TorrentStateUploading,
		qbittorrent.TorrentStatePausedUP,
		qbittorrent.TorrentStateStalledUP:
		// Файлы начали скачиваться, значит можем получить информацию о файлах
	default:
		// Ошибки как таковой нет, придем в следующий раз
		return nil, nil
	}

	torrentFiles, err := s.torrentClient.GetTorrentFiles(params.Hash)
	if err != nil {
		return nil, fmt.Errorf("torrentClient.GetTorrentInfo: %w", err)
	}

	if len(torrentFiles) == 0 {
		return nil, fmt.Errorf("torrentFiles not found: %w", ucerr.NotFound)
	}

	//
	fullPath := filepath.Join(s.config.BasePath, torrentInfo.ContentPath)
	// Вычисляем относительный путь от savePath до currentPath
	// SavePath: /downloads
	// ContentPath /downloads/Vinland Sag
	// relPath Vinland Sag
	relPath, err := filepath.Rel(torrentInfo.SavePath, torrentInfo.ContentPath)
	if err != nil {
		return nil, fmt.Errorf("filepath.Rel: %w", err)
	}

	var prepareTorrentFiles []FileInfo
	for _, file := range torrentFiles {
		relFile, err := filepath.Rel(relPath, file.Name)
		if err != nil {
			return nil, fmt.Errorf("filepath.Rel: %w", err)
		}
		prepareTorrentFiles = append(prepareTorrentFiles, FileInfo{
			RelativePath: relFile,
			FullPath:     filepath.Join(fullPath, relFile),
			//Extension:    strings.ToLower(filepath.Ext(relFile)),
			//Size:         file.Size,
		})
	}

	sort.Slice(prepareTorrentFiles, func(i, j int) bool {
		return prepareTorrentFiles[i].FullPath < prepareTorrentFiles[j].FullPath
	})

	return &TorrentFilesData{
		ContentFullPath: fullPath,
		Files:           prepareTorrentFiles,
	}, nil
}
