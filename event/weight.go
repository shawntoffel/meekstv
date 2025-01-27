package event

type WeightAdjusted struct {
	Name     string
	Current  int64
	Previous int64
	Scale    int64
}

func (e WeightAdjusted) Describe() string {
	diff := e.Current - e.Previous
	change := "reduced"
	if diff > 0 {
		change = "increased"
	} else {
		diff *= -1
	}

	formattedDiff := formatScaledValue(diff, e.Scale)
	formattedCurrent := formatScaledValue(e.Current, e.Scale)

	return e.Name + " weight " + change + " by " + formattedDiff + ". New weight: " + formattedCurrent + "."
}
