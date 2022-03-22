package events

import (
	"fmt"
)

type ExcessUpdated struct {
	Scale  int64
	Excess int64
}

func (e *ExcessUpdated) Process() string {
	excess := formatScaledValue(e.Excess, e.Scale)
	description := fmt.Sprintf("%s excess votes are available after distribution", excess)

	return description
}
