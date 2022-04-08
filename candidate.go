package meekstv

import (
	"sort"

	"github.com/shawntoffel/election"
)

type CandidateStatus string

const (
	Elected      CandidateStatus = "Elected"
	NewlyElected CandidateStatus = "NewlyElected"
	Hopeful      CandidateStatus = "Hopeful"
	Excluded     CandidateStatus = "Excluded"
	Almost       CandidateStatus = "Almost"
)

type MeekCandidates []*MeekCandidate

func (m MeekCandidates) SortedNames() []string {
	names := make([]string, len(m))
	for i, c := range m {
		names[i] = c.Name
	}
	sort.Strings(names)
	return names
}

func (m MeekCandidates) TotalVotes() int64 {
	total := int64(0)
	for _, c := range m {
		total += c.Votes
	}
	return total
}

type MeekCandidate struct {
	election.Candidate
	Status CandidateStatus
	Weight int64
	Votes  int64
}

type ByVotes MeekCandidates

func (c ByVotes) Len() int {
	return len(c)
}

func (c ByVotes) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

func (c ByVotes) Less(i, j int) bool {
	return c[i].Votes < c[j].Votes
}

type BySnapshotVotes []MeekCandidate

func (c BySnapshotVotes) Len() int {
	return len(c)
}

func (c BySnapshotVotes) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

func (c BySnapshotVotes) Less(i, j int) bool {
	return c[i].Votes > c[j].Votes
}

func (meekCandidate *MeekCandidate) AsCandidate() election.Candidate {
	c := election.Candidate{}
	c.Id = meekCandidate.Id
	c.Name = meekCandidate.Name
	c.Rank = meekCandidate.Rank

	return c
}

func (meekCandidates MeekCandidates) AsCandidates() election.Candidates {
	candidates := make(election.Candidates, len(meekCandidates))

	for i, meekCandidate := range meekCandidates {
		candidates[i] = meekCandidate.AsCandidate()
	}

	sorted := candidates
	sort.Sort(election.ByRank(sorted))
	return sorted
}
