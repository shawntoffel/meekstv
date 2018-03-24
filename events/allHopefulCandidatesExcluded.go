package events

import (
	"github.com/shawntoffel/election"
)

type AllHopefulCandidatesExcluded struct{}

func (e *AllHopefulCandidatesExcluded) Process() election.Event {
	description := "All hopeful candidates have been excluded."

	return election.Event{Description: description}
}
