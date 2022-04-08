package events

import (
	"fmt"
)

type PendingElection struct {
	Name string
}

func (e *PendingElection) Process() string {
	description := fmt.Sprintf("%s has reached the quota and is pending election.", e.Name)

	return description
}
