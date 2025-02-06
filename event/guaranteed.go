package event

import (
	"fmt"
	"strings"
)

type GuaranteedCandidatesElected struct {
	Names []string
}

func (e GuaranteedCandidatesElected) Describe() string {
	var candidates string

	if len(e.Names) > 1 {
		candidates = "candidates have"
	} else {
		candidates = "candidate has"
	}

	return fmt.Sprintf("[O] The only remaining %s been elected: %s", candidates, strings.Join(e.Names, ", "))
}
