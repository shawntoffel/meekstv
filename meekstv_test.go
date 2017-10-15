package meekstv

import (
	"fmt"
	"github.com/shawntoffel/election"
	"testing"
)

func TestMeekStv(t *testing.T) {

	var config = election.Config{}

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

	mstv := NewMeekStv()

	mstv.Initialize(config)

	result, err := mstv.Count()

	if err != nil {
		t.Error(err.Error())
	}

	fmt.Println("Events:", len(result.Events))

	for _, e := range result.Events {
		fmt.Println(e.Description)
	}

	for _, c := range result.Candidates {
		fmt.Println(c.Name)
	}

	count := len(result.Candidates)
	expectedCount := 3

	if count != expectedCount {
		t.Errorf("Incorrect number of elected candidates. Expected: %d, Got: %d", expectedCount, count)
	}

}
