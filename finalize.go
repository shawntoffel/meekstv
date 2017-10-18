package meekstv

import (
	"github.com/shawntoffel/meekstv/events"
)

func (m *meekStv) Finalize() {
	m.Pool.ExcludeHopeful()

	m.AddEvent(&events.RemainingCandidatesExcluded{})
}
