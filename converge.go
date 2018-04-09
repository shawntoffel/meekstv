package meekstv

import (
	"github.com/shawntoffel/meekstv/events"
)

func (m *meekStv) Converged() bool {
	converged := true

	candidates := m.Pool.Elected()

	for _, candidate := range candidates {
		converged = m.TryConverge(*candidate)

		m.AddEvent(&events.TriedToConverge{Success: converged})

		m.SettleWeight(*candidate)
	}

	return converged
}

func (m *meekStv) TryConverge(candidate MeekCandidate) bool {

	currentWeight := (m.Quota * m.Scale) / candidate.Votes

	upperBound := m.UpperWeightBound()
	lowerBound := m.LowerWeightBound()

	if currentWeight > upperBound || currentWeight < lowerBound {
		return false
	}

	return true
}

func (m *meekStv) SettleWeight(candidate MeekCandidate) {
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

func (m *meekStv) UpperWeightBound() int64 {
	bound := m.GetScaleBound()

	return m.Scale + bound
}

func (m *meekStv) LowerWeightBound() int64 {
	bound := m.GetScaleBound()

	return m.Scale - bound
}

func (m *meekStv) GetScaleBound() int64 {
	frac := int64(100000)

	bound := m.Scale / frac

	return bound
}
