package runners

import (
	"github.com/kkiling/media-delivery/internal/common"
)

type Type string

const (
	TVShowDelivery Type = "tv_show_delivery"
	TVShowDelete   Type = "tv_show_delete"
)

type Metadata struct {
	ContentID common.ContentID
}

type FailData struct {
}
