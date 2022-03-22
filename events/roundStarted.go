package events

import (
	"fmt"
)

type RoundStarted struct {
	Round int
}

func (e *RoundStarted) Process() string {
	description := fmt.Sprintf("Round %d has started.", e.Round)

	return description
}
