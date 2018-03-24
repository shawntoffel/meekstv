package events

import (
	"github.com/shawntoffel/election"
)

type RemainingCandidatesExcluded struct {
	Names []string
}

func (e *RemainingCandidatesExcluded) Process() election.Event {
	description := "All remaining candidates have been excluded."

	return election.Event{Description: description}
}
