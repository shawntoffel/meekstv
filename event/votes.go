package event

import (
	"fmt"
)

type VotesSummarized struct {
	Name     string
	Current  int64
	Previous int64
	Rank     int
	Scale    int64
	Elected  bool
	Excluded bool
}

func (e VotesSummarized) Describe() string {
	formattedTotal := formatScaledValue(e.Current, e.Scale)
	diff := e.Current - e.Previous

	status := ""
	if e.Elected {
		status = fmt.Sprintf(" (Elected %d)", e.Rank)
	} else if e.Excluded {
		status = " (Excluded)"
	}

	if diff == 0 {
		return "= " + e.Name + status + " votes remain the same. Total: " + formattedTotal + "."
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

	return prefix + " " + e.Name + status + " " + change + " " + formatScaledValue(diff, e.Scale) + " " + vote + ". " + total + ": " + formattedTotal + "."
}
