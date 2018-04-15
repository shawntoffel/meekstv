package meekstv

import (
	"github.com/shawntoffel/meekstv/events"
)

func (m *meekStv) CalculateQuota() int64 {
	total := int64(m.Ballots.Total()) * m.Scale
	excess := m.MeekRound.Excess
	numSeats := int64(m.NumSeats)

	return ((total - excess) / (numSeats + 1)) + 1
}

func (m *meekStv) UpdateQuota() {
	prevQuota := m.Quota

	m.Quota = m.CalculateQuota()

	scaleBound := m.GetScaleBound()

	if m.Quota < scaleBound {
		m.Quota = scaleBound
	}

	if prevQuota != m.Quota {
		m.AddEvent(&events.QuotaUpdated{Quota: m.Quota})
	}
}
