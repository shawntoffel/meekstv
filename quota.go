package meekstv

import (
	"github.com/shawntoffel/meekstv/events"
)

func (m *meekStv) calculateQuota() int64 {
	total := int64(m.Ballots.Total()) * m.Scale
	excess := m.round().Excess
	numSeats := int64(m.NumSeats)

	return ((total - excess) / (numSeats + 1)) + 1
}

func (m *meekStv) updateQuota() {
	prevQuota := m.Quota

	m.Quota = m.calculateQuota()

	scaleBound := m.getScaleBound()

	if m.Quota < scaleBound {
		m.Quota = scaleBound
	}

	if prevQuota != m.Quota {
		m.AddEvent(&events.QuotaUpdated{Quota: m.Quota, Scale: m.Scale})
	}
}

func (m *meekStv) getScaleBound() int64 {
	frac := int64(100000)

	bound := m.Scale / frac

	return bound
}
