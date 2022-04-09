package events

type QuotaSummarized struct {
	Current  int64
	Previous int64
	Scale    int64
}

func (e *QuotaSummarized) Process() string {
	if e.Previous == 0 {
		return "Election quota set to " + formatScaledValue(e.Current, e.Scale) + "."
	}

	if e.Current == e.Previous {
		return "Election quota remains the same: " + formatScaledValue(e.Current, e.Scale) + "."
	}

	return formatDiff(e.Current, e.Previous, e.Scale, "Election", "quota")
}
