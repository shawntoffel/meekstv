package events

import (
	"fmt"
	"strings"
)

type CandidatesExcluded struct {
	Names []string
}

func (e *CandidatesExcluded) Process() string {
	description := fmt.Sprintf("The following candidates have been excluded: %v", strings.Join(e.Names, ", "))

	return description
}
