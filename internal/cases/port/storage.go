package port

import (
	"context"
	en "github.com/100bench/cryptocurrency_provider.git/internal/entities"
)

type Storage interface {
	GetList(ctx context.Context, currencies []string) ([]en.Rate, error)

	// GetStats
	//GetMax24h(ctx context.Context, currencies []string) ([]en.Rate, error)
	//GetMin24h(ctx context.Context, currencies []string) ([]en.Rate, error)
	//GetAvg24h(ctx context.Context, currencies []string) ([]en.Rate, error)
	// Заменил на GetStats, получаем сразу три метода, через мапу смотрим то, что нам надо
	GetStats(ctx context.Context, currencies []string) (map[string]en.Stats, error)

	Save(ctx context.Context, rateChan <-chan en.Rate) error
}
