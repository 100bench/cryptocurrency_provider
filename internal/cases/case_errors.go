package cases

import "errors"

var (
	ErrNilStorage   = errors.New("nil Storage dependency")
	ErrNilProvider  = errors.New("nil PriceProvider dependency")
	ErrNilPublisher = errors.New("nil RatesPublisher dependency")
)
