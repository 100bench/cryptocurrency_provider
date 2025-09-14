package cases

import "errors"

type Aggregation uint8

const (
	AggMin Aggregation = iota
	AggMax
	AggAvg
)

var AllAggs = []Aggregation{AggMin, AggMax, AggAvg}

type GetConfig struct {
	aggs map[Aggregation]struct{}
}

type Option func(*GetConfig)

func NewGet(opts ...Option) GetConfig {
	c := GetConfig{aggs: make(map[Aggregation]struct{}, len(AllAggs))}
	for _, a := range AllAggs {
		c.aggs[a] = struct{}{}
	}
	for _, o := range opts {
		o(&c)
	}
	return c
}

func (c GetConfig) Validate() error {
	if len(c.aggs) == 0 {
		return errors.New("нужна хотя бы одна опция")
	}
	return nil
}

func (c GetConfig) Has(a Aggregation) bool {
	_, ok := c.aggs[a]
	return ok
}

func WithAggs(as ...Aggregation) Option {
	return func(c *GetConfig) {
		c.aggs = make(map[Aggregation]struct{}, len(as))
		for _, a := range as {
			c.aggs[a] = struct{}{}
		}
	}
}
