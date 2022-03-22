package meekstv

import (
	"sort"

	"github.com/shawntoffel/election"
)

type Pool interface {
	Candidate(id string) *MeekCandidate
	SetVotes(id string, votes int64)
	Lowest() MeekCandidates
	Candidates() MeekCandidates
	Count() int
	Elected() MeekCandidates
	Hopeful() MeekCandidates
	NewlyElected() MeekCandidates
	Almost() MeekCandidates
	ElectedCount() int
	ExcludedCount() int
	Elect(id string)
	ElectAllNewlyElected()
	NewlyElect(id string)
	SetAlmost(id string)
	ElectHopeful()
	AddNewCandidates(candidates election.Candidates, scale int64)
	Exclude(id string) *MeekCandidate
	ExcludeByName(name string) MeekCandidates
	ExcludeHopeful()
	SetWeight(id string, weight int64)
	ZeroAllVotes()
	CandidateSummaries() []election.CandidateSummary
}

type pool struct {
	MeekCandidates MeekCandidates
}

func NewPool() Pool {
	return &pool{}
}

func (p *pool) Candidate(id string) *MeekCandidate {
	for _, candidate := range p.MeekCandidates {
		if candidate.Id == id {
			return candidate
		}
	}

	return nil
}

func (p *pool) SetVotes(id string, votes int64) {
	candidate := p.Candidate(id)

	candidate.Votes = votes
}

func (p *pool) SetWeight(id string, weight int64) {
	candidate := p.Candidate(id)

	candidate.Weight = weight
}

func (p *pool) Candidates() MeekCandidates {
	return p.MeekCandidates
}

func (p *pool) CandidatesWithStatus(status CandidateStatus) MeekCandidates {
	candidates := p.Candidates()
	result := MeekCandidates{}

	for _, candidate := range candidates {
		if candidate.Status == status {
			result = append(result, candidate)
		}
	}

	return result
}

func (p *pool) Count() int {
	return len(p.MeekCandidates)
}

func (p *pool) ExcludedCount() int {
	return len(p.CandidatesWithStatus(Excluded))
}

func (p *pool) Elected() MeekCandidates {
	return p.CandidatesWithStatus(Elected)
}

func (p *pool) Hopeful() MeekCandidates {
	return p.CandidatesWithStatus(Hopeful)
}

func (p *pool) NewlyElected() MeekCandidates {
	return p.CandidatesWithStatus(NewlyElected)
}

func (p *pool) Almost() MeekCandidates {
	return p.CandidatesWithStatus(Almost)
}

func (p *pool) ElectedCount() int {
	elected := len(p.CandidatesWithStatus(Elected))
	newlyElected := len(p.CandidatesWithStatus(NewlyElected))

	return elected + newlyElected
}

func (p *pool) Elect(id string) {
	elected := len(p.Elected())

	candidate := p.Candidate(id)
	candidate.Status = Elected
	candidate.Rank = elected + 1
}

func (p *pool) ElectAllNewlyElected() {
	candidates := p.CandidatesWithStatus(NewlyElected)

	for _, candidate := range candidates {
		p.Elect(candidate.Id)
	}
}

func (p *pool) NewlyElect(id string) {
	candidate := p.Candidate(id)

	candidate.Status = NewlyElected
}

func (p *pool) SetAlmost(id string) {
	candidate := p.Candidate(id)

	candidate.Status = Almost
}

func (p *pool) ElectHopeful() {
	candidates := p.Candidates()

	for _, candidate := range candidates {
		if candidate.Status == Hopeful {
			p.Elect(candidate.Id)
		}
	}
}

func (p *pool) ExcludeHopeful() {
	candidates := p.Candidates()

	for _, candidate := range candidates {
		if candidate.Status == Hopeful {
			p.Exclude(candidate.Id)
		}
	}
}

func (p *pool) Lowest() MeekCandidates {
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

func (p *pool) AddNewCandidates(candidates election.Candidates, scale int64) {
	for _, c := range candidates {
		meekCandidate := MeekCandidate{}
		meekCandidate.Id = c.Id
		meekCandidate.Name = c.Name
		meekCandidate.Weight = 1 * scale
		meekCandidate.Status = Hopeful
		meekCandidate.Votes = 0

		p.MeekCandidates = append(p.MeekCandidates, &meekCandidate)
	}
}

func (p *pool) Exclude(id string) *MeekCandidate {
	candidate := p.Candidate(id)

	if candidate == nil {
		return nil
	}

	candidate.Weight = 0
	candidate.Votes = 0
	candidate.Status = Excluded

	return p.Candidate(id)
}

func (p *pool) ExcludeByName(name string) MeekCandidates {
	excluded := MeekCandidates{}
	for _, candidate := range p.Candidates() {
		if candidate.Name == name {
			c := p.Exclude(candidate.Id)
			excluded = append(excluded, c)
		}
	}

	return excluded
}

func (p *pool) ZeroAllVotes() {
	for _, candidate := range p.Candidates() {
		p.SetVotes(candidate.Id, 0)
	}
}

func (p *pool) CandidateSummaries() []election.CandidateSummary {
	summaries := []election.CandidateSummary{}

	for _, candidate := range p.Candidates() {
		summary := election.CandidateSummary{
			Candidate: candidate.AsCandidate(),
			Votes:     candidate.Votes,
			Weight:    candidate.Weight,
			Status:    string(candidate.Status),
		}

		summaries = append(summaries, summary)
	}

	return summaries
}
