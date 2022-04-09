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
		}
	}

	return count
}

func (m *meekStv) processNewlyElectedCandidates() {
	snapshotCandidates := m.Pool.NewlyElected().Snapshot()
	sort.Sort(BySnapshotVotes(snapshotCandidates))

	for _, snapShotCandidate := range snapshotCandidates {
		candidate := m.Pool.Elect(snapShotCandidate.Id)
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
	hopefuls := m.Pool.Hopeful()
	sort.Sort(ByVotes(hopefuls))

	surplus := m.round().Surplus
	numElected := len(m.Pool.Elected())

	for i := 0; i < len(hopefuls); i++ {
		tryExclude := hopefuls[0 : len(hopefuls)-i]

		if numElected+len(hopefuls)-len(tryExclude) < m.NumSeats {
			continue
		}

		totalVotes := tryExclude.TotalVotes()

		if i != 0 && totalVotes+surplus >= hopefuls[len(hopefuls)-i].Votes {
			continue
		}

		toExclude = tryExclude
		break
	}

	count := len(toExclude)

	if count > 0 {
		m.Pool.ExcludeMany(toExclude)
		m.AddEvent(&events.NoChanceCandidatesExcluded{Names: toExclude.SortedNames()})
	}

	return count
}
