package events

import "fmt"

func formatScaledValue(scaled int64, scale int64) string {
	value := float64(scaled) / float64(scale)

	return fmt.Sprintf("%f", value)
}
