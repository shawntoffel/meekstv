package meekstv

import (
	"github.com/shawntoffel/election"
	"github.com/shawntoffel/meekstv/events"
)

func (m *meekStv) DistributeVotes() {
	m.MeekRound.Excess = 0

	cs := m.Pool.Candidates()

	for _, candidate := range cs {
		m.Pool.SetVotes(candidate.Id, 0)
	}

	for _, ballot := range m.Ballots {
		m.DistributeAmongstBallot(ballot)
	}
}

func (m *meekStv) DistributeAmongstBallot(ballot election.RolledUpBallot) {
	value := int64(ballot.Count) * m.Scale

	ended := false

	iter := ballot.Ballot.List.Front()

	for {
		candidate := m.Pool.Candidate(iter.Value.(string))

		if !ended && candidate.Weight > 0 {
			ended = candidate.Status == Hopeful

			value = m.DistributeCandidateVotes(*candidate, value)
		}

		if iter.Next() == nil {
			break
		}

		iter = iter.Next()
	}

	m.MeekRound.Excess = m.MeekRound.Excess + value
}

func (m *meekStv) DistributeCandidateVotes(meekCandidate MeekCandidate, remainder int64) int64 {
	if meekCandidate.Status == Hopeful {
		m.GiveVotesToCandidate(meekCandidate, remainder)

		return 0
	}

	votes := remainder * meekCandidate.Weight / m.Scale

	m.GiveVotesToCandidate(meekCandidate, votes)

	remaining := remainder * (m.Scale - meekCandidate.Weight) / m.Scale

	return remaining
}

func (m *meekStv) GiveVotesToCandidate(meekCandidate MeekCandidate, votes int64) {
	newVotes := meekCandidate.Votes + votes

	m.Pool.SetVotes(meekCandidate.Id, newVotes)

	m.AddEvent(&events.VotesAdjusted{meekCandidate.Name, newVotes})
}
