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

	m.IncrementRound()
	m.UpdateQuota()

	for _, ballot := range m.Ballots {
		remainder := m.Scale

		iter := ballot.Ballot.List.Front()

		for {
			candidate := m.Pool.Candidate(iter.Value.(string))

			votes := remainder * candidate.Weight * ballot.Weight / m.Scale
			m.GiveVotesToCandidate(*candidate, votes)

			remainder = remainder * (m.Scale - candidate.Weight) / m.Scale

			if remainder == 0 {
				break
			}

			if iter.Next() == nil {
				break
			}

			iter = iter.Next()
		}
	}

	m.ElectEligibleCandidates()

	/*

		for {
			m.IncrementRound()

			m.ComputeRound()

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
		}*/
}

func (m *meekStv) ComputeRound() {
	converged := false

	for i := 0; i < m.MaxIterations; i++ {

		m.DistributeVotes()

		m.UpdateQuota()

		converged = m.Converged()

		if converged {
			break
		}
	}

	if !converged {
		m.AddEvent(&events.FailedToConverge{MaxIterations: m.MaxIterations})
	}

	count := m.ElectEligibleCandidates()

	if count > 0 {
		m.MeekRound.AnyElected = true
	}
}

func (m *meekStv) Complete() {
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
