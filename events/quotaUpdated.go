package events

import (
	"fmt"
)

type QuotaUpdated struct {
	Scale int64
	Quota int64
}

func (e *QuotaUpdated) Process() string {
	description := fmt.Sprintf("Quota has been updated to: %s", formatScaledValue(e.Quota, e.Scale))

	return description
}
