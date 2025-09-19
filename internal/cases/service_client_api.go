package cases

import (
	"context"
	"fmt"

	en "github.com/100bench/cryptocurrency_provider.git/internal/entities"
	"github.com/pkg/errors"
)

type ServiceAPI struct {
	client  Client
	storage Storage
}

func NewServiceAPI(prov Client, storage Storage) (*ServiceAPI, error) {
	if prov == nil {
		return nil, errors.Wrap(en.ErrNilDependency, "client")
	}
	return &ServiceAPI{prov, storage}, nil
}

func (s *ServiceAPI) GetRates(ctx context.Context) ([]en.Rate, error) {
	currencies, err := s.storage.GetSymbols(ctx) // []string, error
	if err != nil {
		return nil, errors.Wrap(err, "usecase Service.prov.GetRates: get symbols from storage")
	}
	rates, err := s.client.GetRatesFromClient(ctx, currencies)
	if err != nil {
		return nil, fmt.Errorf("usecase Service.prov.GetRates: get rates for currencies=%v: %w", currencies, err)
	}
	return rates, nil
}
