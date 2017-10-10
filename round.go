package meekstv

import (
	"github.com/shawntoffel/meekstv/events"
)

type MeekRound struct {
	Excess     int64
	AnyElected bool
}

func (m *meekStv) IncrementRound() {
	m.Round = m.Round + 1
	m.MeekRound = MeekRound{}

	m.AddEvent(&events.RoundStarted{m.Round})
}

func (m *meekStv) DoRound() {

}
