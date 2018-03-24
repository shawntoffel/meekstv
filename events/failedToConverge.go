package events

import (
	"fmt"
	"github.com/shawntoffel/election"
)

type FailedToConverge struct {
	MaxIterations int
}

func (e *FailedToConverge) Process() election.Event {
	description := fmt.Sprintf("Failed to converge in %d iterations.", e.MaxIterations)

	return election.Event{Description: description}
}
