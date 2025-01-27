package meekstv

import (
	"io"
	"slices"
	"strings"
	"text/template"

	"github.com/shawntoffel/meekstv/event"
)

type Result struct {
	Elected   []string
	Withdrawn []string
	Detail    *Detail
}

type Detail struct {
	Candidates int
	Ballots    int
	Seats      int
	Precision  int
	Elected    []string
	Withdrawn  []string
	Rounds     []Round
}

func (d Detail) WriteReport(wr io.Writer) error {
	parsed, err := template.
		New("report").
		Funcs(templateFunctions).
		Parse(reportTemplate)
	if err != nil {
		return err
	}

	return parsed.Execute(wr, d)
}

type Round struct {
	Number        int
	Quota         int64
	Exhausted     int64
	Candidates    []Candidate
	previousQuota int64
	scale         int64
	finished      bool
}

func (round Round) Describe() string {
	return event.RoundStarted{
		Round:     round.Number,
		Quota:     round.Quota,
		Exhausted: round.Exhausted,
		Scale:     round.scale,
	}.Describe()
}

func (round Round) Events() []event.Event {
	events := []event.Event{}

	for _, e := range round.voteSummaryEvents() {
		events = append(events, e)
	}

	for _, e := range round.electedCandidateEvents() {
		events = append(events, e)
	}

	if !round.finished {
		for _, e := range round.weightEvents() {
			events = append(events, e)
		}
	}

	for _, e := range round.hopelessCandidatesEvent() {
		events = append(events, e)
	}

	for _, e := range round.lowestCandidateEvents() {
		events = append(events, e)
	}

	return events
}

func (round Round) voteSummaryEvents() []event.VotesSummarized {
	var filterFunc = func(c Candidate) bool {
		return c.Votes > 0 || c.Votes == 0 && c.previousVotes > 0
	}

	var mapFunc = func(c Candidate) event.VotesSummarized {
		return event.VotesSummarized{
			Name:     c.truncatedName,
			Current:  c.Votes,
			Previous: c.previousVotes,
			Rank:     c.Rank,
			Scale:    round.scale,
			Elected:  c.previousState.Has(Elected),
			Excluded: c.previousState.Has(Excluded),
		}
	}

	return filterMap(round.Candidates, filterFunc, mapFunc)
}

func (round Round) electedCandidateEvents() []event.CandidateElected {
	var filterFunc = func(c Candidate) bool {
		return c.newelyElected()
	}

	var mapFunc = func(c Candidate) event.CandidateElected {
		return event.CandidateElected{
			Name: c.Name,
			Rank: c.Rank,
		}
	}

	return filterMap(round.Candidates, filterFunc, mapFunc)
}

func (round Round) weightEvents() []event.WeightAdjusted {
	var filterFunc = func(c Candidate) bool {
		return c.Weight > 0 && c.Weight != c.previousWeight
	}

	var mapFunc = func(c Candidate) event.WeightAdjusted {
		return event.WeightAdjusted{
			Name:     c.truncatedName,
			Current:  c.Weight,
			Previous: c.previousWeight,
			Scale:    round.scale,
		}
	}

	return filterMap(round.Candidates, filterFunc, mapFunc)
}

func (round Round) hopelessCandidatesEvent() []event.HopelessCandidatesExcluded {
	events := []event.HopelessCandidatesExcluded{}

	var filterFunc = func(c Candidate) bool {
		return c.newelyExcluded(Hopeless)
	}

	var mapFunc = func(c Candidate) string {
		return c.Name
	}

	names := filterMap(round.Candidates, filterFunc, mapFunc)
	slices.Sort(names)
	if len(names) > 0 {
		events = append(events, event.HopelessCandidatesExcluded{
			Names: names,
		})
	}

	return events
}

func (round Round) lowestCandidateEvents() []event.LowestCandidateExcluded {
	var filterFunc = func(c Candidate) bool {
		return c.newelyExcluded(Lowest)
	}

	var mapFunc = func(c Candidate) event.LowestCandidateExcluded {
		return event.LowestCandidateExcluded{
			Name:       c.Name,
			RandomUsed: c.State.Has(Random),
		}
	}

	return filterMap(round.Candidates, filterFunc, mapFunc)
}

func filterMap[T any](candidates []Candidate, filterFunc func(Candidate) bool, mapFunc func(Candidate) T) []T {
	result := []T{}

	for _, candidate := range candidates {
		if filterFunc(candidate) {
			result = append(result, mapFunc(candidate))
		}
	}

	return result
}

var templateFunctions = template.FuncMap{
	"add": func(a, b int) int {
		return a + b
	},
	"join": func(elems []string, sep string) string {
		return strings.Join(elems, sep)
	},
}

const reportTemplate = `MeekSTV
==================================================
Candidates: {{.Candidates}}, Ballots: {{.Ballots}}, Seats: {{.Seats}}, Precision: {{.Precision}}
{{ if .Withdrawn -}}
Withdrawn: {{join .Withdrawn ", "}}
{{ end -}}
{{ range .Rounds -}}
--------------------------------------------------
{{ .Describe }}
--------------------------------------------------
{{ range .Events -}}
{{ .Describe }}
{{ end -}}
{{ end -}}
==================================================
Finalized
==================================================
{{range $i, $name := .Elected -}}
{{add $i 1}}: {{$name}}
{{ end -}}
`
