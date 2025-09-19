package app

import (
	"context"
	"github.com/100bench/cryptocurrency_provider.git/internal/cases"
	"github.com/100bench/cryptocurrency_provider.git/scheduler/cron"
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

	c, err := cron.StartEvery5m(ctx, s.api.GetRates) 
	if err != nil {
		return err
	}
	defer c.Stop()

	return nil
}
