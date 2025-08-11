package qbittorrent

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/kkiling/torrent-to-media-server/internal/adapter/apierr"
)

type rawTorrentFile struct {
	Index    int     `json:"index"`
	Name     string  `json:"name"`
	Progress float64 `json:"progress"`
	Size     int64   `json:"size"`
}

// GetTorrentFiles возвращает список файлов в торренте по его хешу
func (api *Api) GetTorrentFiles(hash string) ([]TorrentFile, error) {
	if err := api.login(); err != nil {
		return nil, fmt.Errorf("failed to login: %w", err)
	}

	query := url.Values{}
	query.Add("hash", hash)

	reqURL := api.baseAPIUrl.String() + "/api/v2/torrents/files?" + query.Encode()
	resp, err := api.httpClient.Get(reqURL)
	if err != nil {
		return nil, apierr.HandleStatusCodeError(api.logger, resp)
	}
	defer resp.Body.Close()

	if errStatus := api.handleStatusError(resp); errStatus != nil {
		return nil, errStatus
	}

	var rawFiles []rawTorrentFile
	if err := json.NewDecoder(resp.Body).Decode(&rawFiles); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	files := make([]TorrentFile, len(rawFiles))
	for i, raw := range rawFiles {
		files[i] = TorrentFile{
			Index:    raw.Index,
			Name:     raw.Name,
			Progress: raw.Progress,
			Size:     raw.Size,
		}
	}

	return files, nil
}
