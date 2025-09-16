package public

import (
	"context"

	"github.com/100bench/cryptocurrency_provider.git/pkg/dto"
)

type PublicService interface {
	GetRates(ctx context.Context, req dto.GetRatesRequest) ([]dto.RateItem, error) // по фильтру
	GetAvailableCurrencies(ctx context.Context) ([]string, error)                  // получить список доступных валют
	GetLatestRate(ctx context.Context, currency string) (*dto.RateItem, error)
	GetMaxRate(ctx context.Context, currency string) (*dto.RateItem, error)
	GetMinRate(ctx context.Context, currency string) (*dto.RateItem, error)
	GetAvgRate(ctx context.Context, currency string) (*dto.RateItem, error)
}
