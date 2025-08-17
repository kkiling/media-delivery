package emby

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/kkiling/media-delivery/internal/adapter/apierr"
)

// ResetMetadata сбрасываем всю методату
func (api *API) ResetMetadata(embyID uint64) error {
	queryParams := url.Values{}
	queryParams.Add("api_key", api.apiKey)
	queryParams.Add("ItemIds", strconv.FormatUint(embyID, 10))

	url := fmt.Sprintf("%s/Items/Metadata/Reset?%s", api.baseAPIUrl.String(), queryParams.Encode())

	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Устанавливаем заголовки как в curl
	req.Header.Set("accept", "*/*")
	req.Header.Set("Content-Type", "application/json")

	resp, err := api.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to emby items: %w", apierr.HandleRequestError(api.logger, err))
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return apierr.HandleStatusCodeError(api.logger, resp)
	}

	return nil
}
