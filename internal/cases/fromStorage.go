package cases

import (
	"context"
	en "github.com/100bench/cryptocurrency_provider.git/internal/entities"
)

type Storage interface {
	GetList(ctx context.Context, currencies []string) ([]en.Rate, error) // ?

	// GetStats
	//GetMax24h(ctx context.Context, currencies []string) ([]en.Rate, error)
	//GetMin24h(ctx context.Context, currencies []string) ([]en.Rate, error)
	//GetAvg24h(ctx context.Context, currencies []string) ([]en.Rate, error)
	// Заменил на GetStats, получаем сразу три метода, через мапу смотрим то, что нам надо
	GetStats(ctx context.Context, currencies []string) (map[string]en.Stats, error)
}

type StoreToClient struct {
	store Storage
}

func NewStoreToClient(store Storage) (*StoreToClient, error) {
	if store == nil {
		return nil, ErrNilDependency
	}
	return &StoreToClient{store}, nil
}

func (r *StoreToClient) GetCurrent(ctx context.Context, currencies []string) ([]en.Rate, error) {
	rates, err := r.store.GetList(ctx, currencies)
	if err != nil {
		return nil, err
	}

	return rates, nil
}

func (r *StoreToClient) GetListStats(ctx context.Context, currencies []string) (map[string]en.Stats, error) {
	return r.store.GetStats(ctx, currencies)
}
