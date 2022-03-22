package events

type AllHopefulCandidatesExcluded struct{}

func (e *AllHopefulCandidatesExcluded) Process() string {
	description := "All hopeful candidates have been excluded."

	return description
}
