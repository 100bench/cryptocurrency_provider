package pgx

import (
	"context"
	en "github.com/100bench/cryptocurrency_provider.git/internal/entities"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"
	"time"
)

type PgxStorage struct {
	pool *pgxpool.Pool
}

func NewPgxClient(ctx context.Context, dsn string) (*PgxStorage, error) {
	cfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, errors.Wrap(err, "pgxpool.ParseConfig")
	}

	// pool settings
	cfg.MaxConns = 10
	cfg.MinConns = 2
	cfg.MaxConnLifetime = time.Hour
	cfg.MaxConnIdleTime = time.Minute

	pool, err := pgxpool.ConnectConfig(ctx, cfg)
	if err != nil {
		return nil, errors.Wrap(err, "pgxpool.ConnectConfig")
	}

	return &PgxStorage{pool: pool}, nil
}

func (c *PgxStorage) Close() {
	c.pool.Close()
}

func (c *PgxStorage) GetList(ctx context.Context, currencies []string) ([]en.Rate, error) {

	return nil, nil
}
