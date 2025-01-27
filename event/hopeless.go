package event

import (
	"fmt"
	"strings"
)

type HopelessCandidatesExcluded struct {
	Names []string
}

func (e HopelessCandidatesExcluded) Describe() string {
	var candidates string

	if len(e.Names) > 1 {
		candidates = fmt.Sprintf("%d candidates", len(e.Names))
	} else {
		candidates = "candidate"
	}

	return fmt.Sprintf("[X] The following %s will never reach the quota and have been selected for exclusion: %s", candidates, strings.Join(e.Names, ", "))
}
