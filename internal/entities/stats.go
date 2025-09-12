package entities

type Stats struct {
	Min float64
	Max float64
	Avg float64
}

func NewStats(min float64, max float64, avg float64) (Stats, error) {
	return Stats{min, max, avg}, nil
}
