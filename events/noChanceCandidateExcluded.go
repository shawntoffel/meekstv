package events

import (
	"strings"
)

type NoChanceCandidatesExcluded struct {
	Names []string
}

func (e *NoChanceCandidatesExcluded) Process() string {
	return "Given the remaining surplus of votes, the following candidates will never reach the quota and are excluded: " + strings.Join(e.Names, ", ") + "."
}
