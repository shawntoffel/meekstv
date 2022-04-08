package events

type QuotaAdjusted struct {
	Current  int64
	Previous int64
	Scale    int64
}

func (e *QuotaAdjusted) Process() string {
	if e.Previous == 0 {
		return "Election quota set to " + formatScaledValue(e.Current, e.Scale)
	}

	return formatDiff(e.Current, e.Previous, e.Scale, "Election", "quota")
}
