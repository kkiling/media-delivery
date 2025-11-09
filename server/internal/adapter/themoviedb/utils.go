package themoviedb

import (
	"fmt"
	"time"

	"github.com/samber/lo"
)

func mapImg(path, prefix string, urlImg string) *string {
	if urlImg == "" {
		return nil
	}
	return lo.ToPtr(fmt.Sprintf("%s/%s%s", urlImg, prefix, path))
}

func (api *API) getImage(path string) *Image {
	if path == "" {
		return nil
	}
	urlImg := api.baseImgUrl.String()
	return &Image{
		ID:       path,
		W92:      mapImg(path, "w92", urlImg),
		W154:     mapImg(path, "w154", urlImg),
		W185:     mapImg(path, "w185", urlImg),
		W342:     mapImg(path, "w342", urlImg),
		W500:     mapImg(path, "w500", urlImg),
		W780:     mapImg(path, "w780", urlImg),
		Original: *mapImg(path, "original", urlImg),
	}
}

func parseDate(dateStr string) time.Time {
	const layout = "2006-01-02"
	t, err := time.Parse(layout, dateStr)
	if err != nil {
		return time.Time{}
	}
	return t
}
