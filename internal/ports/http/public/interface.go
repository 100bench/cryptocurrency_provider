package public

import (
	"context"
	en "github.com/100bench/cryptocurrency_provider.git/internal/entities"
)

type PublicService interface {
	GetMinRate(ctx context.Context, currencies []string) ([]en.Rate, error)
	GetMaxRate(ctx context.Context, currencies []string) ([]en.Rate, error)
	GetAvgRate(ctx context.Context, currencies []string) ([]en.Rate, error)
	GetLast(ctx context.Context, currencies []string) ([]en.Rate, error)
}
