package events

import (
	"fmt"
	"strings"

	"github.com/shawntoffel/election"
)

type LosingCandidatesExcluded struct {
	Names []string
}

func (e *LosingCandidatesExcluded) Process() election.Event {
	description := fmt.Sprintf("The following candidates were not included in any ballot have been excluded: %v", strings.Join(e.Names, ", "))

	return election.Event{Description: description}
}
