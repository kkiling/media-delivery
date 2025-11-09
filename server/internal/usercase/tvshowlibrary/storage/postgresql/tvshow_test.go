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

	"github.com/kkiling/media-delivery/internal/usercase/tvshowlibrary"
)

func randID() uint64 {
	return uint64(rand.Uint32())
}

func equalImg(t *testing.T, img1, img2 *tvshowlibrary.Image) {
	require.True(t, (img1 == nil && img2 == nil) || (img1 != nil && img2 != nil))
	if img1 == nil || img2 == nil {
		return
	}
	require.Equal(t, img1.ID, img2.ID)
	require.Equal(t, img1.W92, img2.W92)
	require.Equal(t, img1.W185, img2.W185)
	require.Equal(t, img1.W342, img2.W342)
	require.Equal(t, img1.Original, img2.Original)
}

func equalSeason(t *testing.T, s1, s2 *tvshowlibrary.Season) {
	require.True(t, (s1 == nil && s2 == nil) || (s1 != nil && s2 != nil))
	if s1 == nil || s2 == nil {
		return
	}
	require.WithinDuration(t, s1.AirDate, s2.AirDate, time.Second)
	require.Equal(t, s1.EpisodeCount, s2.EpisodeCount)
	require.Equal(t, s1.Name, s2.Name)
	require.Equal(t, s1.Overview, s2.Overview)
	equalImg(t, s1.Poster, s2.Poster)
	require.Equal(t, s1.SeasonNumber, s2.SeasonNumber)
	require.Equal(t, s1.VoteAverage, s2.VoteAverage)
}

func equalTVShow(t *testing.T, tv1, tv2 *tvshowlibrary.TVShow) {
	require.True(t, (tv1 == nil && tv2 == nil) || (tv1 != nil && tv2 != nil))
	if tv1 == nil || tv2 == nil {
		return
	}

	// Сравниваем TVShowShort часть
	require.Equal(t, tv1.ID, tv2.ID)
	require.Equal(t, tv1.Name, tv2.Name)
	require.Equal(t, tv1.OriginalName, tv2.OriginalName)
	require.Equal(t, tv1.Overview, tv2.Overview)
	equalImg(t, tv1.Poster, tv2.Poster)
	require.WithinDuration(t, tv1.FirstAirDate, tv2.FirstAirDate, time.Second)
	require.Equal(t, tv1.VoteAverage, tv2.VoteAverage)
	require.Equal(t, tv1.VoteCount, tv2.VoteCount)
	require.Equal(t, tv1.Popularity, tv2.Popularity)

	// Сравниваем остальные поля TVShow
	equalImg(t, tv1.Backdrop, tv2.Backdrop)
	require.Equal(t, tv1.Genres, tv2.Genres)
	require.WithinDuration(t, tv1.LastAirDate, tv2.LastAirDate, time.Second)
	require.Equal(t, tv1.NumberOfEpisodes, tv2.NumberOfEpisodes)
	require.Equal(t, tv1.NumberOfSeasons, tv2.NumberOfSeasons)
	require.Equal(t, tv1.OriginCountry, tv2.OriginCountry)
	require.Equal(t, tv1.Status, tv2.Status)
	require.Equal(t, tv1.Tagline, tv2.Tagline)
	require.Equal(t, tv1.Type, tv2.Type)

	// Сравниваем сезоны
	require.Equal(t, len(tv1.Seasons), len(tv2.Seasons))
	for i := range tv1.Seasons {
		equalSeason(t, &tv1.Seasons[i], &tv2.Seasons[i])
	}
}

func equalEpisode(t *testing.T, e1, e2 *tvshowlibrary.Episode) {
	require.True(t, (e1 == nil && e2 == nil) || (e1 != nil && e2 != nil))
	if e1 == nil || e2 == nil {
		return
	}

	require.WithinDuration(t, e1.AirDate, e2.AirDate, time.Second)
	require.Equal(t, e1.EpisodeNumber, e2.EpisodeNumber)
	require.Equal(t, e1.EpisodeType, e2.EpisodeType)
	require.Equal(t, e1.Name, e2.Name)
	require.Equal(t, e1.Overview, e2.Overview)
	require.Equal(t, e1.Runtime, e2.Runtime)
	equalImg(t, e1.Still, e2.Still)
	require.Equal(t, e1.VoteAverage, e2.VoteAverage)
	require.Equal(t, e1.VoteCount, e2.VoteCount)
}

func equalEpisodes(t *testing.T, episodes1, episodes2 []tvshowlibrary.Episode) {
	require.Equal(t, len(episodes1), len(episodes2))
	for i := range episodes1 {
		equalEpisode(t, &episodes1[i], &episodes2[i])
	}
}

func TestStorage_SaveImage(t *testing.T) {
	t.Parallel()

	testStorage := NewTestStorage(testutils.SetupPostgresqlTestDB(t))
	ctx := context.Background()

	t.Run("save image - successfully", func(t *testing.T) {
		t.Parallel()

		id := uuid.NewString()
		img := &tvshowlibrary.Image{
			ID:       id,
			W92:      lo.ToPtr("img/w92"),
			W185:     lo.ToPtr("img/w92"),
			W342:     lo.ToPtr("img/w92"),
			Original: "img/original",
		}
		err := testStorage.saveImage(ctx, img)
		require.NoError(t, err)

		findImg, err := testStorage.getImage(ctx, id)
		require.NoError(t, err)
		equalImg(t, img, findImg)
	})

	t.Run("save image - already exist", func(t *testing.T) {
		t.Parallel()

		id := uuid.NewString()

		img := &tvshowlibrary.Image{
			ID:       id,
			W92:      lo.ToPtr("img/w92"),
			W185:     lo.ToPtr("img/w92"),
			W342:     lo.ToPtr("img/w92"),
			Original: "img/original",
		}
		err1 := testStorage.saveImage(ctx, img)
		require.NoError(t, err1)

		err2 := testStorage.saveImage(ctx, img)
		require.Error(t, err2)
		require.ErrorIs(t, storagebase.ErrAlreadyExists, err2)
	})
}

func TestStorage_SaveSeason(t *testing.T) {
	t.Parallel()

	testStorage := NewTestStorage(testutils.SetupPostgresqlTestDB(t))
	ctx := context.Background()

	t.Run("save season - successfully with poster", func(t *testing.T) {
		t.Parallel()

		tvShowID := randID()
		posterID := uuid.NewString()
		season := tvshowlibrary.Season{
			SeasonNumber: 1,
			AirDate:      time.Now(),
			EpisodeCount: 10,
			Name:         "Season 1",
			Overview:     "Season 1 overview",
			Poster: &tvshowlibrary.Image{
				ID:       posterID,
				W92:      lo.ToPtr("season/w92"),
				W185:     lo.ToPtr("season/w185"),
				W342:     lo.ToPtr("season/w342"),
				Original: "season/original",
			},
			VoteAverage: 8.0,
		}

		err := testStorage.saveSeason(ctx, tvShowID, season)
		require.NoError(t, err)

		// Проверяем что сезон сохранился
		seasons, err := testStorage.getSeasons(ctx, tvShowID)
		require.NoError(t, err)
		require.Len(t, seasons, 1)

		// Проверяем что постер сохранился
		require.NotNil(t, seasons[0].Poster)
		equalSeason(t, &season, &seasons[0])
	})

	t.Run("save season - successfully without poster", func(t *testing.T) {
		t.Parallel()

		tvShowID := randID()
		season := tvshowlibrary.Season{
			SeasonNumber: 2,
			AirDate:      time.Now(),
			EpisodeCount: 8,
			Name:         "Season 2",
			Overview:     "Season 2 overview",
			Poster:       nil,
			VoteAverage:  7.5,
		}

		err := testStorage.saveSeason(ctx, tvShowID, season)
		require.NoError(t, err)

		// Проверяем что сезон сохранился
		seasons, err := testStorage.getSeasons(ctx, tvShowID)
		require.NoError(t, err)
		require.Len(t, seasons, 1)
		require.Nil(t, seasons[0].Poster)
		equalSeason(t, &season, &seasons[0])
	})

	t.Run("save season - already exists", func(t *testing.T) {
		t.Parallel()

		tvShowID := randID()
		season := tvshowlibrary.Season{
			SeasonNumber: 1,
			AirDate:      time.Now(),
			EpisodeCount: 12,
			Name:         "Season 1",
			Overview:     "Season 1 overview",
			VoteAverage:  8.5,
		}

		// Первое сохранение должно быть успешным
		err := testStorage.saveSeason(ctx, tvShowID, season)
		require.NoError(t, err)

		// Второе сохранение того же сезона должно вернуть ошибку
		err = testStorage.saveSeason(ctx, tvShowID, season)
		require.Error(t, err)
		require.ErrorIs(t, storagebase.ErrAlreadyExists, err)
	})

	t.Run("save multiple seasons for same tv show", func(t *testing.T) {
		t.Parallel()

		tvShowID := randID()
		seasons := []tvshowlibrary.Season{
			{
				SeasonNumber: 1,
				AirDate:      time.Now(),
				EpisodeCount: 10,
				Name:         "Season 1",
				Overview:     "First season",
				VoteAverage:  8.0,
			},
			{
				SeasonNumber: 2,
				AirDate:      time.Now().AddDate(0, 6, 0),
				EpisodeCount: 12,
				Name:         "Season 2",
				Overview:     "Second season",
				VoteAverage:  8.5,
			},
			{
				SeasonNumber: 3,
				AirDate:      time.Now().AddDate(1, 0, 0),
				EpisodeCount: 8,
				Name:         "Season 3",
				Overview:     "Third season",
				VoteAverage:  9.0,
			},
		}

		// Сохраняем все сезоны
		for _, season := range seasons {
			err := testStorage.saveSeason(ctx, tvShowID, season)
			require.NoError(t, err)
		}

		// Проверяем что все сезоны сохранились
		savedSeasons, err := testStorage.getSeasons(ctx, tvShowID)
		require.NoError(t, err)
		require.Len(t, savedSeasons, 3)

		// Проверяем что сезоны отсортированы по номеру
		require.Equal(t, uint8(1), savedSeasons[0].SeasonNumber)
		require.Equal(t, uint8(2), savedSeasons[1].SeasonNumber)
		require.Equal(t, uint8(3), savedSeasons[2].SeasonNumber)
	})
}

func TestStorage_SaveTVShow(t *testing.T) {
	t.Parallel()

	testStorage := NewTestStorage(testutils.SetupPostgresqlTestDB(t))
	ctx := context.Background()

	t.Run("save tv show - successfully with images and seasons", func(t *testing.T) {
		t.Parallel()

		posterID := uuid.NewString()
		backdropID := uuid.NewString()
		seasonPosterID := uuid.NewString()

		tvShow := &tvshowlibrary.TVShow{
			TVShowShort: tvshowlibrary.TVShowShort{
				ID:           randID(),
				Name:         "Test TV Show",
				OriginalName: "Test TV Show Original",
				Overview:     "Test overview",
				Poster: &tvshowlibrary.Image{
					ID:       posterID,
					W92:      lo.ToPtr("poster/w92"),
					W185:     lo.ToPtr("poster/w185"),
					W342:     lo.ToPtr("poster/w342"),
					Original: "poster/original",
				},
				FirstAirDate: time.Now(),
				VoteAverage:  8.5,
				VoteCount:    100,
				Popularity:   50.0,
			},
			Backdrop: &tvshowlibrary.Image{
				ID:       backdropID,
				W92:      lo.ToPtr("backdrop/w92"),
				W185:     lo.ToPtr("backdrop/w185"),
				W342:     lo.ToPtr("backdrop/w342"),
				Original: "backdrop/original",
			},
			Genres:           []string{"Drama", "Action"},
			LastAirDate:      time.Now().AddDate(1, 0, 0),
			NumberOfEpisodes: 10,
			NumberOfSeasons:  1,
			OriginCountry:    []string{"US", "UK"},
			Status:           lo.ToPtr("Running"),
			Tagline:          "Great TV show",
			Type:             lo.ToPtr("Scripted"),
			Seasons: []tvshowlibrary.Season{
				{
					SeasonNumber: 1,
					AirDate:      time.Now(),
					EpisodeCount: 10,
					Name:         "Season 1",
					Overview:     "First season",
					Poster: &tvshowlibrary.Image{
						ID:       seasonPosterID,
						W92:      lo.ToPtr("season/w92"),
						W185:     lo.ToPtr("season/w185"),
						W342:     lo.ToPtr("season/w342"),
						Original: "season/original",
					},
					VoteAverage: 8.0,
				},
			},
		}

		err := testStorage.SaveTVShow(ctx, tvShow)
		require.NoError(t, err)

		// Проверяем через GetTVShow
		savedTVShow, err := testStorage.GetTVShow(ctx, tvShow.ID)
		require.NoError(t, err)

		equalTVShow(t, tvShow, savedTVShow)
	})

	t.Run("save tv show - successfully without optional fields", func(t *testing.T) {
		t.Parallel()

		tvShow := &tvshowlibrary.TVShow{
			TVShowShort: tvshowlibrary.TVShowShort{
				ID:           randID(),
				Name:         "Minimal TV Show",
				OriginalName: "Minimal TV Show Original",
				Overview:     "Minimal overview",
				Poster:       nil,
				FirstAirDate: time.Now(),
				VoteAverage:  7.0,
				VoteCount:    50,
				Popularity:   25.0,
			},
			Backdrop:         nil,
			Genres:           []string{},
			LastAirDate:      time.Time{},
			NumberOfEpisodes: 0,
			NumberOfSeasons:  0,
			OriginCountry:    []string{},
			Status:           nil,
			Tagline:          "",
			Type:             nil,
			Seasons:          []tvshowlibrary.Season{},
		}

		err := testStorage.SaveTVShow(ctx, tvShow)
		require.NoError(t, err)

		// Проверяем через GetTVShow
		savedTVShow, err := testStorage.GetTVShow(ctx, tvShow.ID)
		require.NoError(t, err)

		equalTVShow(t, tvShow, savedTVShow)
	})

	t.Run("save tv show - already exists", func(t *testing.T) {
		t.Parallel()

		tvShow := &tvshowlibrary.TVShow{
			TVShowShort: tvshowlibrary.TVShowShort{
				ID:           randID(),
				Name:         "Duplicate TV Show",
				OriginalName: "Duplicate TV Show Original",
				Overview:     "Duplicate overview",
				FirstAirDate: time.Now(),
				VoteAverage:  6.5,
				VoteCount:    30,
				Popularity:   15.0,
			},
		}

		// Первое сохранение
		err := testStorage.SaveTVShow(ctx, tvShow)
		require.NoError(t, err)

		// Второе сохранение того же TVShow
		err = testStorage.SaveTVShow(ctx, tvShow)
		require.Error(t, err)
		require.ErrorIs(t, storagebase.ErrAlreadyExists, err)
	})
}

func TestStorage_GetTVShows(t *testing.T) {
	t.Parallel()

	testStorage := NewTestStorage(testutils.SetupPostgresqlTestDB(t))
	ctx := context.Background()

	t.Run("save multiple tv shows - check GetTVShows", func(t *testing.T) {
		t.Parallel()

		tvShows := []*tvshowlibrary.TVShow{
			{
				TVShowShort: tvshowlibrary.TVShowShort{
					ID:           randID(),
					Name:         "TV Show A",
					OriginalName: "TV Show A Original",
					Overview:     "Overview A",
					FirstAirDate: time.Now(),
					VoteAverage:  8.0,
					VoteCount:    100,
					Popularity:   40.0,
				},
			},
			{
				TVShowShort: tvshowlibrary.TVShowShort{
					ID:           randID(),
					Name:         "TV Show B",
					OriginalName: "TV Show B Original",
					Overview:     "Overview B",
					FirstAirDate: time.Now(),
					VoteAverage:  7.5,
					VoteCount:    80,
					Popularity:   35.0,
				},
			},
		}

		// Сохраняем все TV шоу
		for _, tvShow := range tvShows {
			err := testStorage.SaveTVShow(ctx, tvShow)
			require.NoError(t, err)
		}

		// Проверяем через GetTVShows
		savedTVShows, err := testStorage.GetTVShows(ctx)
		require.NoError(t, err)

		// Находим наши TV шоу среди всех сохраненных
		foundCount := 0
		for _, saved := range savedTVShows {
			for _, original := range tvShows {
				if saved.ID == original.ID {
					foundCount++
					equalTVShow(t, &tvshowlibrary.TVShow{TVShowShort: saved}, original)
				}
			}
		}
		require.Equal(t, 2, foundCount)
	})
}

func TestStorage_getImage(t *testing.T) {
	t.Parallel()

	testStorage := NewTestStorage(testutils.SetupPostgresqlTestDB(t))
	ctx := context.Background()

	t.Run("get image - not found", func(t *testing.T) {
		t.Parallel()

		nonExistentID := uuid.NewString()
		_, err := testStorage.getImage(ctx, nonExistentID)
		require.Error(t, err)
		require.ErrorIs(t, storagebase.ErrNotFound, err)
	})
}

func TestStorage_GetTVShow(t *testing.T) {
	t.Parallel()

	testStorage := NewTestStorage(testutils.SetupPostgresqlTestDB(t))
	ctx := context.Background()

	t.Run("get tvshow - not found", func(t *testing.T) {
		t.Parallel()

		nonExistentID := uint64(999999)
		_, err := testStorage.GetTVShow(ctx, nonExistentID)
		require.Error(t, err)
		require.ErrorIs(t, storagebase.ErrNotFound, err)
	})
}

func TestStorage_SaveEpisodes(t *testing.T) {
	t.Parallel()

	testStorage := NewTestStorage(testutils.SetupPostgresqlTestDB(t))
	ctx := context.Background()

	t.Run("save episodes - successfully with still images", func(t *testing.T) {
		t.Parallel()

		tvShowID := uint64(1)
		seasonNumber := uint8(1)
		stillID := uuid.NewString()

		episodes := []tvshowlibrary.Episode{
			{
				AirDate:       time.Now(),
				EpisodeNumber: 1,
				EpisodeType:   lo.ToPtr("standard"),
				Name:          "Episode 1",
				Overview:      "First episode",
				Runtime:       45,
				Still: &tvshowlibrary.Image{
					ID:       stillID,
					W92:      lo.ToPtr("still/w92"),
					W185:     lo.ToPtr("still/w185"),
					W342:     lo.ToPtr("still/w342"),
					Original: "still/original",
				},
				VoteAverage: 8.0,
				VoteCount:   50,
			},
			{
				AirDate:       time.Now().AddDate(0, 0, 7),
				EpisodeNumber: 2,
				EpisodeType:   lo.ToPtr("standard"),
				Name:          "Episode 2",
				Overview:      "Second episode",
				Runtime:       42,
				Still:         nil,
				VoteAverage:   7.5,
				VoteCount:     40,
			},
		}

		err := testStorage.SaveEpisodes(ctx, tvShowID, seasonNumber, episodes)
		require.NoError(t, err)

		// Проверяем через GetEpisodes
		savedEpisodes, err := testStorage.GetEpisodes(ctx, tvShowID, seasonNumber)
		require.NoError(t, err)

		// Конвертируем []*Episode в []Episode для сравнения
		expectedEpisodes := make([]tvshowlibrary.Episode, len(episodes))
		for i, ep := range episodes {
			expectedEpisodes[i] = ep
		}

		equalEpisodes(t, expectedEpisodes, savedEpisodes)
	})

	t.Run("save episodes - empty slice", func(t *testing.T) {
		t.Parallel()

		tvShowID := uint64(2)
		seasonNumber := uint8(1)

		err := testStorage.SaveEpisodes(ctx, tvShowID, seasonNumber, []tvshowlibrary.Episode{})
		require.NoError(t, err)

		// Проверяем что ничего не сохранилось
		savedEpisodes, err := testStorage.GetEpisodes(ctx, tvShowID, seasonNumber)
		require.NoError(t, err)
		require.Empty(t, savedEpisodes)
	})

	t.Run("save episodes - already exists", func(t *testing.T) {
		t.Parallel()

		tvShowID := uint64(3)
		seasonNumber := uint8(1)

		episodes := []tvshowlibrary.Episode{
			{
				AirDate:       time.Now(),
				EpisodeNumber: 1,
				EpisodeType:   lo.ToPtr("standard"),
				Name:          "Duplicate Episode",
				Overview:      "Duplicate episode",
				Runtime:       45,
				VoteAverage:   8.0,
				VoteCount:     50,
			},
		}

		// Первое сохранение
		err := testStorage.SaveEpisodes(ctx, tvShowID, seasonNumber, episodes)
		require.NoError(t, err)

		// Второе сохранение тех же эпизодов
		err = testStorage.SaveEpisodes(ctx, tvShowID, seasonNumber, episodes)
		require.Error(t, err)
		require.ErrorIs(t, storagebase.ErrAlreadyExists, err)
	})

	t.Run("save episodes - multiple episodes same season", func(t *testing.T) {
		t.Parallel()

		tvShowID := uint64(4)
		seasonNumber := uint8(2)

		episodes := []tvshowlibrary.Episode{
			{
				EpisodeNumber: 1,
				Name:          "Episode 1",
				Overview:      "First episode",
				Runtime:       45,
				VoteAverage:   8.0,
				VoteCount:     50,
			},
			{
				EpisodeNumber: 2,
				Name:          "Episode 2",
				Overview:      "Second episode",
				Runtime:       42,
				VoteAverage:   7.5,
				VoteCount:     40,
			},
			{
				EpisodeNumber: 3,
				Name:          "Episode 3",
				Overview:      "Third episode",
				Runtime:       47,
				VoteAverage:   8.2,
				VoteCount:     60,
			},
		}

		err := testStorage.SaveEpisodes(ctx, tvShowID, seasonNumber, episodes)
		require.NoError(t, err)

		// Проверяем через GetEpisodes
		savedEpisodes, err := testStorage.GetEpisodes(ctx, tvShowID, seasonNumber)
		require.NoError(t, err)
		require.Len(t, savedEpisodes, 3)

		// Проверяем что эпизоды отсортированы по номеру
		require.Equal(t, 1, savedEpisodes[0].EpisodeNumber)
		require.Equal(t, 2, savedEpisodes[1].EpisodeNumber)
		require.Equal(t, 3, savedEpisodes[2].EpisodeNumber)
	})
}
