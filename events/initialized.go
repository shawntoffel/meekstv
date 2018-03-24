package events

import (
	"bytes"
	"fmt"

	"github.com/shawntoffel/election"
)

type Initialized struct {
	Config election.Config
}

func (e *Initialized) Process() election.Event {
	buffer := bytes.Buffer{}

	buffer.WriteString("A new Meek STV count has been created. ")
	buffer.WriteString(fmt.Sprintf("Candidates: %d", len(e.Config.Candidates)))
	buffer.WriteString(fmt.Sprintf(", Ballots: %d", len(e.Config.Ballots)))
	buffer.WriteString(fmt.Sprintf(", Seats: %d", e.Config.NumSeats))
	buffer.WriteString(fmt.Sprintf(", Precision: %d", e.Config.Precision))
	buffer.WriteString(fmt.Sprintf(", Seed: %d", e.Config.Seed))

	description := buffer.String()

	return election.Event{Description: description}
}
