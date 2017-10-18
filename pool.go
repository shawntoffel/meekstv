package meekstv

import (
	"github.com/shawntoffel/election"
	"sort"
)

type Storage map[string]*MeekCandidate

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
	ExcludeHopeful()
	SetWeight(id string, weight int64)
	ZeroAllVotes()
}

type pool struct {
	Storage Storage
}

func NewPool() Pool {
	p := pool{}
	p.Storage = make(Storage)
	return &p
}

func (p *pool) Candidate(id string) *MeekCandidate {
	return p.Storage[id]
}

func (p *pool) SetVotes(id string, votes int64) {
	candidate := p.Candidate(id)

	candidate.Votes = votes

	p.Storage[candidate.Id] = candidate
}

func (p *pool) SetWeight(id string, weight int64) {
	candidate := p.Candidate(id)

	candidate.Weight = weight

	p.Storage[candidate.Id] = candidate
}

func (p *pool) Candidates() MeekCandidates {
	candidates := MeekCandidates{}

	list := List(p.Storage)
	for _, candidate := range list {
		candidates = append(candidates, candidate)
	}

	return candidates
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
	return len(List(p.Storage))
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
	candidate := p.Candidate(id)

	candidate.Status = Elected

	p.Storage[candidate.Id] = candidate
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

	p.Storage[candidate.Id] = candidate
}

func (p *pool) SetAlmost(id string) {
	candidate := p.Candidate(id)

	candidate.Status = Almost

	p.Storage[candidate.Id] = candidate
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

		p.Storage[c.Id] = &meekCandidate
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
	p.Storage[candidate.Id] = candidate

	return p.Candidate(id)
}

func (p *pool) ZeroAllVotes() {
	for _, candidate := range p.Candidates() {
		p.SetVotes(candidate.Id, 0)
	}
}

func List(storage Storage) MeekCandidates {
	list := MeekCandidates{}

	for _, value := range storage {
		list = append(list, *value)
	}

	return list
}
