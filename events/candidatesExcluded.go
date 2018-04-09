package events

import (
	"fmt"
	"strings"

	"github.com/shawntoffel/election"
)

type CandidatesExcluded struct {
	Names []string
}

func (e *CandidatesExcluded) Process() election.Event {
	description := fmt.Sprintf("The following candidates have been excluded: %v", strings.Join(e.Names, ", "))

	return election.Event{Description: description}
}
