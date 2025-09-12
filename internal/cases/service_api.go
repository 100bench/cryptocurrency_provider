package cases

import (
	"context"
	"fmt"
	en "github.com/100bench/cryptocurrency_provider.git/internal/entities"
	"github.com/pkg/errors"
)

type ServiceAPI struct {
	client Client
	broker Broker // сует в kafka (так предпологается)
}

func NewServiceAPI(prov Client, pub Broker) (*ServiceAPI, error) {
	if prov == nil {
		return nil, errors.Wrap(en.ErrNilDependency, "client")
	}
	if pub == nil {
		return nil, errors.Wrap(en.ErrNilDependency, "broker")
	}
	return &ServiceAPI{prov, pub}, nil
}

func (s *ServiceAPI) ProduceToBroker(ctx context.Context, currencies []string) error {
	rates, err := s.client.GetRates(ctx, currencies)
	if err != nil {
		return fmt.Errorf("usecase Service.prov.GetRates: get rates for currencies=%v: %w", currencies, err)
	}
	err = s.broker.Publish(ctx, rates)
	if err != nil {
		return fmt.Errorf("usecase Service.pub.Produce: produce rates for currencies=%v: %w", currencies, err)
	}
	return nil
}
