package meekstv

import (
	"math/rand"
	"sort"

	"github.com/shawntoffel/meekstv/events"
)

func (m *meekStv) electEligibleCandidates() int {
	eligibleCount := m.findEligibleCandidates()

	m.handleMultiwayTie(eligibleCount)
	m.newlyElectAllAlmostCandidates()
	m.processNewlyElectedCandidates()

	return eligibleCount
}

func (m *meekStv) findEligibleCandidates() int {
	count := 0
	candidates := m.Pool.Hopeful()
	for _, candidate := range candidates {
		if candidate.Votes >= m.Quota {
			count = count + 1
			m.Pool.SetAlmost(candidate.Id)

			m.AddEvent(&events.AlmostElected{Name: candidate.Name})
		}
	}

	return count
}

func (m *meekStv) processNewlyElectedCandidates() {
	candidates := m.Pool.NewlyElected()

	for _, candidate := range candidates {
		m.Pool.Elect(candidate.Id)
		m.AddEvent(&events.Elected{Name: candidate.Name, Rank: candidate.Rank})
	}
}

func (m *meekStv) newlyElectAllAlmostCandidates() {
	candidates := m.Pool.Almost()
	for _, candidate := range candidates {
		m.Pool.NewlyElect(candidate.Id)
		m.round().AnyElected = true
	}
}

func (m *meekStv) electAllHopefulCandidates() {
	m.Pool.ElectHopeful()
	m.AddEvent(&events.AllHopefulCandidatesElected{})
}

func (m *meekStv) handleMultiwayTie(eligibleCount int) {
	count := eligibleCount

	for {
		tooMany := m.Pool.ElectedCount()+count > m.NumSeats

		if !tooMany {
			break
		}

		m.Pool.ExcludeHopeful()
		m.AddEvent(&events.AllHopefulCandidatesExcluded{})

		m.excludeLowestCandidate()
		count = count - 1
	}
}

func (m *meekStv) excludeLowestCandidate() {
	lowestCandidates := m.Pool.Lowest()

	toExclude := lowestCandidates[0]

	randomUsed := false

	if len(lowestCandidates) > 1 {
		seed := rand.NewSource(m.Seed)
		r := rand.New(seed)
		i := r.Intn(len(lowestCandidates))
		toExclude = lowestCandidates[i]

		randomUsed = true
	}

	m.Pool.Exclude(toExclude.Id)
	m.AddEvent(&events.LowestCandidateExcluded{Name: toExclude.Name, RandomUsed: randomUsed})
}

func (m *meekStv) excludeAllNoChanceCandidates() int {
	toExclude := MeekCandidates{}

	totalSurplus := int64(0)
	for _, c := range m.Pool.Candidates() {
		if c.Votes > m.Quota {
			totalSurplus += (c.Votes - m.Quota)
		}
	}

	hopefuls := m.Pool.Hopeful()
	sort.Sort(ByVotes(hopefuls))

	for i := 0; i < len(hopefuls); i++ {
		tryExclude := hopefuls[0 : len(hopefuls)-i]

		if len(m.Pool.Elected())+len(m.Pool.Almost())+len(hopefuls)-len(tryExclude) < m.NumSeats {
			continue
		}

		totalVotes := int64(0)
		for _, c := range tryExclude {
			totalVotes += c.Votes
		}

		if i != 0 && totalVotes+totalSurplus >= hopefuls[len(hopefuls)-i].Votes {
			continue
		}

		for _, c := range tryExclude {
			toExclude = append(toExclude, c)
		}
	}

	for _, c := range toExclude {
		m.Pool.Exclude(c.Id)
	}

	if len(toExclude) > 0 {
		m.AddEvent(&events.NoChanceCandidatesExcluded{Names: toExclude.SortedNames()})
	}

	return len(toExclude)
}
