package events

import (
	"fmt"

	"github.com/shawntoffel/election"
)

type ExcessUpdated struct {
	Scale  int64
	Excess int64
}

func (e *ExcessUpdated) Process() election.Event {
	excess := formatScaledValue(e.Excess, e.Scale)
	description := fmt.Sprintf("%s excess votes are available after distribution", excess)

	return election.Event{Description: description}
}
