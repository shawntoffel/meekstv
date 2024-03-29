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
	Title         string
	Quota         int64
	NumSeats      int
	Ballots       election.Ballots
	Pool          Pool
	Precision     int
	Scale         int64
	MaxIterations int
	ElectedAll    bool
	meekRounds    []*MeekRound
	Seed          int64
}

func NewMeekStv() MeekStv {
	return &meekStv{}
}

func (m *meekStv) Initialize(config election.Config) error {
	m.setupTitle(config)
	m.setupNumSeats(config)
	m.setupPrecision(config)
	m.setupScale(config)
	m.setupBallots(config)
	m.setupPool(config)
	m.setupMaxIterations(config)
	m.setupSeed(config)

	m.AddEvent(&events.Initialized{Config: config})

	return m.Error
}

func (m *meekStv) Count() (*election.Result, error) {
	m.performPreliminaryCount()

	for {
		m.doRound()

		if m.hasEnded() {
			break
		}
	}

	m.finalize()

	return m.result()
}

func (m *meekStv) performPreliminaryCount() {
	numCandidates := m.Pool.Count()
	numExcluded := m.Pool.ExcludedCount()

	if numCandidates <= (m.NumSeats + numExcluded) {
		m.electAllHopefulCandidates()
		m.ElectedAll = true
	}

	m.excludeZeroVoteCandidates()
}

func (m *meekStv) hasEnded() bool {
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

	if m.round().Number >= m.MaxIterations {
		m.AddEvent(&events.FailedToConverge{MaxIterations: m.MaxIterations})
		return true
	}

	return false
}

func (m *meekStv) electionFinished() bool {
	return m.Pool.ElectedCount() >= m.NumSeats
}

func (m *meekStv) excludeZeroVoteCandidates() {
	included := make(map[string]bool)

	for _, ballot := range m.Ballots {
		for _, pref := range ballot.Preferences {
			candidate := m.Pool.Candidate(pref)
			included[candidate.Id] = true
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
		for _, name := range excluded {
			m.Pool.ExcludeByName(name)
		}
		m.AddEvent(&events.LosingCandidatesExcluded{Names: excluded})
	}
}

func (m *meekStv) finalize() {
	m.Pool.ExcludeHopeful()
	m.AddEvent(&events.RemainingCandidatesExcluded{})

	elected := m.Pool.Elected().AsCandidates()
	names := make([]string, len(elected))
	for i, candidate := range elected {
		names[i] = candidate.Name
	}
	m.AddEvent(&events.Finalized{Elected: names})
}

func (m *meekStv) result() (*election.Result, error) {
	if m.Error != nil {
		return nil, m.Error
	}

	elected := m.Pool.Elected()

	result := election.Result{
		Title:      m.Title,
		NumSeats:   m.NumSeats,
		NumBallots: m.Ballots.TotalCount(),
		Precision:  m.Precision,
		Seed:       m.Seed,
		Elected:    elected.AsCandidates(),
		Summary:    m.Summary,
	}

	return &result, nil
}
