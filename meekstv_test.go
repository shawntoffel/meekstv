package meekstv

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"os"

	blt "github.com/shawntoffel/goblt"
)

// This is a BLT file combined with expected election results, separated by "---".
const testFileExtension = ".bltr"

func TestAllElections(t *testing.T) {
	files, err := os.ReadDir("testdata")
	if err != nil {
		t.Error(err)
	}

	for _, file := range files {
		if filepath.Ext(file.Name()) != testFileExtension {
			continue
		}

		testName := strings.TrimSuffix(file.Name(), testFileExtension)
		t.Run(testName, func(t *testing.T) {
			config, expected := loadTestData(t, file.Name())
			result, err := Count(config)
			if err != nil {
				t.Error(err)
			}

			verifyResult(t, result, expected)
		})
	}
}

func TestRepeatableElectionOrder(t *testing.T) {
	config, expected := loadTestData(t, "simple.bltr")
	config.DisableDetail = true

	for i := 0; i < 1000; i++ {
		result, err := Count(config)
		if err != nil {
			t.Error(err)
		}

		success := verifyResult(t, result, expected)
		if !success {
			t.Errorf("Failed on iteration: %d", i+1)
			break
		}
	}
}

func TestElectLimit(t *testing.T) {
	p := pool{
		&Candidate{Name: "A", State: Hopeful, Votes: 3},
		&Candidate{Name: "B", State: Hopeful, Votes: 2},
		&Candidate{Name: "C", State: Hopeful, Votes: 1},
	}

	m := meekSTV{pool: p}

	// Max can elect.
	m.Seats = 2

	// Attempt to elect more candidates than seats.
	electedCount := m.elect(p)
	assertEqualBecause(t, electedCount, m.Seats, "should not elect more candidates than seats")

	// Elected in order of most to least votes.
	assertEqual(t, p[0].Rank, 1)
	assertEqual(t, p[1].Rank, 2)
}

func TestExcludeLimit(t *testing.T) {
	p := pool{
		&Candidate{Name: "A", State: Elected, Votes: 3},
		&Candidate{Name: "B", State: Hopeful, Votes: 2},
		&Candidate{Name: "C", State: Hopeful, Votes: 1},
	}

	m := meekSTV{pool: p}
	m.Seats = 2

	// Attempt to exclude 2 candidates, leaving too few to fill remaining seats.
	excludedCount := m.exclude(p[1:], Excluded)
	assertEqualBecause(t, excludedCount, 1, "should leave enough candidates to fill remaining seats")
	assertEqualBecause(t, p[2].State.Has(Excluded), true, "should exclude by lowest votes first")
}

func BenchmarkMeekCount(b *testing.B) {
	config, _ := loadTestData(b, "winter_2022.bltr")
	config.DisableDetail = true

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Count(config)
	}
}

func loadTestData(t testing.TB, filename string) (Config, []string) {
	data, err := os.ReadFile(filepath.Join("testdata", filename))
	if err != nil {
		t.Errorf("failed to read test file %s: %s", filename, err.Error())
	}
	config := generateTestConfig(t, bytes.NewReader(data))
	expected := readExpectedResults(t, bytes.NewReader(data))
	return config, expected
}

func generateTestConfig(t testing.TB, r io.Reader) Config {
	bltElection, err := blt.NewParser(r).Parse()
	if err != nil {
		t.Errorf("failed to load BLT data: %s", err.Error())
	}

	ballots := Ballots{}
	for _, bltBallot := range bltElection.Ballots {
		order := []int{}

		// Equal preferences are not supported.
		for _, preference := range bltBallot.Preferences {
			order = append(order, preference[0])
		}

		ballots = append(ballots, Ballot{
			Count:       bltBallot.Count,
			Preferences: order,
		})
	}

	return Config{
		Seats:               bltElection.NumSeats,
		Ballots:             ballots,
		Candidates:          bltElection.Candidates,
		WithdrawnCandidates: bltElection.Withdrawn,
	}
}

func readExpectedResults(t testing.TB, r io.Reader) []string {
	expected := []string{}

	scanner := bufio.NewScanner(r)

	// Skip the BLT content and separator.
	for scanner.Scan() {
		if scanner.Text() == "---" {
			break
		}
	}

	for scanner.Scan() {
		expected = append(expected, scanner.Text())
	}

	err := scanner.Err()
	if err != nil {
		t.Error(err)
	}

	return expected
}

func verifyResult(t *testing.T, result Result, expected []string) bool {
	t.Helper()

	if result.Detail != nil {
		buf := bytes.Buffer{}
		err := result.Detail.WriteReport(&buf)
		if err != nil {
			t.Error(err)
		}
		t.Log(buf.String())
	}

	if !slices.Equal(result.Elected, expected) {
		t.Errorf("expected:\n%s", rankedNames(expected))
		t.Errorf("got:\n%s", rankedNames(result.Elected))
		return false
	}

	return true
}

func rankedNames(names []string) string {
	result := strings.Builder{}
	for i, name := range names {
		result.WriteString(fmt.Sprintf("%d: %s\n", i+1, name))
	}
	return result.String()
}

func assertEqual[T comparable](t *testing.T, got, expected T) {
	assertEqualBecause(t, got, expected, "")
}

func assertEqualBecause[T comparable](t *testing.T, got, expected T, because string) {
	t.Helper()
	if got != expected {
		t.Error("expected", expected, "got", got, because)
	}
}
