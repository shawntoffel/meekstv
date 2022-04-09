package events

import (
	"strings"
)

type CandidatesExcluded struct {
	Names []string
}

func (e *CandidatesExcluded) Process() string {
	return "The following candidates have been excluded: " + strings.Join(e.Names, ", ") + "."
}
