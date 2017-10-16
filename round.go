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

	m.AddEvent(&events.RoundStarted{m.Round})
}

func (m *meekStv) DoRound() {
	m.IncrementRound()

	converged := false

	for i := 0; i < m.MaxIterations; i++ {

		m.DistributeVotes()

		m.UpdateQuota()

		if m.Converged() {
			converged = true
			break
		}
	}

	if !converged {
		m.AddEvent(&events.FailedToConverge{m.MaxIterations})
	}

	count := m.ElectEligibleCandidates()

	if count > 0 {
		m.MeekRound.AnyElected = true
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
