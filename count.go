package meekstv

import "github.com/shawntoffel/meekstv/events"

func (m *meekStv) PerformPreliminaryCount() {
	numCandidates := m.Pool.Count()
	numExcluded := m.Pool.ExcludedCount()

	if numCandidates <= (m.NumSeats + numExcluded) {
		m.ElectAllHopefulCandidates()
		m.ElectedAll = true
	}

	m.ExcludeZeroVoteCandidates()
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

func (m *meekStv) ExcludeZeroVoteCandidates() {
	included := make(map[string]bool)

	for _, ballot := range m.Ballots {
		iter := ballot.Ballot.List.Front()

		for {
			candidate := m.Pool.Candidate(iter.Value.(string))

			included[candidate.Id] = true

			if iter.Next() == nil {
				break
			}

			iter = iter.Next()
		}

	}

	excluded := []string{}
	for _, c := range m.Pool.Hopeful() {
		_, ok := included[c.Id]

		if !ok {
			excluded = append(excluded, c.Name)
		}
	}

	if len(excluded) < 1 {
		return
	}

	hopeful := len(m.Pool.Hopeful())

	if (hopeful - len(excluded)) > m.NumSeats {

		for _, id := range excluded {
			m.Pool.ExcludeByName(id)
		}
		m.AddEvent(&events.LosingCandidatesExcluded{Names: excluded})
	}
}
