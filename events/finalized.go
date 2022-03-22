package events

import (
	"fmt"
	"strings"
)

type Finalized struct {
	Elected []string
}

func (e *Finalized) Process() string {
	description := ""
	if len(e.Elected) > 0 {
		description = fmt.Sprintf("The following candidates have been elected: %s.", strings.Join(e.Elected, ", "))
	} else {
		description = "No candidates have been elected."
	}

	return description
}
