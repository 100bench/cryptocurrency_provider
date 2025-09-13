package cases

import (
	"context"
	"fmt"
	"github.com/100bench/cryptocurrency_provider.git/internal/cases/port"
	en "github.com/100bench/cryptocurrency_provider.git/internal/entities"
	"github.com/pkg/errors"
)

type ServiceAPI struct {
	client port.Client
}

func NewServiceAPI(prov port.Client) (*ServiceAPI, error) {
	if prov == nil {
		return nil, errors.Wrap(en.ErrNilDependency, "client")
	}
	return &ServiceAPI{prov}, nil
}

func (s *ServiceAPI) GetRates(ctx context.Context, currencies []string) ([]en.Rate, error) {
	rates, err := s.client.GetRatesFromClient(ctx, currencies)
	if err != nil {
		return nil, fmt.Errorf("usecase Service.prov.GetRates: get rates for currencies=%v: %w", currencies, err)
	}
	return rates, nil
}
