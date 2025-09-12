package cases

import (
	"context"
	"fmt"
)

type ServiceAPI struct {
	prov Client
	pub  Broker // сует в kafka (так предпологается)
}

func NewService(prov Client, pub Broker, currency string) (*ServiceAPI, error) {
	if prov == nil {
		return nil, ErrNilProvider
	}
	if pub == nil {
		return nil, ErrNilPublisher
	}
	return &ServiceAPI{prov, pub, currency}, nil
}

func (s *ServiceAPI) ProduceToBroker(ctx context.Context, currencies []string) error {
	rates, err := s.prov.GetRates(ctx, currencies)
	if err != nil {
		return fmt.Errorf("usecase Service.prov.GetRates: get rates for currencies=%v: %w", currencies, err)
	}
	err = s.pub.Publish(ctx, rates)
	if err != nil {
		return fmt.Errorf("usecase Service.pub.Produce: produce rates for currencies=%v: %w", currencies, err)
	}
	return nil
}
