package events

type RemainingCandidatesExcluded struct {
	Names []string
}

func (e *RemainingCandidatesExcluded) Process() string {
	description := "All remaining candidates have been excluded."

	return description
}
