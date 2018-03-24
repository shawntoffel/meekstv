package meekstv

import (
	"testing"

	"github.com/shawntoffel/election"
)

func TestMeekStv(t *testing.T) {
	result, err := runMeekStv(generateTestConfig())

	if err != nil {
		t.Errorf(err.Error())
	}

	for _, e := range result.Events {
		t.Log(e.Description)
	}

	t.Log("Events:", len(result.Events))

	for _, c := range result.Candidates {
		t.Log(c.Name)
	}

	verifyMeekStvResults(result, t)
}

func TestElectionOrder(t *testing.T) {
	for i := 0; i < 1000; i++ {
		result, err := runMeekStv(generateTestConfig())

		if err != nil {
			t.Errorf(err.Error())
		}

		verifyMeekStvResults(result, t)
	}
}

func verifyMeekStvResults(result *election.Result, t *testing.T) {
	count := len(result.Candidates)

	expectedCount := 3

	if count != expectedCount {
		t.Errorf("Incorrect number of elected candidates. Expected: %d, Got: %d", expectedCount, count)
	}

	expected := []string{"Alice", "Bob", "Chris"}

	got := []string{}

	for _, candidate := range result.Candidates {
		got = append(got, candidate.Name)
	}

	orderMatches := verifyElectionOrder(expected, got)

	if !orderMatches {
		t.Errorf("Election order is incorrect. Expected: %v, Got: %v", expected, got)
	}
}

func verifyElectionOrder(expected []string, got []string) bool {
	if len(expected) != len(got) {
		return false
	}

	for i, value := range expected {
		if value != got[i] {
			return false
		}
	}

	return true
}

func generateTestConfig() election.Config {
	config := election.Config{}

	names := []string{"Alice", "Bob", "Chris", "Don", "Eric", "Frank"}

	for _, name := range names {
		c := election.Candidate{}
		c.Id = name
		c.Name = name

		config.Candidates = append(config.Candidates, c)
	}

	var ballots election.Ballots

	for i := 0; i < 28; i++ {
		var ballot = election.NewBallot()
		ballot.PushBack("Alice")
		ballot.PushBack("Bob")
		ballot.PushBack("Chris")
		ballots = append(ballots, ballot)
	}

	for i := 0; i < 26; i++ {
		var ballot = election.NewBallot()
		ballot.PushBack("Bob")
		ballot.PushBack("Alice")
		ballot.PushBack("Chris")
		ballots = append(ballots, ballot)
	}

	for i := 0; i < 3; i++ {
		var ballot = election.NewBallot()
		ballot.PushBack("Chris")
		ballots = append(ballots, ballot)
	}

	for i := 0; i < 2; i++ {
		var ballot = election.NewBallot()
		ballot.PushBack("Don")
		ballots = append(ballots, ballot)
	}

	var ballot = election.NewBallot()
	ballot.PushBack("Eric")
	ballots = append(ballots, ballot)

	config.Ballots = ballots
	config.WithdrawnCandidates = []string{"Frank"}

	config.NumSeats = 3
	config.Precision = 6
	config.Seed = 1

	return config
}

func runMeekStv(config election.Config) (*election.Result, error) {
	mstv := NewMeekStv()

	mstv.Initialize(config)

	return mstv.Count()
}
