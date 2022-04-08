package events

type WeightAdjusted struct {
	Scale    int64
	Name     string
	Previous int64
	Current  int64
}

func (e *WeightAdjusted) Process() string {
	return formatDiff(e.Current, e.Previous, e.Scale, e.Name, "weight")
}
