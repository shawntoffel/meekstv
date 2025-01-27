package meekstv

type Config struct {
	Seats               int
	Ballots             Ballots
	Candidates          []string
	WithdrawnCandidates []int
	DisableDetail       bool
}
