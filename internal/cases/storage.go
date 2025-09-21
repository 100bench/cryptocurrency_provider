package cases

import (
	"context"

	en "github.com/100bench/cryptocurrency_provider.git/internal/entities"
)

type Storage interface {
	Get(ctx context.Context, currencies []string, opt ...Option) ([]en.Rate, error) // реализация метода
	Save(ctx context.Context, rateChan <-chan en.Rate) error                        // save from kafka
	GetSymbols(ctx context.Context) ([]string, error)
}
