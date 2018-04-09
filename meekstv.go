package meekstv

import (
	"github.com/shawntoffel/election"
	"github.com/shawntoffel/meekstv/events"
)

type MeekStv interface {
	election.Counter
}

type meekStv struct {
	election.CounterState
	Quota         int64
	Round         int
	NumSeats      int
	Ballots       election.RolledUpBallots
	Pool          Pool
	Precision     int
	Scale         int64
	MaxIterations int
	ElectedAll    bool
	MeekRound     MeekRound
	Seed          int64
}

func NewMeekStv() MeekStv {
	m := meekStv{}
	m.Pool = NewPool()

	return &m
}

func (m *meekStv) Initialize(config election.Config) error {
	m.SetupNumSeats(config)
	m.SetupPrecision(config)
	m.SetupScale(config)
	m.SetupBallots(config)
	m.SetupPool(config)
	m.SetupMaxIterations(config)
	m.SetupSeed(config)

	m.AddEvent(&events.Initialized{Config: config})

	return m.Error
}

func (m *meekStv) Count() (*election.Result, error) {
	m.PerformPreliminaryCount()

	for _, ballot := range m.Ballots {
		ballot.Weight = m.Scale
	}

	m.DoRound()
	/*
		for {
			if m.HasEnded() {
				break
			}

			m.DoRound()
		}
	*/
	m.Finalize()

	return m.Result()
}
