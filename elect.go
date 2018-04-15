package meekstv

import (
	"math/rand"

	"github.com/shawntoffel/meekstv/events"
)

func (m *meekStv) ElectEligibleCandidates() int {
	eligibleCount := m.FindEligibleCandidates()

	m.HandleMultiwayTie(eligibleCount)

	m.NewlyElectAllAlmostCandidates()

	m.ProcessNewlyElectedCandidates()

	return eligibleCount
}

func (m *meekStv) FindEligibleCandidates() int {
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

func (m *meekStv) ProcessNewlyElectedCandidates() {
	candidates := m.Pool.NewlyElected()

	for _, candidate := range candidates {
		m.Pool.Elect(candidate.Id)
		m.AddEvent(&events.Elected{Name: candidate.Name, Rank: candidate.Rank})
	}
}

func (m *meekStv) NewlyElectAllAlmostCandidates() {
	candidates := m.Pool.Almost()
	for _, candidate := range candidates {
		m.Pool.NewlyElect(candidate.Id)
		m.MeekRound.AnyElected = true
	}
}

func (m *meekStv) ElectAllHopefulCandidates() {
	m.Pool.ElectHopeful()
	m.AddEvent(&events.AllHopefulCandidatesElected{})
}

func (m *meekStv) HandleMultiwayTie(eligibleCount int) {

	count := eligibleCount

	for {
		tooMany := m.Pool.ElectedCount()+count > m.NumSeats

		if !tooMany {
			break
		}

		m.Pool.ExcludeHopeful()
		m.AddEvent(&events.AllHopefulCandidatesExcluded{})

		m.ExcludeLowestCandidate()
		count = count - 1
	}
}

func (m *meekStv) ExcludeLowestCandidate() {
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
