package events

import (
	"fmt"
)

type ExcessAvailable struct {
	Scale  int64
	Excess int64
}

func (e *ExcessAvailable) Process() string {
	excess := formatScaledValue(e.Excess, e.Scale)
	description := fmt.Sprintf("%s excess votes are available after distribution", excess)

	return description
}
