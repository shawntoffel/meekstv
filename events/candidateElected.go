package events

import (
	"fmt"
	"github.com/shawntoffel/election"
)

type Elected struct {
	Name string
}

func (e *Elected) Process() election.Event {
	description := fmt.Sprintf("Candidate %s has been elected.", e.Name)

	return election.Event{description}
}
