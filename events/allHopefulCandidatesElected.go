package events

import (
	"github.com/shawntoffel/election"
)

type AllHopefulCandidatesElected struct{}

func (e *AllHopefulCandidatesElected) Process() election.Event {
	description := "All hopeful candidates have been elected."

	return election.Event{Description: description}
}
