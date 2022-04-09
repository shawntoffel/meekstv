package events

import (
	"strconv"
)

func formatScaledValue(scaled int64, scale int64) string {
	value := float64(scaled) / float64(scale)
	return strconv.FormatFloat(value, 'f', -1, 64)
}

func formatDiff(current, prev int64, scale int64, name, desc string) string {
	diff := current - prev
	change := "decreased"
	if diff > 0 {
		change = "increased"
	} else {
		diff *= -1
	}

	formattedDiff := formatScaledValue(diff, scale)
	formattedCurrent := formatScaledValue(current, scale)

	return name + " " + desc + " " + change + " by " + formattedDiff + ". New " + desc + ": " + formattedCurrent + "."
}
