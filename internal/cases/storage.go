package cases

import (
	"context"
	en "github.com/100bench/cryptocurrency_provider.git/internal/entities"
)

type Storage interface {
	GetList(ctx context.Context) ([]string, error)                                  // получаем сам список валют
	Get(ctx context.Context, currencies []string, opt ...Option) ([]en.Rate, error) // реализация метода
	Save(ctx context.Context, rateChan <-chan en.Rate) error                        // save from kafka
	Store(ctx context.Context, currencies []string) error                           // не уверен, что это правильно
}
