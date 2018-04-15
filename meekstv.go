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

	for {
		m.DoRound()

		if m.HasEnded() {
			break
		}

	}

	m.Finalize()

	return m.Result()
}

func (m *meekStv) PerformPreliminaryCount() {
	numCandidates := m.Pool.Count()
	numExcluded := m.Pool.ExcludedCount()

	if numCandidates <= (m.NumSeats + numExcluded) {
		m.ElectAllHopefulCandidates()
		m.ElectedAll = true
	}

	m.ExcludeZeroVoteCandidates()
}

func (m *meekStv) HasEnded() bool {
	if m.Error != nil {
		return true
	}

	if m.ElectedAll {
		return true
	}

	numElected := m.Pool.ElectedCount()

	if numElected == m.NumSeats {
		return true
	}

	return m.Round >= m.MaxIterations
}

func (m *meekStv) ElectionFinished() bool {
	numElected := m.Pool.ElectedCount()
	return numElected >= m.NumSeats
}

func (m *meekStv) ExcludeZeroVoteCandidates() {
	included := make(map[string]bool)

	for _, ballot := range m.Ballots {
		iter := ballot.Ballot.List.Front()

		for {
			candidate := m.Pool.Candidate(iter.Value.(string))

			included[candidate.Id] = true

			if iter.Next() == nil {
				break
			}

			iter = iter.Next()
		}

	}

	excluded := []string{}
	for _, c := range m.Pool.Hopeful() {
		_, ok := included[c.Id]

		if !ok {
			excluded = append(excluded, c.Name)
		}
	}

	if len(excluded) < 1 {
		return
	}

	hopeful := len(m.Pool.Hopeful())

	if (hopeful - len(excluded)) > m.NumSeats {

		for _, id := range excluded {
			m.Pool.ExcludeByName(id)
		}
		m.AddEvent(&events.LosingCandidatesExcluded{Names: excluded})
	}
}

func (m *meekStv) Finalize() {
	m.Pool.ExcludeHopeful()

	m.AddEvent(&events.RemainingCandidatesExcluded{})

	names := []string{}

	elected := m.Pool.Elected().AsCandidates()

	for _, candidate := range elected {
		names = append(names, candidate.Name)
	}

	m.AddEvent(&events.Finalized{Elected: names})
}

func (m *meekStv) Result() (*election.Result, error) {

	if m.Error != nil {
		return nil, m.Error
	}

	elected := m.Pool.Elected()

	result := election.Result{}
	result.Events = m.Events
	result.Candidates = elected.AsCandidates()

	return &result, nil
}
