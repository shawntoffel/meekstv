package meekstv

import (
	"fmt"
	"math"
)

var (
	Precision     = 8
	MaxIterations = 1000
)

type meekSTV struct {
	Config
	pool      pool
	exhausted int64
	scale     int64
	random    WichmannHillRandom
	summaries []Round
}

func newMeekSTV(config Config) *meekSTV {
	scale := int64(math.Pow10(Precision))

	return &meekSTV{
		Config: config,
		scale:  scale,
		pool:   newPool(config, scale),
		random: NewWichmannHillRandom(
			len(config.Candidates),
			config.Seats+10000,
			config.Ballots.TotalCount()+20000,
		),
	}
}

func Count(config Config) (Result, error) {
	m := newMeekSTV(config)

	// Elect all hopeful candidates if there are enough seats.
	if m.Seats >= m.pool.Count(Hopeful) {
		m.elect(m.pool.Hopeful())
		return m.result(), nil
	}

	for i := 1; i <= MaxIterations; i++ {
		m.doRound(i)
		if m.pool.Count(Elected) >= m.Seats {
			break
		}
	}

	if m.pool.Count(Elected) < m.Seats {
		return Result{}, fmt.Errorf("failed to converge in %d iterations", MaxIterations)
	}

	return m.result(), nil
}

func (m *meekSTV) doRound(number int) {
	m.distributeVotes()

	// Calculate the quota each round using a variation of Droop.
	// quota = [(totalVotes – exhausted)/(seats + 1) + 1]
	// The last 1 is intended to be replaced by 1 in the last decimal place used.
	// Voting matters, Issue 1, p. 9
	quota := ((int64(m.Ballots.TotalCount())*m.scale)-m.exhausted)/int64(m.Seats+1) + 1

	// Elect any candidates who have reached the quota.
	electedCount := m.elect(m.pool.Electable(quota))

	// Adjust the weight of each candidate such that surplus votes (votes > quota) are transferred in the subsequent round.
	for _, candidate := range m.pool {
		candidate.settleWeight(quota, m.scale)
	}

	// Exclude lowest candidates if none were elected.
	if electedCount == 0 && m.canExcludeMoreCandidates() {
		// Exclude all candidates who will never reach the quota.
		excluded := m.exclude(m.pool.Hopeless(quota, m.Seats), Hopeless)

		// Otherwise, exclude a candidate with the lowest number of votes.
		if excluded == 0 {
			m.excludeLowestHopeful()
		}
	}

	// Elect remaining candidates if they are guaranteed to win.
	if m.pool.Count(Elected|Hopeful) == m.Seats {
		m.elect(m.pool.Hopeful())
	}

	m.summarizeRound(number, quota)
}

func (m *meekSTV) distributeVotes() {
	// Votes are counted from zero each round.
	m.pool.ResetVotes()
	m.exhausted = 0

	for _, ballot := range m.Ballots {
		// Each ballot is worth 1 vote.
		transferable := int64(ballot.Count) * m.scale

		for _, candidate := range m.pool.Candidates(ballot.Preferences) {
			// Excluded candidates transfer all votes.
			if candidate.excluded() {
				continue
			}

			// Give a portion of the votes to this candidate.
			candidate.addVotes(transferable * candidate.Weight / m.scale)

			// Transfer the remaining portion of the votes (1 - weight) to the next.
			transferable = (transferable * (m.scale - candidate.Weight)) / m.scale

			// Finished transferring votes.
			if transferable == 0 {
				break
			}
		}

		// No remaining candidates to receive these votes.
		if transferable > 0 {
			m.exhausted += transferable
		}
	}
}

func (m *meekSTV) canExcludeMoreCandidates() bool {
	return m.pool.Count(Elected|Hopeful) > m.Seats
}

// Elect candidates in order of most to least votes.
// Will stop electing if there are no remaining seats.
// Returns the number of elected candidates.
func (m *meekSTV) elect(p pool) int {
	nextRank := m.pool.Count(Elected) + 1

	for i, candidate := range p.OrderByVotesDescending() {
		if nextRank+i > m.Seats {
			return i
		}
		candidate.setElected(nextRank + i)
	}

	return len(p)
}

// Exclude candidates in order of least to most votes.
// Will stop excluding if remaining seats cannot be filled.
// Returns the number of excluded candidates.
func (m *meekSTV) exclude(p pool, reason CandidateState) int {
	remaining := m.pool.Count(Elected | Hopeful)

	for i, candidate := range p.OrderByVotes() {
		if remaining-i <= m.Seats {
			return i
		}

		candidate.setExcluded(reason)
	}

	return len(p)
}

func (m *meekSTV) excludeLowestHopeful() {
	lowly := m.pool.LowlyHopeful()

	if len(lowly) == 1 {
		lowly[0].setExcluded(Lowest)
		return
	}

	i := m.random.NextInt(len(lowly))
	lowly[i].setExcluded(Lowest | Random)
}

func (m *meekSTV) summarizeRound(number int, quota int64) {
	if m.DisableDetail {
		return
	}

	var previousQuota int64
	if len(m.summaries) > 0 {
		previousQuota = m.summaries[len(m.summaries)-1].Quota
	}

	m.summaries = append(m.summaries, Round{
		Number:        number,
		Quota:         quota,
		Exhausted:     m.exhausted,
		Candidates:    m.pool.OrderByVotesDescending().Snapshot(),
		previousQuota: previousQuota,
		scale:         m.scale,
		finished:      m.pool.Count(Elected) >= m.Seats,
	})
}

func (m *meekSTV) result() Result {
	result := Result{
		Elected:   m.pool.Elected().OrderByRank().Names(),
		Withdrawn: m.pool.Withdrawn().Names(),
	}

	if m.DisableDetail {
		return result
	}

	result.Detail = &Detail{
		Candidates: len(m.pool),
		Ballots:    m.Ballots.TotalCount(),
		Seats:      m.Seats,
		Precision:  Precision,
		Elected:    result.Elected,
		Withdrawn:  result.Withdrawn,
		Rounds:     m.summaries,
	}

	return result
}
