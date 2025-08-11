package qbittorrent

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/kkiling/torrent-to-media-server/internal/adapter/apierr"
)

func (api *Api) AddTorrent(opts TorrentAddOptions) error {
	if err := api.login(); err != nil {
		return fmt.Errorf("failed to login: %w", err)
	}

	form := url.Values{}
	form.Set("urls", opts.Magnet)
	form.Set("savepath", opts.SavePath)

	if opts.Category != "" {
		form.Set("category", opts.Category)
	}
	if len(opts.Tags) > 0 {
		form.Set("tags", strings.Join(opts.Tags, ","))
	}
	if opts.Paused {
		form.Set("paused", "true")
	}

	postUrl := api.baseAPIUrl.String() + "/api/v2/torrents/add"
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
