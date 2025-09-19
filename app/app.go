package app

import (
	"context"
	"github.com/100bench/cryptocurrency_provider.git/internal/cases"
)

type Service struct {
	api       cases.ServiceAPI
	consumer  cases.Consumer
	publisher cases.Publisher
	storage   cases.Storage
}

func NewApp(api cases.ServiceAPI, consume cases.Consumer, publish cases.Publisher, storage cases.Storage) *Service {
	return &Service{
		api:       api,
		consumer:  consume,
		publisher: publish,
		storage:   storage,
	}
}

func (s *Service) Run() error {
	ctx := context.Background()
	raets, err := s.api.GetRates(ctx)
}
