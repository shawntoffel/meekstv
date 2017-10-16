package meekstv

import (
	"github.com/shawntoffel/election"
	"github.com/shawntoffel/math"
	"github.com/shawntoffel/meekstv/events"
)

func (m *meekStv) SetupNumSeats(config election.Config) {
	m.NumSeats = config.NumSeats

	if m.NumSeats < 0 {
		m.NumSeats = 0
	}
}

func (m *meekStv) SetupPrecision(config election.Config) {
	m.Precision = config.Precision
}

func (m *meekStv) SetupBallots(config election.Config) {
	m.Ballots = config.Ballots.Rollup()
}

func (m *meekStv) SetupScale(config election.Config) {
	if m.Precision == 0 {
		m.SetupPrecision(config)
	}

	m.Scale = math.Pow64(10, int64(m.Precision))
}

func (m *meekStv) SetupPool(config election.Config) {
	if m.Scale == 0 {
		m.SetupScale(config)
	}

	m.Pool.AddNewCandidates(config.Candidates, m.Scale)
	m.ExcludeWithdrawnCandidates(config.WithdrawnCandidates)
}

func (m *meekStv) SetupMaxIterations(config election.Config) {
	m.MaxIterations = 1000
}

func (m *meekStv) ExcludeWithdrawnCandidates(ids []string) {
	excluded := []string{}

	for _, id := range ids {
		candidate := m.Pool.Exclude(id)

		if candidate != nil {
			excluded = append(excluded, candidate.Name)
		}
	}

	m.AddEvent(&events.CandidatesExcluded{excluded})
}
