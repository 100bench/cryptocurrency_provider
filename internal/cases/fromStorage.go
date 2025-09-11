package cases

import (
	"context"
	en "github.com/100bench/cryptocurrency_provider.git/internal/entities"
)

type Storage interface {
	Get(ctx context.Context, currencies []string) ([]en.Rate, error)
}

type StoreToClient struct {
	store Storage
}

func NewStoreToClient(store Storage) *StoreToClient {
	return &StoreToClient{store}
}

func (r *StoreToClient) Run(ctx context.Context, currencies []string) ([]en.Rate, error) {
	rates, err := r.store.Get(ctx, currencies)
	if err != nil {
		return nil, err
	}

	return rates, nil
}
