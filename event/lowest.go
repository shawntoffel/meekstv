package event

import "fmt"

type LowestCandidateExcluded struct {
	Name       string
	RandomUsed bool
}

func (e LowestCandidateExcluded) Describe() string {
	if e.RandomUsed {
		return fmt.Sprintf("[X] %s was tied for the lowest number of votes and randomly selected for exclusion.", e.Name)
	}
	return fmt.Sprintf("[X] %s has the lowest number of votes and is selected for exclusion.", e.Name)
}
