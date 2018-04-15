package meekstv

import (
	"errors"

	"github.com/shawntoffel/election"
	"github.com/shawntoffel/math"
	"github.com/shawntoffel/meekstv/events"
)

func (m *meekStv) setupNumSeats(config election.Config) {
	m.NumSeats = config.NumSeats

	if m.NumSeats < 1 {
		err := errors.New("At least one seat is required for election.")

		m.Error = err
	}
}

func (m *meekStv) setupPrecision(config election.Config) {
	m.Precision = config.Precision
}

func (m *meekStv) setupBallots(config election.Config) {
	m.Ballots = config.Ballots.Rollup()
}

func (m *meekStv) setupScale(config election.Config) {
	if m.Precision == 0 {
		m.setupPrecision(config)
	}

	m.Scale = math.Pow64(10, int64(m.Precision))
}

func (m *meekStv) setupPool(config election.Config) {
	if m.Scale == 0 {
		m.setupScale(config)
	}

	m.Pool.AddNewCandidates(config.Candidates, m.Scale)
	m.excludeWithdrawnCandidates(config.WithdrawnCandidates)
}

func (m *meekStv) setupMaxIterations(config election.Config) {
	m.MaxIterations = 1000
}

func (m *meekStv) setupSeed(config election.Config) {
	m.Seed = config.Seed
}

func (m *meekStv) excludeWithdrawnCandidates(names []string) {
	excluded := []string{}

	for _, name := range names {
		candidates := m.Pool.ExcludeByName(name)

		for _, candidate := range candidates {

			excluded = append(excluded, candidate.Name)
		}
	}

	m.AddEvent(&events.CandidatesExcluded{Names: excluded})
}
