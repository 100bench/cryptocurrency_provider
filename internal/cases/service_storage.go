package cases

import (
	"context"
	"fmt"
	"github.com/100bench/cryptocurrency_provider.git/internal/cases/port"
	en "github.com/100bench/cryptocurrency_provider.git/internal/entities"
	"github.com/pkg/errors"
)

type StoreToClient struct {
	store port.Storage
}

func NewStoreToClient(store port.Storage) (*StoreToClient, error) {
	if store == nil {
		return nil, errors.Wrap(en.ErrNilDependency, "storage")
	}
	return &StoreToClient{store}, nil
}

func (r *StoreToClient) GetCurrent(ctx context.Context, currencies []string) ([]en.Rate, error) {
	rates, err := r.store.GetList(ctx, currencies)
	if err != nil {
		return nil, fmt.Errorf("usecase StoreToClient.GetCurrent: get current for currencies=%v: %w", currencies, err)
	}

	return rates, nil
}

func (r *StoreToClient) GetListStats(ctx context.Context, currencies []string) (map[string]en.Stats, error) {
	stats, err := r.store.GetStats(ctx, currencies)
	if err != nil {
		return nil, fmt.Errorf("usecase StoreToClient.GetListStats: get stats for currencies=%v: %w", currencies, err)
	}
	return stats, nil
}
