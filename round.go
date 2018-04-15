package meekstv

import (
	"github.com/shawntoffel/meekstv/events"
)

type MeekRound struct {
	Excess     int64
	AnyElected bool
}

func (m *meekStv) IncrementRound() {
	m.Round = m.Round + 1
	m.MeekRound = MeekRound{}

	m.AddEvent(&events.RoundStarted{Round: m.Round})
}

func (m *meekStv) DoRound() {

	for {
		m.UpdateQuota()
		m.IncrementRound()
		m.Pool.ZeroAllVotes()

		for _, ballot := range m.Ballots {
			remainder := m.Scale

			iter := ballot.Ballot.List.Front()

			for {
				candidate := m.Pool.Candidate(iter.Value.(string))

				votes := remainder * candidate.Weight * int64(ballot.Count) / m.Scale

				m.GiveVotesToCandidate(*candidate, votes)

				remainder = remainder * (m.Scale - candidate.Weight) / m.Scale

				if remainder == 0 || iter.Next() == nil {
					break
				}

				iter = iter.Next()
			}

		}

		exhausted := int64(m.Ballots.Total()) * m.Scale

		v := int64(0)

		for _, c := range m.Pool.Candidates() {
			v = v + c.Votes
		}

		exhausted = exhausted - v

		m.MeekRound.Excess = exhausted

		m.AddEvent(&events.ExcessUpdated{Excess: m.MeekRound.Excess})
		m.UpdateQuota()

		count := m.ElectEligibleCandidates()

		if m.ElectionFinished() {
			return
		}

		if count > 0 {
			m.MeekRound.AnyElected = true
		}

		for _, candidate := range m.Pool.Candidates() {
			m.SettleWeight(*candidate)
		}

		if m.RoundHasEnded() {
			break
		}

	}

	if !m.ElectionFinished() {
		m.ExcludeLowestCandidate()

		numCandidates := m.Pool.Count()
		numExcluded := m.Pool.ExcludedCount()

		if (numCandidates - numExcluded) == m.NumSeats {
			m.ElectAllHopefulCandidates()
		}
	}

}

func (m *meekStv) RoundHasEnded() bool {
	if !m.MeekRound.AnyElected {
		return true
	}

	numElected := m.Pool.ElectedCount()

	if numElected >= m.NumSeats {
		return true
	}

	return false
}
