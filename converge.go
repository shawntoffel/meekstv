package meekstv

import (
	"github.com/shawntoffel/meekstv/events"
)

func (m *meekStv) Converged() bool {
	converged := true

	candidates := m.Pool.Candidates()

	for _, candidate := range candidates {
		if candidate.Status == Elected {
			converged = m.TryConverge(candidate)
			m.SettleWeight(candidate)
		}
	}

	return converged
}

func (m *meekStv) TryConverge(candidate MeekCandidate) bool {

	currentWeight := (m.Quota * m.Scale) / candidate.Votes

	if currentWeight > m.UpperWeightBound() || currentWeight < m.LowerWeightBound() {
		return false
	}

	return false
}

func (m *meekStv) SettleWeight(candidate MeekCandidate) {
	newWeight := (m.Quota * candidate.Weight) / candidate.Votes

	remaineder := newWeight % candidate.Votes

	if remaineder > 0 {
		newWeight = newWeight + 1
	}

	if newWeight > m.Scale {
		newWeight = m.Scale
	}

	m.Pool.SetWeight(candidate.Id, newWeight)

	m.AddEvent(&events.WeightAdjusted{candidate.Name, newWeight})
}

func (m *meekStv) UpperWeightBound() int64 {
	//bound := m.GetScaleBound()

	return m.Scale + 1
}

func (m *meekStv) LowerWeightBound() int64 {
	//bound := m.GetScaleBound()

	return m.Scale - 1
}

func (m *meekStv) GetScaleBound() int64 {
	frac := int64(100000)

	bound := m.Scale / (frac * m.Scale) / m.Scale

	return bound
}
