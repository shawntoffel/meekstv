package events

import (
	"fmt"
)

type LowestCandidateExcluded struct {
	Name       string
	RandomUsed bool
}

func (e *LowestCandidateExcluded) Process() string {

	description := ""
	if e.RandomUsed {
		description = fmt.Sprintf("%s was tied for lowest number of votes and was randomly selected for exclusion.", e.Name)
	} else {
		description = fmt.Sprintf("%s has the lowest number of votes and is excluded.", e.Name)
	}

	return description
}
