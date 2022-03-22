package events

type AllHopefulCandidatesElected struct{}

func (e *AllHopefulCandidatesElected) Process() string {
	description := "All hopeful candidates have been elected."

	return description
}
