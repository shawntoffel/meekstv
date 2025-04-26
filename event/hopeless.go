package event

import (
	"fmt"
	"strings"
)

type HopelessCandidatesExcluded struct {
	Names []string
}

func (e HopelessCandidatesExcluded) Describe() string {
	if len(e.Names) == 1 {
		return fmt.Sprintf("[X] %s will never reach the quota and is selected for exclusion.", e.Names[0])
	}

	return fmt.Sprintf("[X] The following %d candidates will never reach the quota and are selected for exclusion: %s", len(e.Names), strings.Join(e.Names, ", "))
}
