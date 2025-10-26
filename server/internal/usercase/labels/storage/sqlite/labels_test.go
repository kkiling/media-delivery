package sqlite

import (
	"context"
	"math/rand"
	"testing"
	"time"

	"github.com/kkiling/goplatform/storagebase/testutils"
	"github.com/kkiling/media-delivery/internal/common"
	"github.com/kkiling/media-delivery/internal/usercase/labels"
	"github.com/stretchr/testify/require"
)

func TestStorage_SaveLabel(t *testing.T) {
	t.Parallel()
	rand.Seed(time.Now().UnixNano())

	testStorage := NewTestStorage(testutils.SetupSqlTestDB(t))
	ctx := context.Background()

	t.Run("ok tvshow", func(t *testing.T) {
		t.Parallel()
		label := labels.Label{
			ContentID: common.ContentID{
				TVShow: &common.TVShowID{
					ID:           rand.Uint64(),
					SeasonNumber: 2,
				},
			},
			TypeLabel: labels.HasVideoContent,
			CreatedAt: time.Now(),
		}
		err := testStorage.SaveLabel(ctx, label)
		require.NoError(t, err)
	})
}
