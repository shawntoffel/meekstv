package meekstv

func (m *meekStv) PerformPreliminaryCount() {
	numCandidates := m.Pool.Count()
	numExcluded := m.Pool.ExcludedCount()

	if numCandidates <= (m.NumSeats + numExcluded) {
		m.ElectAllHopefulCandidates()
		m.ElectedAll = true
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

func (m *meekStv) ElectionFinished() bool {
	numElected := m.Pool.ElectedCount()
	return numElected >= m.NumSeats
}
