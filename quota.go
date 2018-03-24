package meekstv

import (
	"github.com/shawntoffel/meekstv/events"
)

func (m *meekStv) Droop() int64 {
	return m.Scale / (int64(m.NumSeats) + 1)
}

func (m *meekStv) UpdateQuota() {
	total := int64(m.Ballots.Total()) * m.Scale
	prevQuota := m.Quota

	m.Quota = (int64(total) - m.MeekRound.Excess) * m.Droop() / m.Scale

	scaleBound := m.GetScaleBound()

	if m.Quota < scaleBound {
		m.Quota = scaleBound
	}

	if prevQuota != m.Quota {
		m.AddEvent(&events.QuotaUpdated{Quota: m.Quota})
	}
}
