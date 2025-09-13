package pgx

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"
	"time"
)

type PgxClient struct {
	pool *pgxpool.Pool
}

func NewPgxClient(ctx context.Context, dsn string) (*PgxClient, error) {
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

	return &PgxClient{pool: pool}, nil
}

func (c *PgxClient) Close() {
	c.pool.Close()
}
