package meekstv

import (
	"testing"
)

func TestElectableCriteria(t *testing.T) {
	p := pool{
		&Candidate{Name: "A", State: Elected, Votes: 2},
		&Candidate{Name: "B", State: Hopeful, Votes: 2},
		&Candidate{Name: "C", State: Hopeful, Votes: 1},
		&Candidate{Name: "D", State: Excluded, Votes: 2},
		&Candidate{Name: "E", Votes: 2},
	}
	quota := int64(1)
	electable := p.Electable(quota)
	assertEqualBecause(t, len(electable), 1, "only 1 hopeful candidate exceeds the quota")
	assertEqual(t, electable[0].Name, "B")
}

func TestLowlyHopeful(t *testing.T) {
	p := pool{
		&Candidate{Name: "A", State: Hopeful, Votes: 2},
		&Candidate{Name: "B", State: Hopeful, Votes: 1},
		&Candidate{Name: "C", State: Hopeful, Votes: 1},
	}

	lowly := p.LowlyHopeful()
	assertEqualBecause(t, len(lowly), 2, "there are 2 candidates with the same lowest votes")

	// Order should be preserved.
	assertEqual(t, lowly[0].Name, "B")
	assertEqual(t, lowly[1].Name, "C")
}

func TestHopeless(t *testing.T) {
	p := pool{
		&Candidate{Name: "A", State: Elected, Votes: 5},
		&Candidate{Name: "B", State: Hopeful, Votes: 3},
		&Candidate{Name: "C", State: Hopeful, Votes: 2},
		&Candidate{Name: "D", State: Hopeful, Votes: 1},
		&Candidate{Name: "E", State: Hopeful, Votes: 0},
	}

	quota := int64(5)
	seats := 2

	hopeless := p.Hopeless(quota, seats)
	assertEqualBecause(t, len(hopeless), 2, "the sum of D and E votes is fewer than C's votes")

	// Order should be preserved.
	assertEqual(t, hopeless[0].Name, "D")
	assertEqual(t, hopeless[1].Name, "E")
}

func TestCandidatesByID(t *testing.T) {
	p := pool{
		&Candidate{ID: 1, Name: "A"},
		&Candidate{ID: 2, Name: "B"},
		&Candidate{ID: 3, Name: "C"},
	}

	candidates := p.Candidates([]int{3, 2})

	assertEqual(t, len(candidates), 2)
	assertEqual(t, candidates[0].Name, "C")
	assertEqual(t, candidates[1].Name, "B")
}

func TestSumVotes(t *testing.T) {
	p := pool{
		&Candidate{Name: "A", State: Elected, Votes: 3},
		&Candidate{Name: "B", State: Elected, Votes: 2},
		&Candidate{Name: "C", State: Excluded, Votes: 1},
	}

	votes := p.SumVotes()
	assertEqualBecause(t, votes, 5, "excluded candidates should not be included in the sum")
}

func TestTruncateWordBoundary(t *testing.T) {
	got := truncate("Toaru Kagaku no Railgun T", 20)
	expected := "Toaru Kagaku no..."
	assertEqual(t, got, expected)
}
