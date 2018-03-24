package events

import (
	"fmt"
	"github.com/shawntoffel/election"
)

type QuotaUpdated struct {
	Quota int64
}

func (e *QuotaUpdated) Process() election.Event {
	description := fmt.Sprintf("Quota has been updated to: %d", e.Quota)

	return election.Event{Description: description}
}
