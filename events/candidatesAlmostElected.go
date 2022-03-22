package events

import (
	"fmt"
)

type AlmostElected struct {
	Name string
}

func (e *AlmostElected) Process() string {
	description := fmt.Sprintf("%s has reached the quota and is pending election.", e.Name)

	return description
}
