package postgresql

import (
	"context"
	"math/rand"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/kkiling/goplatform/storagebase"
	"github.com/kkiling/goplatform/storagebase/testutils"
	"github.com/samber/lo"
	"github.com/stretchr/testify/require"

	"github.com/kkiling/media-delivery/internal/common"
	"github.com/kkiling/media-delivery/internal/usercase/videocontent/content"
	"github.com/kkiling/media-delivery/internal/usercase/videocontent/runners"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func randID() uint64 {
	return uint64(rand.Uint32())
}

func equalVideoContent(t *testing.T, vc1, vc2 *content.VideoContent) {
	require.True(t, (vc1 == nil && vc2 == nil) || (vc1 != nil && vc2 != nil))
	if vc1 == nil || vc2 == nil {
		return
	}

	require.Equal(t, vc1.ID, vc2.ID)

	// Сравниваем ContentID
	require.Equal(t, vc1.ContentID.MovieID, vc2.ContentID.MovieID)
	if vc1.ContentID.TVShow != nil && vc2.ContentID.TVShow != nil {
		require.Equal(t, vc1.ContentID.TVShow.ID, vc2.ContentID.TVShow.ID)
		require.Equal(t, vc1.ContentID.TVShow.SeasonNumber, vc2.ContentID.TVShow.SeasonNumber)
	} else {
		require.True(t, vc1.ContentID.TVShow == nil && vc2.ContentID.TVShow == nil)
	}

	require.WithinDuration(t, vc1.CreatedAt, vc2.CreatedAt, time.Second)
	require.Equal(t, vc1.DeliveryStatus, vc2.DeliveryStatus)
	require.Equal(t, vc1.State, vc2.State)
}

func equalVideoContents(t *testing.T, contents1, contents2 []content.VideoContent) {
	require.Equal(t, len(contents1), len(contents2))

	// Создаем мапы для сравнения без учета порядка
	contents1Map := make(map[uuid.UUID]content.VideoContent)
	contents2Map := make(map[uuid.UUID]content.VideoContent)

	for _, vc := range contents1 {
		contents1Map[vc.ID] = vc
	}
	for _, vc := range contents2 {
		contents2Map[vc.ID] = vc
	}

	for id, vc1 := range contents1Map {
		vc2, exists := contents2Map[id]
		require.True(t, exists, "video content with ID %s not found", id)
		equalVideoContent(t, &vc1, &vc2)
	}
}

func TestStorage_SaveVideoContent(t *testing.T) {
	t.Parallel()

	testStorage := NewTestStorage(testutils.SetupPostgresqlTestDB(t))
	ctx := context.Background()

	t.Run("save video content - successfully with movie ID", func(t *testing.T) {
		t.Parallel()

		movieID := randID()
		videoContent := &content.VideoContent{
			ID: uuid.New(),
			ContentID: common.ContentID{
				MovieID: &movieID,
				TVShow:  nil,
			},
			CreatedAt:      time.Now(),
			DeliveryStatus: content.DeliveryStatusInProgress,
			State: []content.State{
				{
					StateID: uuid.New(),
					Type:    runners.TVShowDelivery,
				},
				{
					StateID: uuid.New(),
					Type:    runners.TVShowDelivery,
				},
			},
		}

		err := testStorage.SaveVideoContent(ctx, videoContent)
		require.NoError(t, err)

		// Проверяем через GetVideoContents
		savedContents, err := testStorage.GetVideoContents(ctx, videoContent.ContentID)
		require.NoError(t, err)
		require.Len(t, savedContents, 1)

		equalVideoContent(t, videoContent, &savedContents[0])
	})

	t.Run("save video content - successfully with TV show ID", func(t *testing.T) {
		t.Parallel()

		videoContent := &content.VideoContent{
			ID: uuid.New(),
			ContentID: common.ContentID{
				MovieID: nil,
				TVShow: &common.TVShowID{
					ID:           randID(),
					SeasonNumber: 1,
				},
			},
			CreatedAt:      time.Now(),
			DeliveryStatus: content.DeliveryStatusInProgress,
			State: []content.State{
				{
					StateID: uuid.New(),
					Type:    runners.TVShowDelivery,
				},
			},
		}

		err := testStorage.SaveVideoContent(ctx, videoContent)
		require.NoError(t, err)

		// Проверяем через GetVideoContents
		savedContents, err := testStorage.GetVideoContents(ctx, videoContent.ContentID)
		require.NoError(t, err)
		require.Len(t, savedContents, 1)

		equalVideoContent(t, videoContent, &savedContents[0])
	})

	t.Run("save video content - already exists", func(t *testing.T) {
		t.Parallel()

		movieID := randID()
		videoContent := &content.VideoContent{
			ID: uuid.New(),
			ContentID: common.ContentID{
				MovieID: &movieID,
			},
			CreatedAt:      time.Now(),
			DeliveryStatus: content.DeliveryStatusInProgress,
			State:          []content.State{},
		}

		// Первое сохранение
		err := testStorage.SaveVideoContent(ctx, videoContent)
		require.NoError(t, err)

		// Второе сохранение с тем же ID
		err = testStorage.SaveVideoContent(ctx, videoContent)
		require.Error(t, err)
		require.ErrorIs(t, storagebase.ErrAlreadyExists, err)
	})

	t.Run("save multiple video contents for same content ID", func(t *testing.T) {
		t.Parallel()

		movieID := randID()
		contentID := common.ContentID{MovieID: &movieID}

		videoContents := []*content.VideoContent{
			{
				ID:             uuid.New(),
				ContentID:      contentID,
				CreatedAt:      time.Now(),
				DeliveryStatus: content.DeliveryStatusUpdating,
				State: []content.State{
					{StateID: uuid.New(), Type: runners.TVShowDelivery},
				},
			},
			{
				ID:             uuid.New(),
				ContentID:      contentID,
				CreatedAt:      time.Now().Add(time.Hour),
				DeliveryStatus: content.DeliveryStatusInProgress,
				State: []content.State{
					{StateID: uuid.New(), Type: runners.TVShowDelivery},
				},
			},
		}

		// Сохраняем все видео контенты
		for _, vc := range videoContents {
			err := testStorage.SaveVideoContent(ctx, vc)
			require.NoError(t, err)
		}

		// Проверяем через GetVideoContents
		savedContents, err := testStorage.GetVideoContents(ctx, contentID)
		require.NoError(t, err)

		// Конвертируем []*VideoContent в []VideoContent для сравнения
		expectedContents := make([]content.VideoContent, len(videoContents))
		for i, vc := range videoContents {
			expectedContents[i] = *vc
		}

		equalVideoContents(t, expectedContents, savedContents)
	})

	t.Run("save video content - empty state", func(t *testing.T) {
		t.Parallel()

		movieID := randID()
		videoContent := &content.VideoContent{
			ID: uuid.New(),
			ContentID: common.ContentID{
				MovieID: &movieID,
			},
			CreatedAt:      time.Now(),
			DeliveryStatus: content.DeliveryStatusDelivered,
			State:          []content.State{},
		}

		err := testStorage.SaveVideoContent(ctx, videoContent)
		require.NoError(t, err)

		// Проверяем через GetVideoContents
		savedContents, err := testStorage.GetVideoContents(ctx, videoContent.ContentID)
		require.NoError(t, err)
		require.Len(t, savedContents, 1)

		equalVideoContent(t, videoContent, &savedContents[0])
	})
}

func TestStorage_UpdateVideoContent(t *testing.T) {
	t.Parallel()

	testStorage := NewTestStorage(testutils.SetupPostgresqlTestDB(t))
	ctx := context.Background()

	t.Run("update video content - successfully update delivery status", func(t *testing.T) {
		t.Parallel()

		// Сначала создаем видео контент
		movieID := randID()
		videoContent := &content.VideoContent{
			ID: uuid.New(),
			ContentID: common.ContentID{
				MovieID: &movieID,
			},
			CreatedAt:      time.Now(),
			DeliveryStatus: content.DeliveryStatusInProgress,
			State: []content.State{
				{
					StateID: uuid.New(),
					Type:    runners.TVShowDelivery,
				},
			},
		}

		err := testStorage.SaveVideoContent(ctx, videoContent)
		require.NoError(t, err)

		// Обновляем статус
		updateData := &content.UpdateVideoContent{
			DeliveryStatus: content.DeliveryStatusDelivered,
		}

		err = testStorage.UpdateVideoContent(ctx, videoContent.ID, updateData)
		require.NoError(t, err)

		// Проверяем что статус обновился
		savedContents, err := testStorage.GetVideoContents(ctx, videoContent.ContentID)
		require.NoError(t, err)
		require.Len(t, savedContents, 1)

		require.Equal(t, content.DeliveryStatusDelivered, savedContents[0].DeliveryStatus)
		// Проверяем что остальные поля не изменились
		require.Equal(t, videoContent.ID, savedContents[0].ID)
		require.Equal(t, videoContent.ContentID.MovieID, savedContents[0].ContentID.MovieID)
		require.WithinDuration(t, videoContent.CreatedAt, savedContents[0].CreatedAt, time.Second)
		require.Equal(t, videoContent.State, savedContents[0].State)
	})

	t.Run("update video content - not found", func(t *testing.T) {
		t.Parallel()

		nonExistentID := uuid.New()
		updateData := &content.UpdateVideoContent{
			DeliveryStatus: content.DeliveryStatusDelivered,
		}

		err := testStorage.UpdateVideoContent(ctx, nonExistentID, updateData)
		require.Error(t, err)
		require.ErrorIs(t, storagebase.ErrNotFound, err)
	})

	t.Run("update video content - multiple status changes", func(t *testing.T) {
		t.Parallel()

		// Создаем видео контент
		movieID := randID()
		videoContent := &content.VideoContent{
			ID: uuid.New(),
			ContentID: common.ContentID{
				MovieID: &movieID,
			},
			CreatedAt:      time.Now(),
			DeliveryStatus: content.DeliveryStatusUpdating,
			State:          []content.State{},
		}

		err := testStorage.SaveVideoContent(ctx, videoContent)
		require.NoError(t, err)

		// Последовательные обновления статуса
		statusUpdates := []content.DeliveryStatus{
			content.DeliveryStatusInProgress,
			content.DeliveryStatusDelivered,
			content.DeliveryStatusFailed,
		}

		for _, newStatus := range statusUpdates {
			updateData := &content.UpdateVideoContent{
				DeliveryStatus: newStatus,
			}

			err = testStorage.UpdateVideoContent(ctx, videoContent.ID, updateData)
			require.NoError(t, err)

			// Проверяем что статус обновился
			savedContents, err := testStorage.GetVideoContents(ctx, videoContent.ContentID)
			require.NoError(t, err)
			require.Len(t, savedContents, 1)
			require.Equal(t, newStatus, savedContents[0].DeliveryStatus)
		}
	})

	t.Run("update video content - same status", func(t *testing.T) {
		t.Parallel()

		movieID := randID()
		videoContent := &content.VideoContent{
			ID: uuid.New(),
			ContentID: common.ContentID{
				MovieID: &movieID,
			},
			CreatedAt:      time.Now(),
			DeliveryStatus: content.DeliveryStatusInProgress,
			State:          []content.State{},
		}

		err := testStorage.SaveVideoContent(ctx, videoContent)
		require.NoError(t, err)

		// Обновляем на тот же статус
		updateData := &content.UpdateVideoContent{
			DeliveryStatus: content.DeliveryStatusInProgress,
		}

		err = testStorage.UpdateVideoContent(ctx, videoContent.ID, updateData)
		require.NoError(t, err)

		// Проверяем что данные не сломались
		savedContents, err := testStorage.GetVideoContents(ctx, videoContent.ContentID)
		require.NoError(t, err)
		require.Len(t, savedContents, 1)
		equalVideoContent(t, videoContent, &savedContents[0])
	})
}

func TestStorage_GetVideoContents(t *testing.T) {
	t.Parallel()

	testStorage := NewTestStorage(testutils.SetupPostgresqlTestDB(t))
	ctx := context.Background()

	t.Run("get video contents - movie ID not found", func(t *testing.T) {
		t.Parallel()

		nonExistentMovieID := randID()
		contentID := common.ContentID{
			MovieID: &nonExistentMovieID,
		}

		contents, err := testStorage.GetVideoContents(ctx, contentID)
		require.NoError(t, err)
		require.Empty(t, contents)
	})

	t.Run("get video contents - TV show ID not found", func(t *testing.T) {
		t.Parallel()

		contentID := common.ContentID{
			TVShow: &common.TVShowID{
				ID:           randID(),
				SeasonNumber: 1,
			},
		}

		contents, err := testStorage.GetVideoContents(ctx, contentID)
		require.NoError(t, err)
		require.Empty(t, contents)
	})

	t.Run("get video contents - invalid content ID", func(t *testing.T) {
		t.Parallel()

		contentID := common.ContentID{
			MovieID: nil,
			TVShow:  nil,
		}

		contents, err := testStorage.GetVideoContents(ctx, contentID)
		require.Error(t, err)
		require.Nil(t, contents)
		require.Contains(t, err.Error(), "contentID is not valid")
	})

	t.Run("get video contents - multiple contents for same movie", func(t *testing.T) {
		t.Parallel()

		movieID := randID()
		contentID := common.ContentID{MovieID: &movieID}

		videoContents := []*content.VideoContent{
			{
				ID:             uuid.New(),
				ContentID:      contentID,
				CreatedAt:      time.Now(),
				DeliveryStatus: content.DeliveryStatusUpdating,
				State: []content.State{
					{StateID: uuid.New(), Type: runners.TVShowDelivery},
				},
			},
			{
				ID:             uuid.New(),
				ContentID:      contentID,
				CreatedAt:      time.Now().Add(time.Hour),
				DeliveryStatus: content.DeliveryStatusInProgress,
				State: []content.State{
					{StateID: uuid.New(), Type: runners.TVShowDelivery},
					{StateID: uuid.New(), Type: runners.TVShowDelivery},
				},
			},
		}

		// Сохраняем все видео контенты
		for _, vc := range videoContents {
			err := testStorage.SaveVideoContent(ctx, vc)
			require.NoError(t, err)
		}

		// Получаем все контенты для этого movie
		savedContents, err := testStorage.GetVideoContents(ctx, contentID)
		require.NoError(t, err)
		require.Len(t, savedContents, 2)

		// Конвертируем []*VideoContent в []VideoContent для сравнения
		expectedContents := make([]content.VideoContent, len(videoContents))
		for i, vc := range videoContents {
			expectedContents[i] = *vc
		}

		equalVideoContents(t, expectedContents, savedContents)
	})

	t.Run("get video contents - multiple contents for same TV show", func(t *testing.T) {
		t.Parallel()

		tvShowID := randID()
		seasonNumber := uint8(2)
		contentID := common.ContentID{
			TVShow: &common.TVShowID{
				ID:           tvShowID,
				SeasonNumber: seasonNumber,
			},
		}

		videoContents := []*content.VideoContent{
			{
				ID:             uuid.New(),
				ContentID:      contentID,
				CreatedAt:      time.Now(),
				DeliveryStatus: content.DeliveryStatusInProgress,
				State: []content.State{
					{StateID: uuid.New(), Type: runners.TVShowDelivery},
				},
			},
			{
				ID:             uuid.New(),
				ContentID:      contentID,
				CreatedAt:      time.Now().Add(2 * time.Hour),
				DeliveryStatus: content.DeliveryStatusDelivered,
				State:          []content.State{},
			},
		}

		// Сохраняем все видео контенты
		for _, vc := range videoContents {
			err := testStorage.SaveVideoContent(ctx, vc)
			require.NoError(t, err)
		}

		// Получаем все контенты для этого TV show
		savedContents, err := testStorage.GetVideoContents(ctx, contentID)
		require.NoError(t, err)
		require.Len(t, savedContents, 2)

		// Конвертируем []*VideoContent в []VideoContent для сравнения
		expectedContents := make([]content.VideoContent, len(videoContents))
		for i, vc := range videoContents {
			expectedContents[i] = *vc
		}

		equalVideoContents(t, expectedContents, savedContents)
	})

	t.Run("get video contents - filter by specific content ID", func(t *testing.T) {
		t.Parallel()

		movieID1 := randID()
		movieID2 := randID()

		// Создаем контенты для разных movie
		content1 := &content.VideoContent{
			ID: uuid.New(),
			ContentID: common.ContentID{
				MovieID: &movieID1,
			},
			CreatedAt:      time.Now(),
			DeliveryStatus: content.DeliveryStatusInProgress,
			State:          []content.State{},
		}

		content2 := &content.VideoContent{
			ID: uuid.New(),
			ContentID: common.ContentID{
				MovieID: &movieID2,
			},
			CreatedAt:      time.Now(),
			DeliveryStatus: content.DeliveryStatusDelivered,
			State:          []content.State{},
		}

		err := testStorage.SaveVideoContent(ctx, content1)
		require.NoError(t, err)
		err = testStorage.SaveVideoContent(ctx, content2)
		require.NoError(t, err)

		// Получаем только контенты для movieID1
		contentID1 := common.ContentID{MovieID: &movieID1}
		contents1, err := testStorage.GetVideoContents(ctx, contentID1)
		require.NoError(t, err)
		require.Len(t, contents1, 1)
		equalVideoContent(t, content1, &contents1[0])

		// Получаем только контенты для movieID2
		contentID2 := common.ContentID{MovieID: &movieID2}
		contents2, err := testStorage.GetVideoContents(ctx, contentID2)
		require.NoError(t, err)
		require.Len(t, contents2, 1)
		equalVideoContent(t, content2, &contents2[0])
	})
}

func TestStorage_GetVideoContentsByDeliveryStatus(t *testing.T) {
	t.Parallel()

	testStorage := NewTestStorage(testutils.SetupPostgresqlTestDB(t))
	ctx := context.Background()
	const limit = 10000 // Что бы тесты не падали после ретраев

	t.Run("get video contents by delivery status - with other records and update", func(t *testing.T) {
		t.Parallel()

		// Создаем несколько контентов с разными статусами
		targetStatus := content.DeliveryStatusInProgress
		otherStatus := content.DeliveryStatusDelivered

		targetContents := []*content.VideoContent{
			{
				ID: uuid.New(),
				ContentID: common.ContentID{
					MovieID: lo.ToPtr(randID()),
				},
				CreatedAt:      time.Now(),
				DeliveryStatus: targetStatus,
				State: []content.State{
					{StateID: uuid.New(), Type: runners.TVShowDelivery},
				},
			},
			{
				ID: uuid.New(),
				ContentID: common.ContentID{
					TVShow: &common.TVShowID{
						ID:           randID(),
						SeasonNumber: 1,
					},
				},
				CreatedAt:      time.Now().Add(time.Hour),
				DeliveryStatus: targetStatus,
				State:          []content.State{},
			},
		}

		for _, vc := range targetContents {
			err := testStorage.SaveVideoContent(ctx, vc)
			require.NoError(t, err)
		}

		// Получаем контенты с целевым статусом
		initialContents, err := testStorage.GetVideoContentsByDeliveryStatus(ctx, targetStatus, limit)
		require.NoError(t, err)

		// Ищем наши целевые контенты среди полученных
		foundTargets := make(map[uuid.UUID]bool)
		for _, item := range initialContents {
			for _, target := range targetContents {
				if item.ID == target.ID {
					foundTargets[item.ID] = true
					equalVideoContent(t, target, &item)
					break
				}
			}
		}

		// Проверяем что нашли все целевые контенты
		require.Equal(t, len(targetContents), len(foundTargets))

		// Обновляем статус одного из целевых контентов
		contentToUpdate := targetContents[0]
		updateData := &content.UpdateVideoContent{
			DeliveryStatus: otherStatus,
		}

		err = testStorage.UpdateVideoContent(ctx, contentToUpdate.ID, updateData)
		require.NoError(t, err)

		// Снова получаем контенты с целевым статусом
		updatedContents, err := testStorage.GetVideoContentsByDeliveryStatus(ctx, targetStatus, limit)
		require.NoError(t, err)

		// Проверяем что обновленный контент пропал из списка
		foundUpdatedContent := false
		for _, item := range updatedContents {
			if item.ID == contentToUpdate.ID {
				foundUpdatedContent = true
				break
			}
		}
		require.False(t, foundUpdatedContent, "updated content should not be in the list")

		// Проверяем что остальные целевые контенты остались
		remainingTargets := 0
		for _, item := range updatedContents {
			for _, target := range targetContents {
				if item.ID == target.ID && item.ID != contentToUpdate.ID {
					remainingTargets++
					break
				}
			}
		}
		require.Equal(t, len(targetContents)-1, remainingTargets)
	})
}
