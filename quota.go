package meekstv

import (
	"github.com/shawntoffel/meekstv/events"
)

func (m *meekStv) calculateQuota() int64 {
	total := int64(m.Ballots.TotalCount()) * m.Scale
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

	m.AddEvent(&events.QuotaSummarized{Previous: prevQuota, Current: m.Quota, Scale: m.Scale})
}

func (m *meekStv) getScaleBound() int64 {
	return m.Scale / int64(100000)
}
