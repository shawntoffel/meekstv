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
	description := ""
	if len(e.Elected) > 0 {
		description = fmt.Sprintf("The following candidates have been elected: %s.", strings.Join(e.Elected, ", "))
	} else {
		description = "No candidates have been elected."
	}

	return election.Event{Description: description}
}
