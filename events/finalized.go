package events

import (
	"fmt"
	"strings"

	"github.com/shawntoffel/election"
)

type Finalized struct {
	Elected []string
}

func (e *Finalized) Process() election.Event {
	description := fmt.Sprintf("The following candidates have been elected: %s.", strings.Join(e.Elected, ", "))

	return election.Event{Description: description}
}
