package meekstv

import "github.com/shawntoffel/meekstv/events"

func (m *meekStv) distributeVotes() {
	for _, ballot := range m.Ballots {
		remainder := m.Scale

		for _, pref := range ballot.Preferences {
			candidate := m.Pool.Candidate(pref)
			votes := remainder * candidate.Weight * int64(ballot.Count) / m.Scale

			m.giveVotesToCandidate(*candidate, votes)

			remainder = remainder * (m.Scale - candidate.Weight) / m.Scale

			if remainder == 0 {
				break
			}
		}
	}
}

func (m *meekStv) giveVotesToCandidate(meekCandidate MeekCandidate, votes int64) {
	if meekCandidate.Status == Excluded {
		return
	}

	m.Pool.AddVotes(meekCandidate.Id, votes)
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

	m.AddEvent(&events.WeightAdjusted{Name: candidate.Name, NewWeight: newWeight, Scale: m.Scale})
}
