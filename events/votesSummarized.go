package events

type VotesSummarized struct {
	Name     string
	Previous int64
	Current  int64
	Scale    int64
}

func (e *VotesSummarized) Process() string {
	formattedTotal := formatScaledValue(e.Current, e.Scale)
	diff := e.Current - e.Previous
	if diff == 0 {
		return "= " + e.Name + " votes remain the same. Total: " + formattedTotal
	}

	change := "received"
	total := "Total"
	prefix := "+"

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

	return prefix + " " + e.Name + " " + change + " " + formattedDiff + " " + vote + ". " + total + ": " + formattedTotal
}
