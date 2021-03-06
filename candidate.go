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
type MeekCandidate struct {
	election.Candidate
	Status       CandidateStatus
	Weight       int64
	Votes        int64
	InitialVotes int64
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

func (meekCandidate *MeekCandidate) AsCandidate() election.Candidate {
	c := election.Candidate{}
	c.Id = meekCandidate.Id
	c.Name = meekCandidate.Name
	c.Rank = meekCandidate.Rank

	return c
}

func (meekCandidates MeekCandidates) AsCandidates() election.Candidates {
	candidates := election.Candidates{}

	for _, meekCandidate := range meekCandidates {
		candidate := meekCandidate.AsCandidate()
		candidates = append(candidates, candidate)
	}

	sorted := candidates

	sort.Sort(election.ByRank(sorted))

	return sorted
}
