package events

type VotesAdjusted struct {
	Name    string
	Prev    int64
	Current int64
	Scale   int64
}

func (e *VotesAdjusted) Process() string {
	change := "received"
	total := "Total"
	prefix := "+"

	diff := e.Current - e.Prev

	if diff < 0 {
		diff *= -1
		change = "transferred"
		total = "Remaining"
		prefix = "-"
	}

	vote := "vote"
	if diff != e.Scale {
		vote += "s"
	}

	formattedDiff := formatScaledValue(diff, e.Scale)
	formattedTotal := formatScaledValue(e.Current, e.Scale)

	return prefix + " " + e.Name + " " + change + " " + formattedDiff + " " + vote + ". " + total + ": " + formattedTotal
}
