package cases

import (
	"context"
	en "github.com/100bench/cryptocurrency_provider.git/internal/entities"
)

type Broker interface {
	Produce(ctx context.Context, rates []en.Rate) error
	Consume(ctx context.Context) (<-chan en.Rate, error)
}
