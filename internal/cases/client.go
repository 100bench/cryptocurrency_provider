package cases

import (
	"context"
	en "github.com/100bench/cryptocurrency_provider.git/internal/entities"
)

type Client interface {
	GetRatesFromClient(ctx context.Context, currencies []string) ([]en.Rate, error)
}
