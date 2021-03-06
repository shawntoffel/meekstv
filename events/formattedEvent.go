package events

import "strconv"

func formatScaledValue(scaled int64, scale int64) string {
	value := float64(scaled) / float64(scale)

	return strconv.FormatFloat(value, 'f', -1, 64)
}
