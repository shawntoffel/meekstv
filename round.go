package meekstv

import (
	"github.com/shawntoffel/meekstv/events"
)

type MeekRound struct {
	Excess     int64
	AnyElected bool
}

func (m *meekStv) DoRound() {
	m.IncrementRound()

	m.DistributeVotes()

	m.UpdateExcessVotesForRound()

	m.UpdateQuota()

	count := m.ElectEligibleCandidates()

	m.MeekRound.AnyElected = count > 0

	if m.ElectionFinished() {
		return
	}

	for _, candidate := range m.Pool.Candidates() {
		m.SettleWeight(*candidate)
	}

	if m.MeekRound.AnyElected {
		return
	}

	m.ExcludeLowestCandidate()

	if !m.CanExcludeMoreCandidates() {
		m.ElectAllHopefulCandidates()
	}
}

func (m *meekStv) IncrementRound() {
	m.Round = m.Round + 1
	m.MeekRound = MeekRound{}

	m.AddEvent(&events.RoundStarted{Round: m.Round})

	m.Pool.ZeroAllVotes()
}

func (m *meekStv) DistributeVotes() {
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
}

func (m *meekStv) UpdateExcessVotesForRound() {
	exhausted := int64(m.Ballots.Total()) * m.Scale

	votes := int64(0)

	for _, c := range m.Pool.Candidates() {
		votes = votes + c.Votes
	}

	m.MeekRound.Excess = exhausted - votes
	m.AddEvent(&events.ExcessUpdated{Excess: m.MeekRound.Excess})
}

func (m *meekStv) CanExcludeMoreCandidates() bool {
	return m.Pool.Count()-m.Pool.ExcludedCount() > m.NumSeats
}
