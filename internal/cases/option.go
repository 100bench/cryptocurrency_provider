package cases

type Aggregation int

const (
	_ Aggregation = iota
	AggMin
	AggMax
	AggAvg
)

func (a Aggregation) String() string {
	return [...]string{"", "MIN", "MAX", "AVG"}[a]
}

type Options struct {
	agg Aggregation
}

type Option func(*Options)

func WithMin() Option {
	return func(o *Options) {
		o.agg = AggMin
	}
}
