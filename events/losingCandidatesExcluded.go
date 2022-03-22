package events

import (
	"fmt"
	"strings"
)

type LosingCandidatesExcluded struct {
	Names []string
}

func (e *LosingCandidatesExcluded) Process() string {
	description := fmt.Sprintf("The following candidates were not included in any ballot have been excluded: %v", strings.Join(e.Names, ", "))

	return description
}
