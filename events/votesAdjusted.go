package events

import (
	"fmt"

	"github.com/shawntoffel/election"
)

type VotesAdjusted struct {
	Scale    int64
	Name     string
	Existing int64
	Total    int64
}

func (e *VotesAdjusted) Process() election.Event {
	diff := e.Total - e.Existing

	formattedDiff := formatScaledValue(diff, e.Scale)
	formattedTotal := formatScaledValue(e.Total, e.Scale)

	description := fmt.Sprintf("%s received %s votes. Total: %s", e.Name, formattedDiff, formattedTotal)

	return election.Event{Description: description}
}
