package events

import (
	"fmt"
	"github.com/shawntoffel/election"
)

type ExcessUpdated struct {
	Excess int64
}

func (e *ExcessUpdated) Process() election.Event {
	description := fmt.Sprintf("%d excess votes are available after distribution", e.Excess)

	return election.Event{description}
}
