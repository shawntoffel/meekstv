package event

import "fmt"

type CandidateElected struct {
	Name string
	Rank int
}

func (e CandidateElected) Describe() string {
	return fmt.Sprintf("[O] Elected %d: %s", e.Rank, e.Name)
}
