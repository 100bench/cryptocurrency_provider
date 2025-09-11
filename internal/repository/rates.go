package repository

import (
	"context"
	en "github.com/100bench/cryptocurrency_provider.git/internal/entities"
	"github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"
)

type Irate interface {
	Get(ctx context.Context, currencies []string) ([]en.Rate, error)
}
type rates struct {
	postgres *sqlx.DB
	redis    *redis.Client
}

func NewRepository(db *sqlx.DB, redis *redis.Client) Irate {
	return &rates{
		postgres: db,
		redis:    redis,
	}
}

func (r *rates) Get(ctx context.Context, currencies []string) ([]en.Rate, error) {
	// тут я делаю запрос в редис, если там нет или ошибка, то я иду в бд, кладу в редис и потом возращаю клиенту.
	return nil, nil
}
