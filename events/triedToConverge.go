package events

import (
	"github.com/shawntoffel/election"
)

type TriedToConverge struct {
	Success bool
}

func (e *TriedToConverge) Process() election.Event {
	description := ""

	if e.Success {
		description = "Successfully converged."
	} else {
		description = "Failed to converge. Adjusting weights for further vote distribution."
	}

	return election.Event{Description: description}
}
