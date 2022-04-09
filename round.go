package meekstv

import (
	"sort"

	"github.com/shawntoffel/election"
	"github.com/shawntoffel/meekstv/events"
)

type MeekRound struct {
	Number     int
	Excess     int64
	Surplus    int64
	AnyElected bool
	Snapshot   []MeekCandidate
}

func (m *meekStv) doRound() {
	m.incrementRound()
	m.distributeVotes()
	m.summarizeVotes()
	m.updateExcessVotesForRound()
	m.updateQuota()
	m.updateSurplus()

	count := m.electEligibleCandidates()
	m.round().AnyElected = count > 0
	m.round().Snapshot = m.Pool.Snapshot()
	m.summarizeRound()

	if m.electionFinished() {
		return
	}

	for _, candidate := range m.Pool.Candidates() {
		m.settleWeight(*candidate)
	}

	if m.round().AnyElected {
		return
	}

	excluded := m.excludeAllNoChanceCandidates()
	if excluded < 1 && m.canExcludeMoreCandidates() {
		m.excludeLowestCandidate()
	}

	if !m.canExcludeMoreCandidates() {
		m.electAllHopefulCandidates()
	}
}

func (m *meekStv) incrementRound() {
	roundNumber := len(m.meekRounds) + 1
	m.meekRounds = append(m.meekRounds, &MeekRound{Number: roundNumber})

	m.AddEvent(&events.RoundStarted{Round: roundNumber})

	m.Pool.ZeroAllVotes()
}

func (m *meekStv) summarizeRound() {
	round := m.round()

	roundSummary := election.RoundSummary{
		Number:     round.Number,
		Quota:      m.Quota,
		Excess:     round.Excess,
		Surplus:    round.Surplus,
		Candidates: m.Pool.CandidateSummaries(),
	}

	m.Summary.AddRound(roundSummary)
}

func (m *meekStv) summarizeVotes() {
	prev := m.previousRound()
	if prev == nil {
		candidates := m.Pool.Snapshot()
		sort.Sort(BySnapshotVotes(candidates))
		for _, c := range candidates {
			if c.Votes > 0 {
				m.AddEvent(&events.VotesSummarized{
					Name:    c.Name,
					Current: c.Votes,
					Scale:   m.Scale,
				})
			}
		}
		return
	}

	snapshot := prev.Snapshot
	sort.Sort(BySnapshotVotes(snapshot))

	for _, previous := range snapshot {
		current := m.Pool.Candidate(previous.Id)
		if current.Votes > 0 {
			m.AddEvent(&events.VotesSummarized{
				Name:     current.Name,
				Elected:  current.Status == Elected,
				Rank:     current.Rank,
				Previous: previous.Votes,
				Current:  current.Votes,
				Scale:    m.Scale,
			})
		}
	}
}

func (m *meekStv) updateExcessVotesForRound() {
	exhausted := int64(m.Ballots.TotalCount()) * m.Scale

	currentVotes := int64(0)
	for _, c := range m.Pool.Candidates() {
		currentVotes = currentVotes + c.Votes
	}

	m.round().Excess = exhausted - currentVotes
	if m.round().Excess > 0 {
		m.AddEvent(&events.ExcessAvailable{Scale: m.Scale, Excess: m.round().Excess})
	}
}

func (m *meekStv) canExcludeMoreCandidates() bool {
	return m.Pool.Count()-m.Pool.ExcludedCount() > m.NumSeats
}

func (m *meekStv) updateSurplus() {
	candidates := append(m.Pool.Elected(), m.Pool.Hopeful()...)
	round := m.round()

	for _, candidate := range candidates {
		if candidate.Votes > m.Quota {
			round.Surplus = round.Surplus + (candidate.Votes - m.Quota)
		}
	}
}

func (m *meekStv) round() *MeekRound {
	round := m.getRoundBack(1)
	if round == nil {
		return &MeekRound{}
	}
	return round
}

func (m *meekStv) previousRound() *MeekRound {
	return m.getRoundBack(2)
}

func (m *meekStv) getRoundBack(back int) *MeekRound {
	count := len(m.meekRounds)
	if count < back {
		return nil
	}

	return m.meekRounds[count-back]
}
