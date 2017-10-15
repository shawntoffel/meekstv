package events

import (
	"fmt"
	"github.com/shawntoffel/election"
)

type WeightAdjusted struct {
	Name      string
	NewWeight int64
}

func (e *WeightAdjusted) Process() election.Event {
	description := fmt.Sprintf("%s weight has been adjusted to %d", e.Name, e.NewWeight)

	return election.Event{description}
}
