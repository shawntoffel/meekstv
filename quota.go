package meekstv

import (
	"github.com/shawntoffel/meekstv/events"
)

func (m *meekStv) Droop() int64 {
	return m.Scale / (int64(m.NumSeats) + 1)
}

func (m *meekStv) UpdateQuota() {
	total := m.Ballots.Total()

	prevQuota := m.Quota

	m.Quota = (int64(total) - m.MeekRound.Excess) * m.Droop() / m.Scale * m.Scale

	if prevQuota != m.Quota {
		m.AddEvent(&events.QuotaUpdated{m.Quota})
	}
}
