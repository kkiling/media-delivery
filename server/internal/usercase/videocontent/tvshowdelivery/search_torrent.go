package tvshowdelivery

import (
	"context"
	"fmt"
	"sort"
)

type SearchTorrentParams struct {
	SearchQuery string
}

// Функция для расчета рейтинга
func calculateScore(t TorrentSearch) float64 {
	// Веса можно настроить под ваши предпочтения
	seedWeight := 0.5
	downloadWeight := 0.3
	ratioWeight := 0.2

	ratio := float64(t.Seeds) / float64(t.Leeches+1)

	// Нормализуем значения (можно добавить нормализацию)
	return float64(t.Seeds)*seedWeight +
		float64(t.Downloads)*downloadWeight +
		ratio*ratioWeight
}

// SearchTorrent Делаем запрос к торрент сайту, получаем список раздач
func (s *Service) SearchTorrent(_ context.Context, params SearchTorrentParams) ([]TorrentSearch, error) {
	searchResult, err := s.torrentSite.SearchTorrents(params.SearchQuery)
	if err != nil {
		return nil, fmt.Errorf("torrentSite.SearchTorrents: %w", err)
	}

	var result []TorrentSearch
	for _, item := range searchResult.Results {
		result = append(result, TorrentSearch{
			Title:      item.Title,
			Href:       item.Href,
			SizeBytes:  item.SizeBytes,
			SizePretty: item.SizePretty,
			Seeds:      item.Seeds,
			Leeches:    item.Leeches,
			Downloads:  item.Downloads,
			AddedDate:  item.AddedDate,
			Category:   item.Category,
		})
	}

	// Сортировка
	sort.Slice(result, func(i, j int) bool {
		return calculateScore(result[i]) > calculateScore(result[j])
	})

	return result, nil
}
