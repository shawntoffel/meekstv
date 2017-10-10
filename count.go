package meekstv

import (
	"github.com/shawntoffel/meekstv/events"
)

func (m *meekStv) PerformPreliminaryCount() {
	numCandidates := m.Pool.Count()
	numExcluded := m.Pool.ExcludedCount()

	if numCandidates <= (m.NumSeats + numExcluded) {
		m.Pool.ElectHopeful()
		m.ElectedAll = true

		m.AddEvent(&events.AllHopefulCandidatesElected{})
	}
}

func (m *meekStv) HasEnded() bool {
	if m.Error != nil {
		return true
	}

	if m.ElectedAll {
		return true
	}

	numElected := m.Pool.ElectedCount()

	return numElected == m.NumSeats
}
