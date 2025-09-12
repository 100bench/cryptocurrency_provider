package cases

import (
	"context"
	en "github.com/100bench/cryptocurrency_provider.git/internal/entities"
)

type PriceProvider interface {
	GetRates(ctx context.Context, currency string) (Rates []en.Rate, err error)
}

type RatesPublisher interface {
	Publish(ctx context.Context, rates []en.Rate) error
}

type Service struct {
	prov     PriceProvider
	pub      RatesPublisher // сует в kafka (так предпологается)
	currency string
}

func NewService(prov PriceProvider, pub RatesPublisher, currency string) (*Service, error) {
	if prov == nil {
		return nil, ErrNilDependency
	}
	if pub == nil {
		return nil, ErrNilDependency
	}
	return &Service{prov, pub, currency}, nil
}

func (s *Service) GetRates(ctx context.Context) error {
	rates, err := s.prov.GetRates(ctx, s.currency)
	if err != nil {
		return err
	}
	return s.pub.Publish(ctx, rates)
}
