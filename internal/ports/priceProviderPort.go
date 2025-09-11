package ports

import (
	"context"
	en "github.com/100bench/cryptocurrency_provider.git/internal/entities"
)

type priceProvider interface {
	GetRates(ctx context.Context, currency string) (Rates []en.Rate, err error)
}
