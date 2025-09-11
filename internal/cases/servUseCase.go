package cases

import (
	"context"
	en "github.com/100bench/cryptocurrency_provider.git/internal/entities"
)

type priceProvider interface {
	GetRates(ctx context.Context, currency string) (Rates []en.Rate, err error)
}

type RatesPublisher interface {
	Publish(ctx context.Context, rates []en.Rate) error
}

type Service struct {
	prov     priceProvider
	pub      RatesPublisher // сует в kafka (так предпологается)
	currency string
}

func NewService(prov priceProvider, pub RatesPublisher, currency string) *Service {
	return &Service{prov, pub, currency}
}

func (s *Service) GetRates(ctx context.Context) error {
	rates, err := s.prov.GetRates(ctx, s.currency)
	if err != nil {
		return err
	}
	return s.pub.Publish(ctx, rates)
}
