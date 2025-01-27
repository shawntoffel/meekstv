package meekstv

import (
	"bytes"
	"fmt"
	"io"
	"slices"
	"strings"
	"testing"

	"os"

	blt "github.com/shawntoffel/goblt"
)

func TestSimple(t *testing.T) {
	expected := []string{
		"Alice",
		"Bob",
		"Chris",
	}

	testMeekStv(t, "simple.blt", expected)
}

func TestSummer2017(t *testing.T) {
	expected := []string{
		"Made in Abyss",
		"Youkoso Jitsuryoku Shijou Shugi no Kyoushitsu e",
		"Owarimonogatari (Ge)",
		"Princess Principal",
		"Fate/Apocrypha",
	}

	testMeekStv(t, "summer_2017.blt", expected)
}

func TestFall2017(t *testing.T) {
	expected := []string{
		"Mahoutsukai no Yome",
		"Shoujo Shuumatsu Ryokou",
		"Kujira no Kora wa Sajou ni Utau",
		"Kino no Tabi  -the Beautiful World- the Animated Series",
		"Infini-T Force",
	}

	testMeekStv(t, "fall_2017.blt", expected)
}

func TestWinter2018(t *testing.T) {
	expected := []string{
		"Kokkoku",
		"Hakata Tonkotsu Ramens",
		"Darling in the Franxx",
	}

	testMeekStv(t, "winter_2018.blt", expected)
}

func TestSpring2018(t *testing.T) {
	expected := []string{
		"Golden Kamui",
		"Hinamatsuri",
		"Mahou Shoujo Ore",
		"Megalo Box",
		"Persona 5 The Animation",
	}

	testMeekStv(t, "spring_2018.blt", expected)
}

func TestSummer2018(t *testing.T) {
	expected := []string{
		"Shingeki no Kyojin 3",
		"BANANA FISH",
		"Asobi Asobase",
		"Satsuriku no Tenshi",
		"Shoujo☆Kageki Revue Starlight",
	}

	testMeekStv(t, "summer_2018.blt", expected)
}

func TestFall2018(t *testing.T) {
	expected := []string{
		"JoJo no Kimyou na Bouken: Ougon no Kaze",
		"Zombie Land Saga",
		"Golden Kamuy 2",
		"🍌🐟",
		"Goblin Slayer",
	}

	testMeekStv(t, "fall_2018.blt", expected)
}

func TestWinter2019(t *testing.T) {
	expected := []string{
		"Yakusoku no Neverland",
		"Mob Psycho 100 II",
		"Dororo",
		"Kouya no Kotobuki Hikoutai",
		"Mahou Shoujo Tokushusen Asuka",
	}

	testMeekStv(t, "winter_2019.blt", expected)
}

func TestSpring2019(t *testing.T) {
	expected := []string{
		"One Punch Man 2",
		"Kono Yo no Hate de Koi o Utau Shoujo YU-NO",
		"Kimetsu no Yaiba",
		"Shingeki no Kyojin 3 Part 2",
		"Sarazanmai",
	}

	testMeekStv(t, "spring_2019.blt", expected)
}
func TestSummer2019(t *testing.T) {
	expected := []string{
		"Vinland Saga",
		"Toaru Kagaku no Accelerator",
		"Cop Craft",
		"Lord El-Melloi II-sei no Jikenbo: Rail Zeppelin Grace note",
		"Joshikousei no Mudazukai",
	}

	testMeekStv(t, "summer_2019.blt", expected)
}

func TestFall2019(t *testing.T) {
	expected := []string{
		"PSYCHO-PASS 3",
		"No Guns Life",
		"BEASTARS",
		"Fate/Grand Order: Zettai Majuu Sensen Babylonia",
		"Watashi, Nouryoku wa Heikinchi de tte Itta yo ne!",
	}

	testMeekStv(t, "fall_2019.blt", expected)
}

func TestWinter2020(t *testing.T) {
	expected := []string{
		"Toaru Kagaku no Railgun T",
		"Magia Record: Mahou Shoujo Madoka☆Magica Gaiden",
		"Eizouken ni wa Te wo Dasu na!",
		"ID: INVADED",
		"BanG Dream! 3rd Season",
	}

	testMeekStv(t, "winter_2020.blt", expected)
}

func TestSpring2020(t *testing.T) {
	expected := []string{
		"Fugou Keiji: Balance:UNLIMITED",
		"LISTENERS",
		"Princess Connect! Re:Dive",
		"Shin Sakura Taisen the Animation",
		"Otome Game no Hametsu Flag shika Nai Akuyaku Reijou ni Tensei shiteshimatta…",
	}

	testMeekStv(t, "spring_2020.blt", expected)
}

func TestSummer2020(t *testing.T) {
	expected := []string{
		"Deca-Dence",
		"Re:Zero kara Hajimeru Isekai Seikatsu 2nd Season",
		"THE GOD OF HIGH SCHOOL",
		"Monster Musume no Oisha-san",
		"Uzaki-chan wa Asobitai!",
	}

	testMeekStv(t, "summer_2020.blt", expected)
}

func TestFall2020(t *testing.T) {
	expected := []string{
		"Akudama Drive",
		"Munou na Nana",
		"Majo no Tabitabi",
		"Golden Kamuy 3",
		"Dragon Quest: Dai no Daibouken",
	}

	testMeekStv(t, "fall_2020.blt", expected)
}

func TestWinter2021(t *testing.T) {
	expected := []string{
		"Re:Zero kara Hajimeru Isekai Seikatsu 2nd Season Part 2",
		"Ura Sekai Picnic",
		"Wonder Egg Priority",
		"Yakusoku no Neverland 2",
		"Mushoku Tensei: Isekai Ittara Honki Dasu",
	}

	testMeekStv(t, "winter_2021.blt", expected)
}
func TestSpring2021(t *testing.T) {
	expected := []string{
		"Zombie Land Saga: Revenge",
		"Vivy: Fluorite Eye's Song",
		"86: Eighty Six",
		"Thunderbolt Fantasy: Touriken Yuuki 3",
		"Shadows House",
	}

	testMeekStv(t, "spring_2021.blt", expected)
}

func TestSummer2021(t *testing.T) {
	expected := []string{
		"Sonny Boy",
		"Magia Record: Mahou Shoujo Madoka☆Magica Gaiden 2nd Season - Kakusei Zenya",
		"NIGHT HEAD 2041",
		"Bokutachi no Remake",
		"Shiroi Suna no Aquatope",
	}

	testMeekStv(t, "summer_2021.blt", expected)
}

func TestFall2021(t *testing.T) {
	expected := []string{
		"JoJo no Kimyou na Bouken: Stone Ocean",
		"86: Eighty Six Part 2",
		"Sakugan",
		"takt op.Destiny",
		"Taishou Otome Otogibanashi",
	}

	testMeekStv(t, "fall_2021.blt", expected)
}

func TestWinter2022(t *testing.T) {
	expected := []string{
		"Sabikui Bisco",
		"Tensai Ouji no Akaji Kokka Saisei Jutsu",
		"Shuumatsu no Harem",
		"Princess Connect! Re:Dive Season 2",
		"Tokyo 24-ku",
	}

	testMeekStv(t, "winter_2022.blt", expected)
}

func TestRepeatableElectionOrder(t *testing.T) {
	expected := []string{
		"Alice",
		"Bob",
		"Chris",
	}

	config := generateTestConfig(t, "simple.blt")

	for i := 0; i < 1000; i++ {
		result, err := runMeekStv(config)
		if err != nil {
			t.Error(err)
		}

		success := verifyResults(t, result, expected)
		if !success {
			t.Errorf("Failed on iteration: %d", i+1)
			break
		}
	}
}

func testMeekStv(t *testing.T, filename string, expected []string) bool {
	result, err := runMeekStv(generateTestConfig(t, filename))
	if err != nil {
		t.Error(err)
	}

	if result.Detail == nil {
		t.Log("Detail is disabled.")
	} else {
		t.Log("Rounds:", len(result.Detail.Rounds))
		buf := bytes.Buffer{}
		err = result.Detail.WriteReport(&buf)
		if err != nil {
			t.Error(err)
		}
		t.Log(buf.String())
	}

	return verifyResults(t, result, expected)
}

func generateTestConfig(t *testing.T, filename string) Config {
	bltElection, err := blt.NewParser(loadFileData(filename)).Parse()
	if err != nil {
		t.Errorf("failed to load BLT file %s: %s", filename, err.Error())
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

func verifyResults(t *testing.T, result Result, elected []string) bool {
	if !slices.Equal(result.Elected, elected) {
		t.Errorf("expected:\n%v\n", rankedNames(elected))
		t.Errorf("got:\n%v\n", rankedNames(result.Elected))
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

func loadFileData(filename string) io.Reader {
	bytes, _ := os.ReadFile("testdata/" + filename)
	return strings.NewReader(string(bytes))
}

func runMeekStv(config Config) (Result, error) {
	return Count(config)
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
