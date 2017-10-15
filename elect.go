package meekstv

import (
	"github.com/shawntoffel/meekstv/events"
	"math/rand"
	"time"
)

func (m *meekStv) ElectEligibleCandidates() {
	eligibleCount := m.FindEligibleCandidates()

	m.HandleMultiwayTie(eligibleCount)

	m.NewlyElectAllAlmostCandidates()

	m.AdjustNewlyElectedWeight()

	m.Pool.ElectAllNewlyElected()
}

func (m *meekStv) FindEligibleCandidates() int {
	count := 0
	candidates := m.Pool.Candidates()
	for _, candidate := range candidates {
		if candidate.Status == Hopeful && candidate.Votes > m.Quota {
			count = count + 1
			m.Pool.Almost(candidate.Id)

			m.AddEvent(&events.AlmostElected{candidate.Name})
		}
	}

	return count
}

func (m *meekStv) AdjustNewlyElectedWeight() {
	candidates := m.Pool.CandidatesWithStatus(NewlyElected)

	for _, candidate := range candidates {
		newWeight := (m.Quota * m.Scale) / candidate.Votes

		m.Pool.SetWeight(candidate.Id, newWeight)

		m.AddEvent(&events.WeightAdjusted{candidate.Name, newWeight})
	}
}

func (m *meekStv) NewlyElectAllAlmostCandidates() {
	candidates := m.Pool.CandidatesWithStatus(Almost)
	for _, candidate := range candidates {
		m.Pool.NewlyElect(candidate.Id)
		m.AddEvent(&events.Elected{candidate.Name})
	}
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
		seed := rand.NewSource(time.Now().Unix())
		r := rand.New(seed)
		i := r.Intn(len(lowestCandidates))
		toExclude = lowestCandidates[i]

		randomUsed = true
	}

	m.Pool.Exclude(toExclude.Id)
	m.AddEvent(&events.LowestCandidateExcluded{toExclude.Name, randomUsed})
}
