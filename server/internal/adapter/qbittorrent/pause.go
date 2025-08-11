package qbittorrent

import (
	"fmt"
	"net/url"

	"github.com/kkiling/torrent-to-media-server/internal/adapter/apierr"
)

func (api *Api) PauseTorrent(hash string) error {
	if err := api.login(); err != nil {
		return fmt.Errorf("failed to login: %w", err)
	}

	form := url.Values{}
	form.Set("hashes", hash)

	postUrl := api.baseAPIUrl.String() + "/api/v2/torrents/stop"
	resp, err := api.httpClient.PostForm(postUrl, form)
	if err != nil {
		return apierr.HandleStatusCodeError(api.logger, resp)
	}
	defer resp.Body.Close()

	if errStatus := api.handleStatusError(resp); errStatus != nil {
		return errStatus
	}

	return nil
}
