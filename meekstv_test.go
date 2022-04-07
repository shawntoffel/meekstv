package meekstv

import (
	"fmt"
	"io"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/shawntoffel/election"
)

func TestSimple(t *testing.T) {
	expected := buildExpected(
		"Alice",
		"Bob",
		"Chris",
	)

	testMeekStv(t, "simple.blt", expected, true)
}

func TestSummer2017(t *testing.T) {
	expected := buildExpected(
		"Youkoso Jitsuryoku Shijou Shugi no Kyoushitsu e",
		"Owarimonogatari (Ge)",
		"Made in Abyss",
		"Princess Principal",
		"Fate/Apocrypha",
	)

	testMeekStv(t, "summer_2017.blt", expected, true)
}

func TestWinter2018(t *testing.T) {
	expected := buildExpected(
		"Kokkoku",
		"Hakata Tonkotsu Ramens",
		"Darling in the Franxx",
	)

	testMeekStv(t, "winter_2018.blt", expected, true)
}

func TestWinter2019(t *testing.T) {
	expected := buildExpected(
		"Mob Psycho 100 II",
		"Dororo",
		"Yakusoku no Neverland",
		"Kouya no Kotobuki Hikoutai",
		"Mahou Shoujo Tokushusen Asuka",
	)

	testMeekStv(t, "winter_2019.blt", expected, true)
}

func TestSummer2020(t *testing.T) {
	expected := buildExpected(
		"Deca-Dence",
		"Re:Zero kara Hajimeru Isekai Seikatsu 2nd Season",
		"THE GOD OF HIGH SCHOOL",
		"Monster Musume no Oisha-san",
		"Uzaki-chan wa Asobitai!",
	)

	testMeekStv(t, "summer_2020.blt", expected, true)
}

func TestWinter2022(t *testing.T) {
	expected := buildExpected(
		"Sabikui Bisco",
		"Princess Connect! Re:Dive Season 2",
		"Shuumatsu no Harem",
		"Hakozume: Kouban Joshi no Gyakushuu",
		"Tensai Ouji no Akaji Kokka Saisei Jutsu",
	)

	testMeekStv(t, "winter_2022.blt", expected, true)
}

func TestConsistentElectionOrder(t *testing.T) {
	expected := buildExpected(
		"Alice",
		"Bob",
		"Chris",
	)

	config := generateTestConfig(t, "simple.blt")

	for i := 0; i < 1000; i++ {
		result, err := runMeekStv(config)
		if err != nil {
			t.Errorf(err.Error())
		}

		success := verifyResults(t, result, expected, false)
		if !success {
			t.Errorf("Failed on iteration: %d", i+1)
			break
		}
	}
}

func testMeekStv(t *testing.T, filename string, expected election.Candidates, log bool) bool {
	result, err := runMeekStv(generateTestConfig(t, filename))
	if err != nil {
		t.Errorf(err.Error())
	}

	return verifyResults(t, result, expected, log)
}

func generateTestConfig(t *testing.T, filename string) election.Config {
	data := loadFileData(filename)
	c, err := election.LoadConfigFromBlt(data)
	if err != nil {
		t.Errorf("failed to load config from file %s: %s", filename, err.Error())
	}
	c.Precision = 8
	c.Seed = 1
	return c
}

func verifyResults(t *testing.T, result *election.Result, e []election.Candidate, log bool) bool {
	if result == nil {
		t.Error("nil election result")
		return false
	}
	if log {
		for _, e := range result.Summary.Events {
			t.Logf("%s: %s", e.Type, e.Description)
		}

		t.Log("Events:", len(result.Summary.Events))
		t.Log("Rounds:", len(result.Summary.Rounds))

		for _, c := range result.Elected {
			t.Log(c.Rank, c.Name)
		}
	}

	if !result.Elected.Equals(e) {
		t.Errorf("expected:\n%v\n", printCandidate(e))
		t.Errorf("got:\n%v\n", printCandidate(result.Elected))

		return false
	}

	return true
}

func printCandidate(c []election.Candidate) string {
	result := strings.Builder{}
	for _, candidate := range c {
		result.WriteString(fmt.Sprintf("%d: %s\n", candidate.Rank, candidate.Name))
	}

	return result.String()
}

func buildExpected(names ...string) []election.Candidate {
	candidates := []election.Candidate{}

	for i, name := range names {
		candidates = append(candidates, election.Candidate{
			Id:   name,
			Name: name,
			Rank: i + 1,
		})
	}

	return candidates
}

func loadFileData(filename string) io.Reader {
	bytes, _ := ioutil.ReadFile("testdata/" + filename)
	return strings.NewReader(string(bytes))
}

func runMeekStv(config election.Config) (*election.Result, error) {
	mstv := NewMeekStv()

	mstv.Initialize(config)

	return mstv.Count()
}
