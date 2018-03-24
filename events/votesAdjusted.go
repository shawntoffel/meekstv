package events

import (
	"fmt"
	"github.com/shawntoffel/election"
)

type VotesAdjusted struct {
	Name  string
	Votes int64
}

func (e *VotesAdjusted) Process() election.Event {
	description := fmt.Sprintf("%s now has %d votes.", e.Name, e.Votes)

	return election.Event{Description: description}
}
