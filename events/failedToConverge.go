package events

import (
	"fmt"
)

type FailedToConverge struct {
	MaxIterations int
}

func (e *FailedToConverge) Process() string {
	description := fmt.Sprintf("Failed to converge in %d iterations.", e.MaxIterations)

	return description
}
