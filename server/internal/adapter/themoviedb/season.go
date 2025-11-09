package themoviedb

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/kkiling/media-delivery/internal/adapter/apierr"
)

func (api *API) GetSeason(_ context.Context, tvID uint64, seasonNumber uint8, language Language) (*SeasonWithEpisodes, error) {
	queryParams := url.Values{}
	queryParams.Add("api_key", api.apiKey)
	queryParams.Add("language", string(language))

	getUrl := fmt.Sprintf("%s/tv/%d/season/%d?%s", api.baseAPIUrl.String(), tvID, seasonNumber, queryParams.Encode())
	resp, err := api.httpClient.Get(getUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to get season episodes: %w", apierr.HandleRequestError(api.logger, err))
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, apierr.HandleStatusCodeError(api.logger, resp)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var result struct {
		AirDate      string  `json:"air_date"`
		ID           uint64  `json:"id"`
		Name         string  `json:"name"`
		Overview     string  `json:"overview"`
		Poster       string  `json:"poster_path"`
		SeasonNumber uint8   `json:"season_number"`
		VoteAverage  float32 `json:"vote_average"`
		Episodes     []struct {
			AirDate       string  `json:"air_date"`
			EpisodeNumber int     `json:"episode_number"`
			EpisodeType   *string `json:"episode_type,omitempty"`
			ID            uint64  `json:"id"`
			Name          string  `json:"name"`
			Overview      string  `json:"overview"`
			Runtime       uint32  `json:"runtime"`
			Still         string  `json:"still_path"`
			VoteAverage   float32 `json:"vote_average"`
			VoteCount     uint32  `json:"vote_count"`
		} `json:"episodes"`
	}

	if err = json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	season := SeasonWithEpisodes{
		Season: Season{
			AirDate:      parseDate(result.AirDate),
			EpisodeCount: uint32(len(result.Episodes)),
			ID:           result.ID,
			Name:         result.Name,
			Overview:     result.Overview,
			Poster:       api.getImage(result.Poster),
			SeasonNumber: result.SeasonNumber,
			VoteAverage:  result.VoteAverage,
		},
		Episodes: make([]Episode, 0, len(result.Episodes)),
	}
	for _, ep := range result.Episodes {
		season.Episodes = append(season.Episodes, Episode{
			AirDate:       parseDate(ep.AirDate),
			EpisodeNumber: ep.EpisodeNumber,
			EpisodeType:   ep.EpisodeType,
			ID:            ep.ID,
			Name:          ep.Name,
			Overview:      ep.Overview,
			Runtime:       ep.Runtime,
			Still:         api.getImage(ep.Still),
			VoteAverage:   ep.VoteAverage,
			VoteCount:     ep.VoteCount,
		})
	}

	return &season, nil
}
