package events

import (
	"fmt"

	"github.com/shawntoffel/election"
)

type WeightAdjusted struct {
	Scale     int64
	Name      string
	NewWeight int64
}

func (e *WeightAdjusted) Process() election.Event {
	description := fmt.Sprintf("%s weight has been adjusted to %s", e.Name, formatScaledValue(e.NewWeight, e.Scale))

	return election.Event{Description: description}
}
