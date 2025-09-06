package matchtvshow

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
)

var (
	videoExtensions     = []string{".mp4", ".mkv", ".avi", ".mov", ".wmv", ".flv", ".webm"}
	audioExtensions     = []string{".mka", ".mp3", ".flac", ".m4a", ".ogg", ".wav", ".ac3"}
	subtitlesExtensions = []string{".ass", ".srt", ".ssa", ".sub", ".txt"}

	// Паттерны для определения языка (только четкие указания языка)
	languagePatterns = map[string]*regexp.Regexp{
		"ru": regexp.MustCompile(`(rus|russian|рус)`),
		"en": regexp.MustCompile(`(eng|english|англ)`),
		"jp": regexp.MustCompile(`(jap|japanese|яп|япон)`),
	}

	// Паттерны для определения сезона и серии
	seasonEpisodePatterns = []*regexp.Regexp{
		// SXXEXX формат
		regexp.MustCompile(`[s](\d+)[\s\-_]*[e](\d+)`),
		// Season X Episode Y
		regexp.MustCompile(`season[\.\s]*(\d+)[\s\-_]*episode[\.\s]*(\d+)`),
		// Серия X из Y сезона (русский)
		regexp.MustCompile(`серия[\.\s]*(\d+)[\s\-_]*из[\.\s]*(\d+)[\s\-_]*сезона`),
		// Просто номер серии в конце
		regexp.MustCompile(`[\s\-_\[\]](\d{2})[\s\.\[\]]`),
		regexp.MustCompile(`[\s\-_\[\]](\d{2,})[^a-z0-9]`),
	}

	// Паттерны для определения сезона (вынесены в переменные)
	seasonPatterns = []*regexp.Regexp{
		regexp.MustCompile(`season[\.\s]*(\d+)`),
		regexp.MustCompile(`сезон[\.\s]*(\d+)`),
		regexp.MustCompile(`[s](\d+)`),
	}
)

type Service struct {
}

func NewService() *Service {
	return &Service{}
}

// extractSeasonAndEpisode извлекает номер сезона и серии из строки
func extractSeasonAndEpisode(filename string) (season uint8, episode int, found bool) {
	for _, pattern := range seasonEpisodePatterns {
		matches := pattern.FindStringSubmatch(filename)
		if len(matches) >= 3 {
			// Для паттернов, которые находят и сезон и серию
			season = parseUint8(matches[1])
			episode = int(parseUint8(matches[2]))
			return season, episode, true
		} else if len(matches) >= 2 {
			// Для паттернов, которые находят только серию
			episode = int(parseUint8(matches[1]))
			season = detectSeasonFromString(filename)
			return season, episode, true
		}
	}
	return 0, 0, false
}

// detectSeasonFromString пытается определить сезон из названия
func detectSeasonFromString(filename string) uint8 {
	for _, pattern := range seasonPatterns {
		matches := pattern.FindStringSubmatch(filename)
		if len(matches) >= 2 {
			return parseUint8(matches[1])
		}
	}

	// Если сезон явно не указан, предполагаем 1-й
	return 1
}

// extractTrackName извлекает название трека из пути
func extractTrackName(filePath string) string {
	// Берем имя родительской директории как название
	dir := filepath.Dir(filePath)
	parentDir := filepath.Base(dir)

	// Если это не корневая директория, используем её имя
	if parentDir != "." && parentDir != "/" {
		return parentDir
	}

	// Иначе берем имя файла без расширения
	filename := filepath.Base(filePath)
	return strings.TrimSuffix(filename, filepath.Ext(filename))
}

// parseUint8 преобразует строку в uint8
func parseUint8(s string) uint8 {
	var result uint8
	for _, char := range s {
		if char >= '0' && char <= '9' {
			result = result*10 + uint8(char-'0')
		}
	}
	return result
}

// isVideoFile проверяет, является ли файл видео
func isVideoFile(filename string) bool {
	ext := filepath.Ext(filename)
	for _, videoExt := range videoExtensions {
		if ext == videoExt {
			return true
		}
	}
	return false
}

// isAudioFile проверяет, является ли файл аудио
func isAudioFile(filename string) bool {
	ext := filepath.Ext(filename)
	for _, audioExt := range audioExtensions {
		if ext == audioExt {
			return true
		}
	}
	return false
}

// isSubtitlesFile проверяет, является ли файл субтитрами
func isSubtitlesFile(filename string) bool {
	ext := filepath.Ext(filename)
	for _, subExt := range subtitlesExtensions {
		if ext == subExt {
			return true
		}
	}
	return false
}

// detectLanguage определяет язык трека - УПРОЩЕННАЯ ВЕРСИЯ
func detectLanguage(filePath string) *string {
	// Преобразуем весь путь в lowercase для поиска
	lowerPath := strings.ToLower(filePath)

	for lang, pattern := range languagePatterns {
		if pattern.MatchString(lowerPath) {
			return &lang
		}
	}
	return nil
}

// MatchEpisodeFiles сопоставляет файлы с эпизодами
func (s *Service) MatchEpisodeFiles(torrentFiles []string) ([]Episode, error) {
	episodesMap := make(map[string]*Episode)

	for _, file := range torrentFiles {
		// Преобразуем в lowercase один раз в начале
		lowerFile := strings.ToLower(file)
		filename := filepath.Base(lowerFile)

		// Пытаемся определить сезон и серию
		season, episode, found := extractSeasonAndEpisode(lowerFile) // Используем весь путь
		if !found {
			continue
		}

		// Создаем ключ для мапы эпизодов
		key := fmt.Sprintf("S%02dE%02d", season, episode)

		if _, exists := episodesMap[key]; !exists {
			episodesMap[key] = &Episode{
				EpisodeNumber: episode,
				SeasonNumber:  season,
				AudioFiles:    []Track{},
				Subtitles:     []Track{},
			}
		}

		ep := episodesMap[key]

		// Определяем тип файла
		switch {
		case isVideoFile(filename):
			ep.VideoFile = file // Сохраняем оригинальное имя файла

		case isAudioFile(filename):
			language := detectLanguage(lowerFile) // Используем lowercase путь
			trackName := extractTrackName(file)   // Оригинальное имя для красивого отображения
			ep.AudioFiles = append(ep.AudioFiles, Track{
				Name:     trackName,
				Language: language,
				File:     file,
			})

		case isSubtitlesFile(filename):
			language := detectLanguage(lowerFile)
			trackName := extractTrackName(file)
			ep.Subtitles = append(ep.Subtitles, Track{
				Name:     trackName,
				Language: language,
				File:     file,
			})
		}
	}

	// Конвертируем мапу в слайс
	episodes := make([]Episode, 0, len(episodesMap))
	for _, episode := range episodesMap {
		episodes = append(episodes, *episode)
	}

	return episodes, nil
}
