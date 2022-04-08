package events

import (
	"fmt"
	"strings"
)

type NoChanceCandidatesExcluded struct {
	Names []string
}

func (e *NoChanceCandidatesExcluded) Process() string {
	return fmt.Sprintf("The following candidates will never reach the quota and were exluded: %s", strings.Join(e.Names, ", "))
}
