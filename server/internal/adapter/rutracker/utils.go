package rutracker

import (
	"fmt"
	"io"
	"math"
	"strconv"
	"strings"
	"unicode"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html/charset"
)

func readerDocument(body io.Reader) (*goquery.Document, error) {
	utf8Reader, err := charset.NewReader(body, "text/html")
	if err != nil {
		return nil, fmt.Errorf("failed to create charset reader: %v", err)
	}

	doc, err := goquery.NewDocumentFromReader(utf8Reader)
	if err != nil {
		return nil, fmt.Errorf("failed to create document: %v", err)
	}

	return doc, nil
}

func emptyTorrentResponse() *TorrentResponse {
	return &TorrentResponse{
		Results:      []Torrent{},
		Page:         1,
		TotalResults: 1,
	}
}

// convertFileSizeToBytes преобразует строку объема файла в байты
func convertFileSizeToBytes(sizeStr string) (uint64, error) {
	// Удаляем все нечисловые символы кроме точки и пробелов
	cleaned := strings.Map(func(r rune) rune {
		if unicode.IsDigit(r) || r == '.' || unicode.IsSpace(r) {
			return r
		}
		return -1 // удаляем символ
	}, sizeStr)

	// Разделяем число и единицу измерения
	parts := strings.Fields(cleaned)
	if len(parts) == 0 {
		return 0, fmt.Errorf("пустая строка после очистки")
	}

	// Парсим числовую часть
	numberStr := parts[0]
	value, err := strconv.ParseFloat(numberStr, 64)
	if err != nil {
		return 0, fmt.Errorf("ошибка парсинга числа '%s': %v", numberStr, err)
	}

	// Определяем множитель на основе единицы измерения
	multiplier := uint64(1)
	if len(parts) > 1 {
		unit := strings.ToUpper(parts[1])
		multiplier = getMultiplierFromUnit(unit)
	} else {
		// Если единица измерения не указана, пытаемся определить из исходной строки
		multiplier = detectMultiplierFromOriginalString(sizeStr)
	}

	// Вычисляем результат
	result := uint64(value * float64(multiplier))
	return result, nil
}

// getMultiplierFromUnit возвращает множитель для единицы измерения
func getMultiplierFromUnit(unit string) uint64 {
	switch strings.ToUpper(unit) {
	case "B", "BYTE", "BYTES":
		return 1
	case "KB", "KILOBYTE", "KILOBYTES":
		return 1 << 10
	case "MB", "MEGABYTE", "MEGABYTES":
		return 1 << 20
	case "GB", "GIGABYTE", "GIGABYTES":
		return 1 << 30
	case "TB", "TERABYTE", "TERABYTES":
		return 1 << 40
	case "PB", "PETABYTE", "PETABYTES":
		return 1 << 50
	case "EB", "EXABYTE", "EXABYTES":
		return 1 << 60
	default:
		return 1 // по умолчанию считаем байтами
	}
}

// detectMultiplierFromOriginalString пытается определить множитель из исходной строки
func detectMultiplierFromOriginalString(original string) uint64 {
	originalUpper := strings.ToUpper(original)

	switch {
	case strings.Contains(originalUpper, "KB"):
		return 1 << 10
	case strings.Contains(originalUpper, "MB"):
		return 1 << 20
	case strings.Contains(originalUpper, "GB"):
		return 1 << 30
	case strings.Contains(originalUpper, "TB"):
		return 1 << 40
	case strings.Contains(originalUpper, "PB"):
		return 1 << 50
	case strings.Contains(originalUpper, "EB"):
		return 1 << 60
	default:
		return 1
	}
}

// convertTextToUint32 преобразует текст в uint32, игнорируя пробелы и запятые
func convertTextToUint32(text string) (uint32, error) {
	// Удаляем все пробелы и запятые из строки
	cleaned := strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) || r == ',' {
			return -1 // удаляем символ
		}
		return r
	}, text)

	// Проверяем, что строка не пустая после очистки
	if cleaned == "" {
		return 0, fmt.Errorf("пустая строка после очистки")
	}

	// Преобразуем в uint64 сначала, чтобы проверить переполнение
	result, err := strconv.ParseUint(cleaned, 10, 32)
	if err != nil {
		return 0, fmt.Errorf("ошибка преобразования: %v", err)
	}

	return uint32(result), nil
}

// FormatBytesWithPrecision преобразует с указанной точностью
func formatBytesWithPrecision(bytes uint64, precision int) string {
	if bytes == 0 {
		return "0B"
	}

	units := []string{"B", "KB", "MB", "GB", "TB", "PB", "EB"}
	base := float64(1024)

	exponent := math.Floor(math.Log(float64(bytes)) / math.Log(base))
	if exponent < 0 {
		exponent = 0
	}
	if exponent > float64(len(units)-1) {
		exponent = float64(len(units) - 1)
	}

	value := float64(bytes) / math.Pow(base, exponent)
	unit := units[int(exponent)]

	if exponent == 0 {
		return fmt.Sprintf("%dB", bytes)
	}

	return fmt.Sprintf("%.*f%s", precision, value, unit)
}
