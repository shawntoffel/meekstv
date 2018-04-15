package meekstv

import (
	"github.com/shawntoffel/meekstv/events"
)

type MeekRound struct {
	Excess     int64
	AnyElected bool
}

func (m *meekStv) doRound() {
	m.incrementRound()
	m.distributeVotes()
	m.updateExcessVotesForRound()
	m.updateQuota()

	count := m.electEligibleCandidates()
	m.MeekRound.AnyElected = count > 0

	if m.electionFinished() {
		return
	}

	for _, candidate := range m.Pool.Candidates() {
		m.settleWeight(*candidate)
	}

	if m.MeekRound.AnyElected {
		return
	}

	m.excludeLowestCandidate()

	if !m.canExcludeMoreCandidates() {
		m.electAllHopefulCandidates()
	}
}

func (m *meekStv) incrementRound() {
	m.Round = m.Round + 1
	m.MeekRound = MeekRound{}

	m.AddEvent(&events.RoundStarted{Round: m.Round})

	m.Pool.ZeroAllVotes()
}

func (m *meekStv) updateExcessVotesForRound() {
	exhausted := int64(m.Ballots.Total()) * m.Scale

	votes := int64(0)

	for _, c := range m.Pool.Candidates() {
		votes = votes + c.Votes
	}

	m.MeekRound.Excess = exhausted - votes
	m.AddEvent(&events.ExcessUpdated{Excess: m.MeekRound.Excess})
}

func (m *meekStv) canExcludeMoreCandidates() bool {
	return m.Pool.Count()-m.Pool.ExcludedCount() > m.NumSeats
}
