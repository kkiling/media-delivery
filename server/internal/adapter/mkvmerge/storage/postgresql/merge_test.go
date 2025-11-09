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

	"github.com/kkiling/media-delivery/internal/adapter/mkvmerge"
)

func randID() uint64 {
	return uint64(rand.Uint32())
}

func TestStorage_Create(t *testing.T) {
	t.Parallel()

	testStorage := NewTestStorage(testutils.SetupPostgresqlTestDB(t))
	ctx := context.Background()

	t.Run("create mkv merge - successfully", func(t *testing.T) {
		t.Parallel()

		create := &mkvmerge.CreateMergeResult{
			ID:             uuid.New(),
			IdempotencyKey: uuid.NewString(),
			Params: mkvmerge.MergeParams{
				VideoInputFile:  "/input/video.mkv",
				VideoOutputFile: "/output/video.mkv",
				AudioTracks: []mkvmerge.Track{
					{
						Path:     "/audio/eng.aac",
						Language: lo.ToPtr("eng"),
						Name:     "English",
						Default:  true,
					},
					{
						Path:     "/audio/rus.aac",
						Language: lo.ToPtr("rus"),
						Name:     "Russian",
						Default:  false,
					},
				},
				SubtitleTracks: []mkvmerge.Track{
					{
						Path:     "/subs/eng.srt",
						Language: lo.ToPtr("eng"),
						Name:     "English",
						Default:  true,
					},
				},
				KeepOriginalAudio:     true,
				KeepOriginalSubtitles: false,
			},
			Status:    mkvmerge.PendingStatus,
			CreatedAt: time.Now(),
		}

		err := testStorage.Create(ctx, create)
		require.NoError(t, err)

		// Проверяем через GetByID
		result, err := testStorage.GetByID(ctx, create.ID)
		require.NoError(t, err)

		require.Equal(t, create.ID, result.ID)
		require.Equal(t, create.IdempotencyKey, result.IdempotencyKey)
		require.Equal(t, create.Params, result.Params)
		require.Equal(t, create.Status, result.Status)
		require.WithinDuration(t, create.CreatedAt, result.CreatedAt, time.Second)
		require.Nil(t, result.Error)
		require.Nil(t, result.CompletedAt)
		require.Nil(t, result.Progress)
	})

	t.Run("create mkv merge - already exists", func(t *testing.T) {
		t.Parallel()

		create := &mkvmerge.CreateMergeResult{
			ID:             uuid.New(),
			IdempotencyKey: uuid.NewString(),
			Params: mkvmerge.MergeParams{
				VideoInputFile:  "/input/video.mkv",
				VideoOutputFile: "/output/video.mkv",
				AudioTracks:     []mkvmerge.Track{},
				SubtitleTracks:  []mkvmerge.Track{},
			},
			Status:    mkvmerge.PendingStatus,
			CreatedAt: time.Now(),
		}

		// Первое создание
		err := testStorage.Create(ctx, create)
		require.NoError(t, err)

		// Второе создание с тем же ID
		err = testStorage.Create(ctx, create)
		require.Error(t, err)
		require.ErrorIs(t, storagebase.ErrAlreadyExists, err)
	})

	t.Run("create mkv merge - with empty params", func(t *testing.T) {
		t.Parallel()

		create := &mkvmerge.CreateMergeResult{
			ID:             uuid.New(),
			IdempotencyKey: uuid.NewString(),
			Params: mkvmerge.MergeParams{
				VideoInputFile:        "/input/video.mkv",
				VideoOutputFile:       "/output/video.mkv",
				AudioTracks:           []mkvmerge.Track{},
				SubtitleTracks:        []mkvmerge.Track{},
				KeepOriginalAudio:     false,
				KeepOriginalSubtitles: false,
			},
			Status:    mkvmerge.RunningStatus,
			CreatedAt: time.Now(),
		}

		err := testStorage.Create(ctx, create)
		require.NoError(t, err)

		// Проверяем через GetByID
		result, err := testStorage.GetByID(ctx, create.ID)
		require.NoError(t, err)

		require.Equal(t, create.ID, result.ID)
		require.Equal(t, create.Params, result.Params)
		require.Equal(t, create.Status, result.Status)
		require.Empty(t, result.Params.AudioTracks)
		require.Empty(t, result.Params.SubtitleTracks)
	})

	t.Run("create mkv merge - with idempotency key conflict", func(t *testing.T) {
		t.Parallel()

		idempotencyKey := uuid.NewString()

		create1 := &mkvmerge.CreateMergeResult{
			ID:             uuid.New(),
			IdempotencyKey: idempotencyKey,
			Params: mkvmerge.MergeParams{
				VideoInputFile:  "/input/video1.mkv",
				VideoOutputFile: "/output/video1.mkv",
				AudioTracks:     []mkvmerge.Track{},
				SubtitleTracks:  []mkvmerge.Track{},
			},
			Status:    mkvmerge.PendingStatus,
			CreatedAt: time.Now(),
		}

		create2 := &mkvmerge.CreateMergeResult{
			ID:             uuid.New(),
			IdempotencyKey: idempotencyKey, // Тот же idempotency key
			Params: mkvmerge.MergeParams{
				VideoInputFile:  "/input/video2.mkv",
				VideoOutputFile: "/output/video2.mkv",
				AudioTracks:     []mkvmerge.Track{},
				SubtitleTracks:  []mkvmerge.Track{},
			},
			Status:    mkvmerge.PendingStatus,
			CreatedAt: time.Now(),
		}

		// Первое создание
		err := testStorage.Create(ctx, create1)
		require.NoError(t, err)

		// Второе создание с тем же idempotency key
		err = testStorage.Create(ctx, create2)
		require.Error(t, err)
		require.ErrorIs(t, storagebase.ErrAlreadyExists, err)
	})
}

func TestStorage_Update(t *testing.T) {
	t.Parallel()

	testStorage := NewTestStorage(testutils.SetupPostgresqlTestDB(t))
	ctx := context.Background()

	t.Run("update mkv merge - successfully update status and completed time", func(t *testing.T) {
		t.Parallel()

		// Сначала создаем запись
		create := &mkvmerge.CreateMergeResult{
			ID:             uuid.New(),
			IdempotencyKey: uuid.NewString(),
			Params: mkvmerge.MergeParams{
				VideoInputFile:  "/input/video.mkv",
				VideoOutputFile: "/output/video.mkv",
				AudioTracks:     []mkvmerge.Track{},
				SubtitleTracks:  []mkvmerge.Track{},
			},
			Status:    mkvmerge.PendingStatus,
			CreatedAt: time.Now(),
		}

		err := testStorage.Create(ctx, create)
		require.NoError(t, err)

		// Обновляем запись
		completedTime := time.Now().Add(time.Hour)
		update := &mkvmerge.UpdateMergeResult{
			Status:    mkvmerge.CompleteStatus,
			Error:     nil,
			Completed: &completedTime,
		}

		err = testStorage.Update(ctx, create.ID, update)
		require.NoError(t, err)

		// Проверяем через GetByID
		result, err := testStorage.GetByID(ctx, create.ID)
		require.NoError(t, err)

		require.Equal(t, mkvmerge.CompleteStatus, result.Status)
		require.Nil(t, result.Error)
		require.NotNil(t, result.CompletedAt)
		require.WithinDuration(t, completedTime, *result.CompletedAt, time.Second)
		// Проверяем что остальные поля не изменились
		require.Equal(t, create.ID, result.ID)
		require.Equal(t, create.IdempotencyKey, result.IdempotencyKey)
		require.Equal(t, create.Params, result.Params)
		require.WithinDuration(t, create.CreatedAt, result.CreatedAt, time.Second)
	})

	t.Run("update mkv merge - successfully update with error", func(t *testing.T) {
		t.Parallel()

		// Сначала создаем запись
		create := &mkvmerge.CreateMergeResult{
			ID:             uuid.New(),
			IdempotencyKey: uuid.NewString(),
			Params: mkvmerge.MergeParams{
				VideoInputFile:  "/input/video.mkv",
				VideoOutputFile: "/output/video.mkv",
				AudioTracks:     []mkvmerge.Track{},
				SubtitleTracks:  []mkvmerge.Track{},
			},
			Status:    mkvmerge.RunningStatus,
			CreatedAt: time.Now(),
		}

		err := testStorage.Create(ctx, create)
		require.NoError(t, err)

		// Обновляем запись с ошибкой
		errorMsg := "merge failed: invalid input file"
		update := &mkvmerge.UpdateMergeResult{
			Status:    mkvmerge.ErrorStatus,
			Error:     &errorMsg,
			Completed: nil,
		}

		err = testStorage.Update(ctx, create.ID, update)
		require.NoError(t, err)

		// Проверяем через GetByID
		result, err := testStorage.GetByID(ctx, create.ID)
		require.NoError(t, err)

		require.Equal(t, mkvmerge.ErrorStatus, result.Status)
		require.NotNil(t, result.Error)
		require.Equal(t, errorMsg, *result.Error)
		require.Nil(t, result.CompletedAt)
	})

	t.Run("update mkv merge - not found", func(t *testing.T) {
		t.Parallel()

		nonExistentID := uuid.New()
		update := &mkvmerge.UpdateMergeResult{
			Status: mkvmerge.CompleteStatus,
			Error:  nil,
		}

		err := testStorage.Update(ctx, nonExistentID, update)
		require.Error(t, err)
		require.ErrorIs(t, storagebase.ErrNotFound, err)
	})

	t.Run("update mkv merge - partial update", func(t *testing.T) {
		t.Parallel()

		// Сначала создаем запись
		create := &mkvmerge.CreateMergeResult{
			ID:             uuid.New(),
			IdempotencyKey: uuid.NewString(),
			Params: mkvmerge.MergeParams{
				VideoInputFile:  "/input/video.mkv",
				VideoOutputFile: "/output/video.mkv",
				AudioTracks:     []mkvmerge.Track{},
				SubtitleTracks:  []mkvmerge.Track{},
			},
			Status:    mkvmerge.PendingStatus,
			CreatedAt: time.Now(),
		}

		err := testStorage.Create(ctx, create)
		require.NoError(t, err)

		// Обновляем только статус
		update := &mkvmerge.UpdateMergeResult{
			Status:    mkvmerge.RunningStatus,
			Error:     nil,
			Completed: nil,
		}

		err = testStorage.Update(ctx, create.ID, update)
		require.NoError(t, err)

		// Проверяем через GetByID
		result, err := testStorage.GetByID(ctx, create.ID)
		require.NoError(t, err)

		require.Equal(t, mkvmerge.RunningStatus, result.Status)
		require.Nil(t, result.Error)
		require.Nil(t, result.CompletedAt)
		// Проверяем что остальные поля не изменились
		require.Equal(t, create.ID, result.ID)
		require.Equal(t, create.IdempotencyKey, result.IdempotencyKey)
		require.Equal(t, create.Params, result.Params)
	})

	t.Run("update mkv merge - multiple updates", func(t *testing.T) {
		t.Parallel()

		// Сначала создаем запись
		create := &mkvmerge.CreateMergeResult{
			ID:             uuid.New(),
			IdempotencyKey: uuid.NewString(),
			Params: mkvmerge.MergeParams{
				VideoInputFile:  "/input/video.mkv",
				VideoOutputFile: "/output/video.mkv",
				AudioTracks:     []mkvmerge.Track{},
				SubtitleTracks:  []mkvmerge.Track{},
			},
			Status:    mkvmerge.PendingStatus,
			CreatedAt: time.Now(),
		}

		err := testStorage.Create(ctx, create)
		require.NoError(t, err)

		// Первое обновление - запускаем процесс
		update1 := &mkvmerge.UpdateMergeResult{
			Status: mkvmerge.RunningStatus,
		}
		err = testStorage.Update(ctx, create.ID, update1)
		require.NoError(t, err)

		result1, err := testStorage.GetByID(ctx, create.ID)
		require.NoError(t, err)
		require.Equal(t, mkvmerge.RunningStatus, result1.Status)

		// Второе обновление - завершаем успешно
		completedTime := time.Now()
		update2 := &mkvmerge.UpdateMergeResult{
			Status:    mkvmerge.CompleteStatus,
			Completed: &completedTime,
		}
		err = testStorage.Update(ctx, create.ID, update2)
		require.NoError(t, err)

		result2, err := testStorage.GetByID(ctx, create.ID)
		require.NoError(t, err)
		require.Equal(t, mkvmerge.CompleteStatus, result2.Status)
		require.NotNil(t, result2.CompletedAt)
		require.WithinDuration(t, completedTime, *result2.CompletedAt, time.Second)
	})
}

func TestStorage_UpdateProgress(t *testing.T) {
	t.Parallel()

	testStorage := NewTestStorage(testutils.SetupPostgresqlTestDB(t))
	ctx := context.Background()

	t.Run("update progress - successfully", func(t *testing.T) {
		t.Parallel()

		// Сначала создаем запись
		create := &mkvmerge.CreateMergeResult{
			ID:             uuid.New(),
			IdempotencyKey: uuid.NewString(),
			Params: mkvmerge.MergeParams{
				VideoInputFile:  "/input/video.mkv",
				VideoOutputFile: "/output/video.mkv",
				AudioTracks:     []mkvmerge.Track{},
				SubtitleTracks:  []mkvmerge.Track{},
			},
			Status:    mkvmerge.RunningStatus,
			CreatedAt: time.Now(),
		}

		err := testStorage.Create(ctx, create)
		require.NoError(t, err)

		// Обновляем прогресс
		progress := float32(50.5)
		err = testStorage.UpdateProgress(ctx, create.ID, progress)
		require.NoError(t, err)

		// Проверяем через GetByID
		result, err := testStorage.GetByID(ctx, create.ID)
		require.NoError(t, err)

		require.NotNil(t, result.Progress)
		require.Equal(t, progress, *result.Progress)
		// Проверяем что остальные поля не изменились
		require.Equal(t, create.ID, result.ID)
		require.Equal(t, create.Status, result.Status)
		require.Equal(t, create.Params, result.Params)
	})

	t.Run("update progress - not found", func(t *testing.T) {
		t.Parallel()

		nonExistentID := uuid.New()
		progress := float32(25.0)

		err := testStorage.UpdateProgress(ctx, nonExistentID, progress)
		require.Error(t, err)
		require.ErrorIs(t, storagebase.ErrNotFound, err)
	})

	t.Run("update progress - multiple updates", func(t *testing.T) {
		t.Parallel()

		// Сначала создаем запись
		create := &mkvmerge.CreateMergeResult{
			ID:             uuid.New(),
			IdempotencyKey: uuid.NewString(),
			Params: mkvmerge.MergeParams{
				VideoInputFile:  "/input/video.mkv",
				VideoOutputFile: "/output/video.mkv",
				AudioTracks:     []mkvmerge.Track{},
				SubtitleTracks:  []mkvmerge.Track{},
			},
			Status:    mkvmerge.RunningStatus,
			CreatedAt: time.Now(),
		}

		err := testStorage.Create(ctx, create)
		require.NoError(t, err)

		// Последовательные обновления прогресса
		progressUpdates := []float32{10.0, 25.5, 50.0, 75.3, 100.0}

		for _, progress := range progressUpdates {
			err = testStorage.UpdateProgress(ctx, create.ID, progress)
			require.NoError(t, err)

			// Проверяем что прогресс обновился
			result, err := testStorage.GetByID(ctx, create.ID)
			require.NoError(t, err)

			require.NotNil(t, result.Progress)
			require.Equal(t, progress, *result.Progress)
		}

		// Финальная проверка
		result, err := testStorage.GetByID(ctx, create.ID)
		require.NoError(t, err)
		require.Equal(t, float32(100.0), *result.Progress)
	})

	t.Run("update progress - zero progress", func(t *testing.T) {
		t.Parallel()

		// Сначала создаем запись
		create := &mkvmerge.CreateMergeResult{
			ID:             uuid.New(),
			IdempotencyKey: uuid.NewString(),
			Params: mkvmerge.MergeParams{
				VideoInputFile:  "/input/video.mkv",
				VideoOutputFile: "/output/video.mkv",
				AudioTracks:     []mkvmerge.Track{},
				SubtitleTracks:  []mkvmerge.Track{},
			},
			Status:    mkvmerge.RunningStatus,
			CreatedAt: time.Now(),
		}

		err := testStorage.Create(ctx, create)
		require.NoError(t, err)

		// Обновляем прогресс на 0
		err = testStorage.UpdateProgress(ctx, create.ID, 0.0)
		require.NoError(t, err)

		// Проверяем через GetByID
		result, err := testStorage.GetByID(ctx, create.ID)
		require.NoError(t, err)

		require.NotNil(t, result.Progress)
		require.Equal(t, float32(0.0), *result.Progress)
	})

	t.Run("update progress - negative progress", func(t *testing.T) {
		t.Parallel()

		// Сначала создаем запись
		create := &mkvmerge.CreateMergeResult{
			ID:             uuid.New(),
			IdempotencyKey: uuid.NewString(),
			Params: mkvmerge.MergeParams{
				VideoInputFile:  "/input/video.mkv",
				VideoOutputFile: "/output/video.mkv",
				AudioTracks:     []mkvmerge.Track{},
				SubtitleTracks:  []mkvmerge.Track{},
			},
			Status:    mkvmerge.RunningStatus,
			CreatedAt: time.Now(),
		}

		err := testStorage.Create(ctx, create)
		require.NoError(t, err)

		// Обновляем прогресс на отрицательное значение
		err = testStorage.UpdateProgress(ctx, create.ID, -10.0)
		require.NoError(t, err)

		// Проверяем через GetByID
		result, err := testStorage.GetByID(ctx, create.ID)
		require.NoError(t, err)

		require.NotNil(t, result.Progress)
		require.Equal(t, float32(-10.0), *result.Progress)
	})

	t.Run("update progress - combined with status update", func(t *testing.T) {
		t.Parallel()

		// Сначала создаем запись
		create := &mkvmerge.CreateMergeResult{
			ID:             uuid.New(),
			IdempotencyKey: uuid.NewString(),
			Params: mkvmerge.MergeParams{
				VideoInputFile:  "/input/video.mkv",
				VideoOutputFile: "/output/video.mkv",
				AudioTracks:     []mkvmerge.Track{},
				SubtitleTracks:  []mkvmerge.Track{},
			},
			Status:    mkvmerge.PendingStatus,
			CreatedAt: time.Now(),
		}

		err := testStorage.Create(ctx, create)
		require.NoError(t, err)

		// Обновляем статус
		statusUpdate := &mkvmerge.UpdateMergeResult{
			Status: mkvmerge.RunningStatus,
		}
		err = testStorage.Update(ctx, create.ID, statusUpdate)
		require.NoError(t, err)

		// Обновляем прогресс
		progress := float32(33.3)
		err = testStorage.UpdateProgress(ctx, create.ID, progress)
		require.NoError(t, err)

		// Проверяем результат
		result, err := testStorage.GetByID(ctx, create.ID)
		require.NoError(t, err)

		require.Equal(t, mkvmerge.RunningStatus, result.Status)
		require.NotNil(t, result.Progress)
		require.Equal(t, progress, *result.Progress)
		// Проверяем что остальные поля не изменились
		require.Equal(t, create.ID, result.ID)
		require.Equal(t, create.Params, result.Params)
	})
}

func TestStorage_GetByID(t *testing.T) {
	t.Parallel()

	testStorage := NewTestStorage(testutils.SetupPostgresqlTestDB(t))
	ctx := context.Background()

	t.Run("get by id - successfully", func(t *testing.T) {
		t.Parallel()

		// Сначала создаем запись
		create := &mkvmerge.CreateMergeResult{
			ID:             uuid.New(),
			IdempotencyKey: uuid.NewString(),
			Params: mkvmerge.MergeParams{
				VideoInputFile:  "/input/video.mkv",
				VideoOutputFile: "/output/video.mkv",
				AudioTracks: []mkvmerge.Track{
					{
						Path:     "/audio/eng.aac",
						Language: lo.ToPtr("eng"),
						Name:     "English",
						Default:  true,
					},
				},
				SubtitleTracks:        []mkvmerge.Track{},
				KeepOriginalAudio:     true,
				KeepOriginalSubtitles: false,
			},
			Status:    mkvmerge.RunningStatus,
			CreatedAt: time.Now(),
		}

		err := testStorage.Create(ctx, create)
		require.NoError(t, err)

		// Получаем запись по ID
		result, err := testStorage.GetByID(ctx, create.ID)
		require.NoError(t, err)

		require.Equal(t, create.ID, result.ID)
		require.Equal(t, create.IdempotencyKey, result.IdempotencyKey)
		require.Equal(t, create.Params, result.Params)
		require.Equal(t, create.Status, result.Status)
		require.WithinDuration(t, create.CreatedAt, result.CreatedAt, time.Second)
		require.Nil(t, result.Error)
		require.Nil(t, result.CompletedAt)
		require.Nil(t, result.Progress)
	})

	t.Run("get by id - not found", func(t *testing.T) {
		t.Parallel()

		nonExistentID := uuid.New()
		result, err := testStorage.GetByID(ctx, nonExistentID)
		require.Error(t, err)
		require.Nil(t, result)
		require.ErrorIs(t, storagebase.ErrNotFound, err)
	})
}

func TestStorage_GetByIdempotencyKey(t *testing.T) {
	t.Parallel()

	testStorage := NewTestStorage(testutils.SetupPostgresqlTestDB(t))
	ctx := context.Background()

	t.Run("get by idempotency key - successfully", func(t *testing.T) {
		t.Parallel()

		idempotencyKey := uuid.NewString()

		// Сначала создаем запись
		create := &mkvmerge.CreateMergeResult{
			ID:             uuid.New(),
			IdempotencyKey: idempotencyKey,
			Params: mkvmerge.MergeParams{
				VideoInputFile:  "/input/video.mkv",
				VideoOutputFile: "/output/video.mkv",
				AudioTracks: []mkvmerge.Track{
					{
						Path:     "/audio/rus.aac",
						Language: lo.ToPtr("rus"),
						Name:     "Russian",
						Default:  false,
					},
				},
				SubtitleTracks: []mkvmerge.Track{
					{
						Path:     "/subs/eng.srt",
						Language: lo.ToPtr("eng"),
						Name:     "English",
						Default:  true,
					},
				},
				KeepOriginalAudio:     false,
				KeepOriginalSubtitles: true,
			},
			Status:    mkvmerge.CompleteStatus,
			CreatedAt: time.Now(),
		}

		err := testStorage.Create(ctx, create)
		require.NoError(t, err)

		// Получаем запись по idempotency key
		result, err := testStorage.GetByIdempotencyKey(ctx, idempotencyKey)
		require.NoError(t, err)

		require.Equal(t, create.ID, result.ID)
		require.Equal(t, create.IdempotencyKey, result.IdempotencyKey)
		require.Equal(t, create.Params, result.Params)
		require.Equal(t, create.Status, result.Status)
		require.WithinDuration(t, create.CreatedAt, result.CreatedAt, time.Second)
		require.Nil(t, result.Error)
		require.Nil(t, result.CompletedAt)
		require.Nil(t, result.Progress)
	})

	t.Run("get by idempotency key - not found", func(t *testing.T) {
		t.Parallel()

		nonExistentKey := uuid.NewString()
		result, err := testStorage.GetByIdempotencyKey(ctx, nonExistentKey)
		require.Error(t, err)
		require.Nil(t, result)
		require.ErrorIs(t, storagebase.ErrNotFound, err)
	})
}
