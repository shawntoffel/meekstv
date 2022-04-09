package events

import (
	"strings"
)

type LosingCandidatesExcluded struct {
	Names []string
}

func (e *LosingCandidatesExcluded) Process() string {
	return "The following candidates were not included in any ballot and are excluded: " + strings.Join(e.Names, ", ") + "."
}
