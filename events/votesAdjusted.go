package events

type VotesAdjusted struct {
	Name    string
	Prev    int64
	Current int64
	Scale   int64
}

func (e *VotesAdjusted) Process() string {
	diff := e.Current - e.Prev

	formattedTotal := formatScaledValue(e.Current, e.Scale)

	change := "received"
	total := "Total"
	prefix := "+"
	vote := "vote"

	if diff != e.Scale {
		vote += "s"
	}

	if diff < 0 {
		diff *= -1
		change = "transferred"
		total = "Remaining"
		prefix = "-"
	}

	formattedDiff := formatScaledValue(diff, e.Scale)

	return prefix + " " + e.Name + " " + change + " " + formattedDiff + " " + vote + ". " + total + ": " + formattedTotal
}
