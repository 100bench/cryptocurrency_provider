package cases

import (
	"github.com/100bench/cryptocurrency_provider.git/internal/cases/port"
	en "github.com/100bench/cryptocurrency_provider.git/internal/entities"
	"github.com/pkg/errors"
)

type StorageService struct {
	store port.Storage
}

func NewStorageService(store port.Storage) (*StorageService, error) {
	if store == nil {
		return nil, errors.Wrap(en.ErrNilDependency, "storage")
	}
	return &StorageService{store}, nil
}
