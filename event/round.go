package event

import "fmt"

type RoundStarted struct {
	Round     int
	Quota     int64
	Exhausted int64
	Scale     int64
}

func (e RoundStarted) Describe() string {
	description := fmt.Sprintf("Round: %d, Quota: %s", e.Round, formatScaledValue(e.Quota, e.Scale))
	if e.Exhausted > 0 {
		description += ", Excess: " + formatScaledValue(e.Exhausted, e.Scale)
	}
	return description
}
