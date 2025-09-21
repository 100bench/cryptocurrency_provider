package cases

import (
	"context"

	en "github.com/100bench/cryptocurrency_provider.git/internal/entities"
	"github.com/pkg/errors"
)

type StorageService struct {
	store Storage
}

func NewStorageService(store Storage) (*StorageService, error) {
	if store == nil {
		return nil, errors.Wrap(en.ErrNilDependency, "storage")
	}
	return &StorageService{store}, nil
}

func (s *StorageService) GetMinRate(ctx context.Context, currencies []string) ([]en.Rate, error) {
	rates, err := s.store.Get(ctx, currencies, WithMin())
	if err != nil {
		return nil, errors.Wrap(err, "storage.Get")
	}
	return rates, nil

}

func (s *StorageService) GetMaxRate(ctx context.Context, currencies []string) ([]en.Rate, error) {
	rates, err := s.store.Get(ctx, currencies, WithMax())
	if err != nil {
		return nil, errors.Wrap(err, "storage.Get")
	}
	return rates, nil
}

func (s *StorageService) GetAvgRate(ctx context.Context, currencies []string) ([]en.Rate, error) {
	rates, err := s.store.Get(ctx, currencies, WithAvg())
	if err != nil {
		return nil, errors.Wrap(err, "storage.Get")
	}
	return rates, nil
}

func (s *StorageService) GetLast(ctx context.Context, currencies []string) ([]en.Rate, error) {
	rates, err := s.store.Get(ctx, currencies)
	if err != nil {
		return nil, errors.Wrap(err, "storage.Get")
	}
	return rates, nil
}
