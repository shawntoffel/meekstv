package meekstv

import (
	"slices"
	"sort"
	"strings"
	"unicode"
)

type pool []*Candidate

// Counts the number of candidates with ANY of the provided states.
func (p pool) Count(state CandidateState) int {
	count := 0
	for _, candidate := range p {
		if candidate.State.HasAny(state) {
			count++
		}
	}
	return count
}

func (p pool) Elected() pool {
	return p.filter(func(c *Candidate) bool {
		return c.elected()
	})
}

func (p pool) Hopeful() pool {
	return p.filter(func(c *Candidate) bool {
		return c.State.Has(Hopeful)
	})
}

// LowlyHopeful returns all hopeful candidates who have the lowest number of votes in the pool.
func (p pool) LowlyHopeful() pool {
	lowly := pool{}

	for _, candidate := range p.Hopeful().OrderByVotes() {
		if len(lowly) > 0 && candidate.Votes != lowly[0].Votes {
			break
		}
		lowly = append(lowly, candidate)
	}

	return lowly
}

func (p pool) Withdrawn() pool {
	return p.filter(func(c *Candidate) bool {
		return c.State.Has(Excluded | Withdrawn)
	})
}

func (p pool) Electable(quota int64) pool {
	return p.filter(func(c *Candidate) bool {
		return c.Votes > quota && c.State.Has(Hopeful)
	})
}

func (p pool) Candidates(ids []int) []*Candidate {
	result := make([]*Candidate, len(ids))
	for i, id := range ids {
		result[i] = p[id-1]
	}
	return result
}

func (p pool) SumVotes() int64 {
	total := int64(0)
	for _, candidate := range p {
		if !candidate.excluded() {
			total += candidate.Votes
		}
	}
	return total
}

func (p pool) ResetVotes() {
	for _, candidate := range p {
		candidate.removeVotes()
	}
}

func (p pool) Names() []string {
	names := make([]string, len(p))
	for i, c := range p {
		names[i] = c.Name
	}
	return names
}

func (p pool) OrderByVotes() pool {
	candidates := p.copy()
	sort.SliceStable(candidates, func(i, j int) bool {
		return candidates[i].Votes < candidates[j].Votes
	})
	return candidates
}

func (p pool) OrderByVotesDescending() pool {
	candidates := p.copy()
	sort.SliceStable(candidates, func(i, j int) bool {
		return candidates[i].Votes > candidates[j].Votes
	})
	return candidates
}

func (p pool) OrderByRank() pool {
	candidates := p.copy()
	sort.SliceStable(candidates, func(i, j int) bool {
		return candidates[i].Rank < candidates[j].Rank
	})
	return candidates
}

// Short-cut exclusion rule
// Voting matters, Issue 22, p. 9
//
// "...if two or more lowest candidates have a total number of
// votes that, together with the current surplus, is less than
// the votes of the candidate next above, it is safe to exclude
// them all at once, provided that enough would remain
// to fill all seats."
func (p pool) Hopeless(quota int64, numSeats int) pool {
	surplus := p.calculateSurplusVotes(quota)

	// Order hopefuls by votes from high to low.
	hopefuls := p.Hopeful().OrderByVotesDescending()

	// Ensure we leave enough candidates to fill all seats.
	start := numSeats - p.Count(Elected)

	for i := start; i < len(hopefuls); i++ {
		// Get all hopefuls with votes lower than hopeful[i].
		lower := hopefuls[i+1:]

		// Calculate the amount of votes any one lower candidate can possibly receive.
		possibleVotes := surplus + lower.SumVotes()

		// If hopeful[i] has more votes than any lower candidate can receive, none of the lower candidates can win.
		if hopefuls[i].Votes > possibleVotes {
			return lower
		}
	}

	return pool{}
}

func (p pool) Snapshot() []Candidate {
	snapshot := make([]Candidate, len(p))
	for i, c := range p {
		snapshot[i] = *c
	}
	return snapshot
}

func (p pool) calculateSurplusVotes(quota int64) int64 {
	surplus := int64(0)
	for _, candidate := range p.Elected() {
		if candidate.Votes > quota {
			surplus += (candidate.Votes - quota)
		}
	}
	return surplus
}

func (p pool) filter(match func(*Candidate) bool) pool {
	result := pool{}
	for _, candidate := range p {
		if match(candidate) {
			result = append(result, candidate)
		}
	}
	return result
}

func (p pool) copy() pool {
	n := make(pool, len(p))
	copy(n, p)
	return n
}

func newPool(config Config, scale int64) pool {
	p := pool{}

	for i, name := range config.Candidates {
		candidate := newCandidate(i, name, scale)

		if slices.Contains(config.WithdrawnCandidates, candidate.ID+1) {
			candidate.setExcluded(Withdrawn)
		}

		if !config.DisableDetail {
			const defaultTruncationLength = 25
			candidate.truncatedName = truncate(name, findMinTruncationLength(
				config.Candidates,
				name,
				defaultTruncationLength,
			))
		}

		p = append(p, candidate)
	}

	return p
}

// Find a minimum truncation length per name such that all names remain unique.
func findMinTruncationLength(names []string, name string, start int) int {
	for i := start; i <= len(name); i++ {
		if count(truncateAll(names, i), truncate(name, i)) == 1 {
			return i
		}
	}
	return len(name)
}

func count(names []string, name string) int {
	count := 0
	for _, n := range names {
		if strings.EqualFold(n, name) {
			count++
		}
	}
	return count
}

func truncateAll(names []string, length int) []string {
	result := make([]string, len(names))
	for i, s := range names {
		result[i] = truncate(s, length)
	}
	return result
}

// Truncate to length or the nearest word boundary < length.
func truncate(name string, length int) string {
	runes := []rune(name)
	if len(runes) <= length || length < 0 {
		return name
	}

	for i := length - 1; i >= 0; i-- {
		if unicode.IsSpace(runes[i]) {
			return string(runes[:i]) + "..."
		}
	}

	return string(runes[:length]) + "..."
}
