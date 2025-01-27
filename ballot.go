package meekstv

type Ballot struct {
	Count       int
	Preferences []int
}

type Ballots []Ballot

func (b Ballots) TotalCount() int {
	total := 0
	for _, ballot := range b {
		total += ballot.Count
	}
	return total
}
