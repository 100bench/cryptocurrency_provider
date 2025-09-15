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
	s.store.Get(ctx, currencies, WithMin())
	return nil, nil

}
