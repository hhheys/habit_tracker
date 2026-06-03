package achievement

type Achievement struct {
	ID          uint
	Name        string
	Description string
	Conditions  []Condition
}

type Condition struct {
	MetricCode  uint
	TargetValue int
}
