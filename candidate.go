package meekstv

type CandidateState uint8

const (
	Hopeful CandidateState = 1 << iota
	Elected
	Excluded
	Lowest
	Random
	Hopeless
	Withdrawn
)

// Returns true if ALL states in s are active. Logical AND.
func (c CandidateState) Has(s CandidateState) bool {
	return c&s == s
}

// Returns true if ANY state in s is active. Logical OR.
func (c CandidateState) HasAny(s CandidateState) bool {
	return c&s != 0
}

func newCandidate(id int, name string, scale int64) *Candidate {
	return &Candidate{
		ID:             id,
		Name:           name,
		Weight:         scale,
		previousWeight: scale,
		State:          Hopeful,
	}
}

type Candidate struct {
	ID     int
	Name   string
	Weight int64
	Votes  int64
	State  CandidateState
	Rank   int

	truncatedName  string
	previousWeight int64
	previousVotes  int64
	previousState  CandidateState
}

func (c Candidate) elected() bool {
	return c.State.Has(Elected)
}

func (c Candidate) excluded() bool {
	return c.State.Has(Excluded)
}

func (c Candidate) newelyElected() bool {
	return c.elected() && !c.previousState.Has(Elected)
}
func (c Candidate) newelyExcluded(reason CandidateState) bool {
	return c.State.Has(Excluded|reason) && !c.previousState.Has(Excluded)
}

func (c *Candidate) addVotes(votes int64) {
	if c.excluded() {
		return
	}
	c.Votes += votes
}

func (c *Candidate) removeVotes() {
	c.previousState = c.State
	c.previousVotes = c.Votes
	c.Votes = 0
}

func (c *Candidate) settleWeight(quota int64, scale int64) {
	if c.Votes == 0 || !c.elected() {
		return
	}

	newWeight := (quota * c.Weight) / c.Votes
	if newWeight%c.Votes > 0 {
		newWeight += 1
	}

	// Ensure the weight cannot exceed 1.
	if newWeight > scale {
		newWeight = scale
	}

	c.previousWeight = c.Weight
	c.Weight = newWeight
}

func (c *Candidate) setElected(rank int) {
	c.setState(Elected)
	c.Rank = rank
}

func (c *Candidate) setExcluded(reason CandidateState) {
	c.setState(Excluded | reason)
	c.Weight = 0
}

func (c *Candidate) setState(s CandidateState) {
	c.previousState = c.State
	c.State = s
}
