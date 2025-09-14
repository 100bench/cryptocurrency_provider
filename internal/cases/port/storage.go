package port

import (
	"context"
	en "github.com/100bench/cryptocurrency_provider.git/internal/entities"
	"golang.org/x/text/cases"
)

type Storage interface {
	GetList(ctx context.Context) ([]string, error)                                        // получаем сам список валют
	Get(ctx context.Context, currencies []string, opt ...cases.Option) ([]en.Rate, error) // реализация метода
	Save(ctx context.Context, rateChan <-chan en.Rate) error                              // save from kafka
	// Store() реализовать метод, который будет сторить данные в список валют

}
