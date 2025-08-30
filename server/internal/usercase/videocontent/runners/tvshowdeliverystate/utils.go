package tvshowdeliverystate

import (
	"fmt"
	"math"
)

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
