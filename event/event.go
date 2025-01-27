package event

import (
	"strconv"
)

type Event interface {
	Describe() string
}

func formatScaledValue(scaled int64, scale int64) string {
	return strconv.FormatFloat(float64(scaled)/float64(scale), 'f', -1, 64)
}
