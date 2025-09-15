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
	Agg Aggregation
}

type Option func(*Options)

func WithMin() Option {
	return func(o *Options) {
		o.Agg = AggMin
	}
}

func WithMax() Option {
	return func(o *Options) {
		o.Agg = AggMax
	}
}

func WithAvg() Option {
	return func(o *Options) {
		o.Agg = AggAvg
	}
}
