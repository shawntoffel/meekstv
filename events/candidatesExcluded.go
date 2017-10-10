package events

import (
	"fmt"
	"github.com/shawntoffel/election"
)

type CandidatesExcluded struct {
	Names []string
}

func (e *CandidatesExcluded) Process() election.Event {
	description := fmt.Sprintf("The following candidates have been excluded: %v", e.Names)

	return election.Event{description}
}
