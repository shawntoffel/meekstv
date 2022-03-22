package events

import (
	"fmt"
)

type WeightAdjusted struct {
	Scale     int64
	Name      string
	NewWeight int64
}

func (e *WeightAdjusted) Process() string {
	description := fmt.Sprintf("%s weight has been adjusted to %s", e.Name, formatScaledValue(e.NewWeight, e.Scale))

	return description
}
