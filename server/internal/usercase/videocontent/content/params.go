package content

import (
	"github.com/kkiling/media-delivery/internal/common"
)

type CreateVideoContentParams struct {
	ContentID common.ContentID
}

type DeliveryVideoContentParams struct {
	ContentID common.ContentID
}

type DeleteVideoContentFilesParams struct {
	ContentID common.ContentID
}
