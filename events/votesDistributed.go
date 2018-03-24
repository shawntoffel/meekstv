package events

import (
	"github.com/shawntoffel/election"
)

type VotesDistributed struct{}

func (e *VotesDistributed) Process() election.Event {
	description := "Finished distributing votes for iteration."

	return election.Event{Description: description}
}
