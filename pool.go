package meekstv

import (
	"sort"

	"github.com/shawntoffel/election"
)

type Pool struct {
	MeekCandidates MeekCandidates
}

func (p *Pool) Candidate(id string) *MeekCandidate {
	for _, candidate := range p.MeekCandidates {
		if candidate.Id == id {
			return candidate
		}
	}

	return nil
}

func (p *Pool) AddVotes(id string, votes int64) {
	p.Candidate(id).Votes += votes
}

func (p *Pool) SetWeight(id string, weight int64) {
	p.Candidate(id).Weight = weight
}

func (p *Pool) NewlyElect(id string) {
	p.Candidate(id).Status = NewlyElected
}

func (p *Pool) SetAlmost(id string) {
	p.Candidate(id).Status = Almost
}

func (p *Pool) Candidates() MeekCandidates {
	return p.MeekCandidates
}

func (p *Pool) Snapshot() []MeekCandidate {
	candidates := p.MeekCandidates
	m := make([]MeekCandidate, len(candidates))
	for i, c := range candidates {
		m[i] = *c
	}
	return m
}

func (p *Pool) CandidatesWithStatus(status CandidateStatus) MeekCandidates {
	result := MeekCandidates{}
	for _, candidate := range p.Candidates() {
		if candidate.Status == status {
			result = append(result, candidate)
		}
	}
	return result
}

func (p *Pool) Count() int {
	return len(p.MeekCandidates)
}

func (p *Pool) ExcludedCount() int {
	return len(p.CandidatesWithStatus(Excluded))
}

func (p *Pool) Elected() MeekCandidates {
	return p.CandidatesWithStatus(Elected)
}

func (p *Pool) Hopeful() MeekCandidates {
	return p.CandidatesWithStatus(Hopeful)
}

func (p *Pool) NewlyElected() MeekCandidates {
	return p.CandidatesWithStatus(NewlyElected)
}

func (p *Pool) Almost() MeekCandidates {
	return p.CandidatesWithStatus(Almost)
}

func (p *Pool) ElectedCount() int {
	elected := len(p.CandidatesWithStatus(Elected))
	newlyElected := len(p.CandidatesWithStatus(NewlyElected))

	return elected + newlyElected
}

func (p *Pool) Elect(id string) {
	elected := len(p.Elected())

	candidate := p.Candidate(id)
	candidate.Status = Elected
	candidate.Rank = elected + 1
}

func (p *Pool) ElectAllNewlyElected() {
	for _, candidate := range p.Candidates() {
		if candidate.Status == NewlyElected {
			p.Elect(candidate.Id)
		}
	}
}

func (p *Pool) ElectHopeful() {
	for _, candidate := range p.Candidates() {
		if candidate.Status == Hopeful {
			p.Elect(candidate.Id)
		}
	}
}

func (p *Pool) ExcludeHopeful() {
	for _, candidate := range p.Candidates() {
		if candidate.Status == Hopeful {
			p.Exclude(candidate.Id)
		}
	}
}

func (p *Pool) Lowest() MeekCandidates {
	candidates := p.Candidates()
	sort.Sort(ByVotes(candidates))

	lowest := MeekCandidates{}

	for _, candidate := range candidates {
		if candidate.Status == Excluded {
			continue
		}
		if len(lowest) > 0 && candidate.Votes != lowest[0].Votes {
			break
		}
		lowest = append(lowest, candidate)
	}

	return lowest
}

func (p *Pool) AddNewCandidates(candidates election.Candidates, scale int64) {
	for _, c := range candidates {
		meekCandidate := MeekCandidate{
			Candidate: election.Candidate{
				Id:   c.Id,
				Name: c.Name,
			},
			Weight: 1 * scale,
			Status: Hopeful,
			Votes:  0,
		}

		p.MeekCandidates = append(p.MeekCandidates, &meekCandidate)
	}
}

func (p *Pool) Exclude(id string) *MeekCandidate {
	candidate := p.Candidate(id)
	if candidate == nil {
		return nil
	}

	candidate.Weight = 0
	candidate.Votes = 0
	candidate.Status = Excluded

	return candidate
}

func (p *Pool) ExcludeByName(name string) MeekCandidates {
	excluded := MeekCandidates{}
	for _, candidate := range p.Candidates() {
		if candidate.Name == name {
			c := p.Exclude(candidate.Id)
			excluded = append(excluded, c)
		}
	}

	return excluded
}

func (p *Pool) ZeroAllVotes() {
	for _, candidate := range p.Candidates() {
		candidate.Votes = 0
	}
}

func (p *Pool) CandidateSummaries() []election.CandidateSummary {
	candidates := p.Candidates()
	summaries := make([]election.CandidateSummary, len(candidates))

	for i, candidate := range candidates {
		summary := election.CandidateSummary{
			Candidate: candidate.AsCandidate(),
			Votes:     candidate.Votes,
			Weight:    candidate.Weight,
			Status:    string(candidate.Status),
		}

		summaries[i] = summary
	}

	return summaries
}
