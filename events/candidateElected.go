package events

import (
	"fmt"
)

type Elected struct {
	Name string
	Rank int
}

func (e *Elected) Process() string {
	description := fmt.Sprintf("%s has reached the quota and is elected with rank %d.", e.Name, e.Rank)

	return description
}
