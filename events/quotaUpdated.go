package events

import (
	"fmt"

	"github.com/shawntoffel/election"
)

type QuotaUpdated struct {
	Scale int64
	Quota int64
}

func (e *QuotaUpdated) Process() election.Event {
	description := fmt.Sprintf("Quota has been updated to: %s", formatScaledValue(e.Quota, e.Scale))

	return election.Event{Description: description}
}
