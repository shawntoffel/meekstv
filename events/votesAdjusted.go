package events

import (
	"fmt"

	"github.com/shawntoffel/election"
)

type VotesAdjusted struct {
	Name     string
	Existing int64
	Total    int64
}

func (e *VotesAdjusted) Process() election.Event {
	diff := e.Total - e.Existing
	description := fmt.Sprintf("%s received %d votes. Total: %d", e.Name, diff, e.Total)

	return election.Event{Description: description}
}
