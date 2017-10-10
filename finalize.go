package meekstv

import (
	"github.com/shawntoffel/meekstv/events"
)

func (m *meekStv) Finalize() {
	m.ExcludeRemainingCandidates()
}

func (m *meekStv) ExcludeRemainingCandidates() {
	candidates := m.Pool.Candidates()

	for _, candidate := range candidates {
		if candidate.Status != Elected {
			m.Pool.Exclude(candidate.Id)
		}
	}

	m.AddEvent(&events.RemainingCandidatesExcluded{})
}
