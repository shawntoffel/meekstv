package events

import (
	"fmt"
	"github.com/shawntoffel/election"
)

type AlmostElected struct {
	Name string
}

func (e *AlmostElected) Process() election.Event {
	description := fmt.Sprintf("%s has reached the quota and is pending election.", e.Name)

	return election.Event{description}
}
