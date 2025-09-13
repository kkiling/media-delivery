package delivery

import (
	"testing"

	"github.com/samber/lo"
	"github.com/stretchr/testify/require"
)

// Тестирование функции ValidateContentMatch
func TestValidateContentMatch(t *testing.T) {
	s := &Service{}

	videoTrack := Track{
		Type: TrackTypeVideo,
		Name: lo.ToPtr("video1"),
		File: FileInfo{RelativePath: "video1.mkv", FullPath: "/path/video1.mkv"},
	}
	audioTrack := Track{
		Type: TrackTypeAudio,
		Name: lo.ToPtr("audio1"),
		File: FileInfo{RelativePath: "audio1.aac", FullPath: "/path/audio1.aac"},
	}
	subtitleTrack := Track{
		Type: TrackTypeSubtitle,
		Name: lo.ToPtr("sub1"),
		File: FileInfo{RelativePath: "sub1.srt", FullPath: "/path/sub1.srt"},
	}
	episode := EpisodeInfo{
		SeasonNumber:  1,
		EpisodeNumber: 1,
		FullPath:      "/media/ep1.mkv",
		RelativePath:  "ep1.mkv",
	}

	oldMatches := &ContentMatches{
		Matches: []ContentMatch{
			{
				Episode:     episode,
				Video:       &videoTrack,
				AudioTracks: []Track{audioTrack},
				Subtitles:   []Track{subtitleTrack},
			},
		},
		Unallocated: []Track{},
		Options: ContentMatchesOptions{
			DefaultAudioTrackName: lo.ToPtr("audio1"),
			DefaultSubtitleTrack:  lo.ToPtr("sub1"),
		},
	}

	t.Run("valid content match", func(t *testing.T) {
		// Проверяется корректный случай: все треки и эпизоды совпадают, ошибки быть не должно
		newMatches := &ContentMatches{
			Matches: []ContentMatch{
				{
					Episode:     episode,
					Video:       &videoTrack,
					AudioTracks: []Track{audioTrack},
					Subtitles:   []Track{subtitleTrack},
				},
			},
			Unallocated: []Track{},
			Options: ContentMatchesOptions{
				DefaultAudioTrackName: lo.ToPtr("audio1"),
				DefaultSubtitleTrack:  lo.ToPtr("sub1"),
			},
		}
		err := s.ValidateContentMatch(oldMatches, newMatches)
		require.NoError(t, err)
	})

	t.Run("episode count mismatch", func(t *testing.T) {
		// Проверяется случай, когда количество эпизодов не совпадает — должна быть ошибка
		newMatches := &ContentMatches{
			Matches: []ContentMatch{},
			Options: oldMatches.Options,
		}
		err := s.ValidateContentMatch(oldMatches, newMatches)
		require.ErrorContains(t, err, "old and new contents matches do not match")
	})

	t.Run("episode season mismatch", func(t *testing.T) {
		// Проверяется случай, когда номер сезона в эпизоде не совпадает — должна быть ошибка
		newMatches := &ContentMatches{
			Matches: []ContentMatch{
				{
					Episode: EpisodeInfo{
						SeasonNumber:  2, // Измененный номер сезона
						EpisodeNumber: episode.EpisodeNumber,
						FullPath:      episode.FullPath,
						RelativePath:  episode.RelativePath,
					},
					Video:       &videoTrack,
					AudioTracks: []Track{audioTrack},
					Subtitles:   []Track{subtitleTrack},
				},
			},
			Options: oldMatches.Options,
		}
		err := s.ValidateContentMatch(oldMatches, newMatches)
		require.ErrorContains(t, err, "episode mismatch")
	})

	t.Run("episode number mismatch", func(t *testing.T) {
		// Проверяется случай, когда номер эпизода не совпадает — должна быть ошибка
		newMatches := &ContentMatches{
			Matches: []ContentMatch{
				{
					Episode: EpisodeInfo{
						SeasonNumber:  episode.SeasonNumber,
						EpisodeNumber: 2, // Измененный номер эпизода
						FullPath:      episode.FullPath,
						RelativePath:  episode.RelativePath,
					},
					Video:       &videoTrack,
					AudioTracks: []Track{audioTrack},
					Subtitles:   []Track{subtitleTrack},
				},
			},
			Options: oldMatches.Options,
		}
		err := s.ValidateContentMatch(oldMatches, newMatches)
		require.ErrorContains(t, err, "episode mismatch")
	})

	t.Run("episode full path mismatch", func(t *testing.T) {
		// Проверяется случай, когда полный путь к файлу эпизода не совпадает — должна быть ошибка
		newMatches := &ContentMatches{
			Matches: []ContentMatch{
				{
					Episode: EpisodeInfo{
						SeasonNumber:  episode.SeasonNumber,
						EpisodeNumber: episode.EpisodeNumber,
						FullPath:      "/media/ep2.mkv", // Измененный путь
						RelativePath:  episode.RelativePath,
					},
					Video:       &videoTrack,
					AudioTracks: []Track{audioTrack},
					Subtitles:   []Track{subtitleTrack},
				},
			},
			Options: oldMatches.Options,
		}
		err := s.ValidateContentMatch(oldMatches, newMatches)
		require.ErrorContains(t, err, "episode mismatch")
	})

	t.Run("episode relative path mismatch", func(t *testing.T) {
		// Проверяется случай, когда относительный путь к файлу эпизода не совпадает — должна быть ошибка
		newMatches := &ContentMatches{
			Matches: []ContentMatch{
				{
					Episode: EpisodeInfo{
						SeasonNumber:  episode.SeasonNumber,
						EpisodeNumber: episode.EpisodeNumber,
						FullPath:      episode.FullPath,
						RelativePath:  "ep2.mkv", // Измененный относительный путь
					},
					Video:       &videoTrack,
					AudioTracks: []Track{audioTrack},
					Subtitles:   []Track{subtitleTrack},
				},
			},
			Options: oldMatches.Options,
		}
		err := s.ValidateContentMatch(oldMatches, newMatches)
		require.ErrorContains(t, err, "episode mismatch")
	})

	t.Run("track type mismatch", func(t *testing.T) {
		// Проверяется случай, когда тип трека не совпадает с ожидаемым (аудио трек помечен как субтитры) — должна быть ошибка
		badAudio := audioTrack
		badAudio.Type = TrackTypeSubtitle
		newMatches := &ContentMatches{
			Matches: []ContentMatch{
				{
					Episode:     episode,
					Video:       &videoTrack,
					AudioTracks: []Track{badAudio},
					Subtitles:   []Track{subtitleTrack},
				},
			},
			Options: oldMatches.Options,
		}
		err := s.ValidateContentMatch(oldMatches, newMatches)
		require.ErrorContains(t, err, "track type mismatch")
	})

	t.Run("no video with audio", func(t *testing.T) {
		// Проверяется случай, когда есть аудиодорожка, но нет видео — должна быть ошибка
		newMatches := &ContentMatches{
			Matches: []ContentMatch{
				{
					Episode:     episode,
					Video:       nil,
					AudioTracks: []Track{audioTrack},
					Subtitles:   []Track{},
				},
			},
			Options: oldMatches.Options,
		}
		err := s.ValidateContentMatch(oldMatches, newMatches)
		require.ErrorContains(t, err, "episode has audio or subtitle tracks but no video assigned")
	})

	t.Run("default audio does not exist", func(t *testing.T) {
		// Проверяется случай, когда дефолтная аудиодорожка указана, но её нет в списке треков — должна быть ошибка
		newMatches := &ContentMatches{
			Matches: []ContentMatch{
				{
					Episode:     episode,
					Video:       &videoTrack,
					AudioTracks: []Track{audioTrack}, // Есть аудиодорожка, чтобы пройти проверку 'no video with audio tracks'
					Subtitles:   []Track{subtitleTrack},
				},
			},
			Unallocated: []Track{},
			Options: ContentMatchesOptions{
				DefaultAudioTrackName: lo.ToPtr("non-existent-audio"), // Указана несуществующая дорожка
				DefaultSubtitleTrack:  lo.ToPtr("sub1"),
			},
		}
		err := s.ValidateContentMatch(oldMatches, newMatches)
		require.ErrorContains(t, err, "default audio does not exist")
	})

	t.Run("video track type mismatch", func(t *testing.T) {
		// Проверяется случай, когда видео трек имеет неверный тип — должна быть ошибка
		badVideo := videoTrack
		badVideo.Type = TrackTypeAudio // Устанавливаем неверный тип
		newMatches := &ContentMatches{
			Matches: []ContentMatch{
				{
					Episode:     episode,
					Video:       &badVideo,
					AudioTracks: []Track{audioTrack},
					Subtitles:   []Track{subtitleTrack},
				},
			},
			Options: oldMatches.Options,
		}
		err := s.ValidateContentMatch(oldMatches, newMatches)
		require.ErrorContains(t, err, "track type mismatch")
	})

	t.Run("subtitle track type mismatch", func(t *testing.T) {
		// Проверяется случай, когда трек субтитров имеет неверный тип — должна быть ошибка
		badSubtitle := subtitleTrack
		badSubtitle.Type = TrackTypeVideo // Устанавливаем неверный тип
		newMatches := &ContentMatches{
			Matches: []ContentMatch{
				{
					Episode:     episode,
					Video:       &videoTrack,
					AudioTracks: []Track{audioTrack},
					Subtitles:   []Track{badSubtitle},
				},
			},
			Options: oldMatches.Options,
		}
		err := s.ValidateContentMatch(oldMatches, newMatches)
		require.ErrorContains(t, err, "track type mismatch")
	})
	t.Run("unallocated file", func(t *testing.T) {
		// Проверяется случай, когда в newMatches появляется новый трек, которого не было в oldMatches — должна быть ошибка
		newTrack := Track{
			Type: TrackTypeAudio,
			Name: lo.ToPtr("new_audio"),
			File: FileInfo{RelativePath: "new_audio.aac", FullPath: "/path/new_audio.aac"},
		}
		newMatches := &ContentMatches{
			Matches: []ContentMatch{
				{
					Episode:     episode,
					Video:       &videoTrack,
					AudioTracks: []Track{audioTrack, newTrack}, // Добавляем новый трек
					Subtitles:   []Track{subtitleTrack},
				},
			},
			Options: oldMatches.Options,
		}
		err := s.ValidateContentMatch(oldMatches, newMatches)
		require.ErrorContains(t, err, "episode has unallocated file")
	})

	t.Run("track name mismatch", func(t *testing.T) {
		// Проверяется случай, когда имя трека было изменено — должна быть ошибка
		changedAudio := audioTrack
		changedAudio.Name = lo.ToPtr("changed_name") // Меняем имя
		newMatches := &ContentMatches{
			Matches: []ContentMatch{
				{
					Episode:     episode,
					Video:       &videoTrack,
					AudioTracks: []Track{changedAudio},
					Subtitles:   []Track{subtitleTrack},
				},
			},
			Options: oldMatches.Options,
		}
		err := s.ValidateContentMatch(oldMatches, newMatches)
		require.ErrorContains(t, err, "track name mismatch")
	})

	t.Run("track type mismatch in track list", func(t *testing.T) {
		// Проверяется случай, когда тип трека был изменен (но путь остался прежним) — должна быть ошибка

		// Создаем копию audioTrack, чтобы не изменять исходный объект для других тестов
		changedAudio := audioTrack
		changedAudio.Type = TrackTypeVideo // Меняем тип

		// Создаем копию oldMatches, чтобы не изменять исходный объект
		localOldMatches := *oldMatches
		localOldMatches.Matches = []ContentMatch{
			{
				Episode:     episode,
				Video:       &videoTrack,
				AudioTracks: []Track{audioTrack}, // Здесь audioTrack все еще с правильным типом
				Subtitles:   []Track{subtitleTrack},
			},
		}

		newMatches := &ContentMatches{
			Matches: []ContentMatch{
				{
					Episode:     episode,
					Video:       &videoTrack,
					AudioTracks: []Track{changedAudio}, // А здесь уже с измененным
					Subtitles:   []Track{subtitleTrack},
				},
			},
			Options: oldMatches.Options,
		}

		// Сначала сработает проверка на тип аудиодорожки
		err := s.ValidateContentMatch(&localOldMatches, newMatches)
		require.ErrorContains(t, err, "track type mismatch")

		// Теперь проверим случай, когда ошибка именно "track mismatch"
		// Для этого нужно, чтобы трек с измененным типом прошел первые проверки
		// и попал в общий список newTracks для сравнения с oldTrackMap.
		changedAudioInNew := audioTrack
		changedAudioInNew.Type = TrackTypeSubtitle // Меняем тип на субтитры

		newMatchesForMismatch := &ContentMatches{
			Matches: []ContentMatch{
				{
					Episode:     episode,
					Video:       &videoTrack,
					AudioTracks: []Track{audioTrack},        // Оставляем корректный аудио трек
					Subtitles:   []Track{changedAudioInNew}, // Добавляем "бывший" аудио трек как субтитр
				},
			},
			Options: oldMatches.Options,
		}

		err = s.ValidateContentMatch(oldMatches, newMatchesForMismatch)
		require.ErrorContains(t, err, "track mismatch")
	})

	t.Run("track full path mismatch", func(t *testing.T) {
		// Проверяется случай, когда полный путь трека был изменен (но относительный путь остался прежним) — должна быть ошибка
		changedAudio := audioTrack
		changedAudio.File.FullPath = "/new/path/audio1.aac" // Меняем полный путь
		newMatches := &ContentMatches{
			Matches: []ContentMatch{
				{
					Episode:     episode,
					Video:       &videoTrack,
					AudioTracks: []Track{changedAudio},
					Subtitles:   []Track{subtitleTrack},
				},
			},
			Options: oldMatches.Options,
		}
		err := s.ValidateContentMatch(oldMatches, newMatches)
		require.ErrorContains(t, err, "track mismatch")
	})

	t.Run("no audio tracks with video", func(t *testing.T) {
		// Проверяется случай, когда есть видео, но нет ни одной аудиодорожки — должна быть ошибка
		newMatches := &ContentMatches{
			Matches: []ContentMatch{
				{
					Episode:     episode,
					Video:       &videoTrack,
					AudioTracks: []Track{}, // Убираем аудиодорожки
					Subtitles:   []Track{subtitleTrack},
				},
			},
			Options: ContentMatchesOptions{},
		}
		err := s.ValidateContentMatch(oldMatches, newMatches)
		require.ErrorContains(t, err, "no video with audio tracks")
	})

	t.Run("default subtitle does not exist", func(t *testing.T) {
		// Проверяется случай, когда дефолтные субтитры указаны, но их нет в списке треков — должна быть ошибка
		newMatches := &ContentMatches{
			Matches: []ContentMatch{
				{
					Episode:     episode,
					Video:       &videoTrack,
					AudioTracks: []Track{audioTrack},
					Subtitles:   []Track{}, // Нет субтитров
				},
			},
			Options: ContentMatchesOptions{
				DefaultAudioTrackName: lo.ToPtr("audio1"),
				DefaultSubtitleTrack:  lo.ToPtr("non-existent-subtitle"), // Указаны несуществующие субтитры
			},
		}
		err := s.ValidateContentMatch(oldMatches, newMatches)
		require.ErrorContains(t, err, "default subtitle does not exist")
	})
}
