package meekstv

import (
	"github.com/shawntoffel/meekstv/events"
)

type MeekRound struct {
	Number     int
	Excess     int64
	Surplus    int64
	AnyElected bool
}

func (m *meekStv) doRound() {
	m.incrementRound()
	m.distributeVotes()
	m.updateExcessVotesForRound()
	m.updateQuota()
	m.updateSurplus()

	count := m.electEligibleCandidates()
	m.round().AnyElected = count > 0

	if m.electionFinished() {
		return
	}

	for _, candidate := range m.Pool.Candidates() {
		m.settleWeight(*candidate)
	}

	if m.round().AnyElected {
		return
	}

	m.excludeLowestCandidate()

	if !m.canExcludeMoreCandidates() {
		m.electAllHopefulCandidates()
	}
}

func (m *meekStv) incrementRound() {
	round := len(m.meekRounds) + 1
	m.meekRounds = append(m.meekRounds, &MeekRound{Number: round})

	m.AddEvent(&events.RoundStarted{Round: round})

	m.Pool.ZeroAllVotes()
}

func (m *meekStv) updateExcessVotesForRound() {
	exhausted := int64(m.Ballots.Total()) * m.Scale

	votes := int64(0)

	for _, c := range m.Pool.Candidates() {
		votes = votes + c.Votes
	}

	m.round().Excess = exhausted - votes
	m.AddEvent(&events.ExcessUpdated{Excess: m.round().Excess})
}

func (m *meekStv) canExcludeMoreCandidates() bool {
	return m.Pool.Count()-m.Pool.ExcludedCount() > m.NumSeats
}

func (m *meekStv) updateSurplus() {
	candidates := append(m.Pool.Elected(), m.Pool.Hopeful()...)
	round := m.round()

	for _, candidate := range candidates {
		if candidate.Votes > m.Quota {
			round.Surplus = round.Surplus + (candidate.Votes - m.Quota)
		}
	}
}

func (m *meekStv) findCandidatesToEliminate() MeekCandidates {

	hopefulVotes := int64(0)

	for _, c := range m.Pool.Hopeful() {
		hopefulVotes += c.Votes
	}

	if hopefulVotes == 0 && m.round().Surplus == 0 {
		return m.Pool.Hopeful()
	}

	return MeekCandidates{}
}

func (m *meekStv) round() *MeekRound {
	round := len(m.meekRounds)

	if round < 1 {
		return &MeekRound{}
	}

	return m.meekRounds[round-1]
}
