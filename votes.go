package meekstv

import "github.com/shawntoffel/meekstv/events"

func (m *meekStv) distributeVotes() {
	for _, ballot := range m.Ballots {
		remainder := m.Scale

		iter := ballot.Ballot.List.Front()

		for {
			candidate := m.Pool.Candidate(iter.Value.(string))

			votes := remainder * candidate.Weight * int64(ballot.Count) / m.Scale

			m.giveVotesToCandidate(*candidate, votes, remainder)

			remainder = remainder * (m.Scale - candidate.Weight) / m.Scale

			if remainder == 0 || iter.Next() == nil {
				break
			}

			iter = iter.Next()
		}
	}
}

func (m *meekStv) giveVotesToCandidate(meekCandidate MeekCandidate, votes int64, remainder int64) {
	if meekCandidate.Status == Excluded {
		return
	}

	oldVotes := meekCandidate.Votes
	newVotes := oldVotes + votes

	m.Pool.SetVotes(meekCandidate.Id, newVotes)

	if oldVotes != 0 && remainder < m.Scale || m.round().Number == 1 {
		m.AddEvent(&events.VotesAdjusted{Name: meekCandidate.Name, Existing: oldVotes, Total: newVotes})
	}
}

func (m *meekStv) settleWeight(candidate MeekCandidate) {
	if candidate.Votes == 0 {
		return
	}

	if candidate.Status != "Elected" {
		return
	}
	newWeight := (m.Quota * candidate.Weight) / candidate.Votes

	remainder := newWeight % candidate.Votes

	if remainder > 0 {
		newWeight = newWeight + 1
	}

	if newWeight > m.Scale {
		newWeight = m.Scale
	}

	m.Pool.SetWeight(candidate.Id, newWeight)

	m.AddEvent(&events.WeightAdjusted{Name: candidate.Name, NewWeight: newWeight})
}
