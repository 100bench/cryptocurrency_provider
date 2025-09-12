package cases

import (
	"context"
	en "github.com/100bench/cryptocurrency_provider.git/internal/entities"
)

type Broker interface {
	Publish(ctx context.Context, rates []en.Rate) error
}
