package events

import (
	"fmt"
	"github.com/shawntoffel/election"
)

type RoundStarted struct {
	Round int
}

func (e *RoundStarted) Process() election.Event {
	description := fmt.Sprintf("Round %d has started.", e.Round)

	return election.Event{Description: description}
}
