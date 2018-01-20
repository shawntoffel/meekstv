package events

import (
	"fmt"
	"github.com/shawntoffel/election"
)

type Elected struct {
	Name string
	Rank int
}

func (e *Elected) Process() election.Event {
	description := fmt.Sprintf("%s has been elected with rank %d.", e.Name, e.Rank)

	return election.Event{description}
}
