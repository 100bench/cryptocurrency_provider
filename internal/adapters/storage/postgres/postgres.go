package postgres

import (
	"context"
	"fmt"
	"github.com/100bench/cryptocurrency_provider.git/internal/cases"
	en "github.com/100bench/cryptocurrency_provider.git/internal/entities"
	"github.com/jackc/pgx/v4"
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

func (c *PgxStorage) GetList(ctx context.Context) ([]string, error) {
	const q = `
		SELECT DISTINCT code FROM symbols
	`
	rows, err := c.pool.Query(ctx, q)
	if err != nil {
		return nil, errors.Wrap(err, "pgx.query")
	}
	defer rows.Close()

	var currencies []string
	for rows.Next() {
		var code string
		if err := rows.Scan(&code); err != nil {
			return nil, errors.Wrap(err, "rows.Scan")
		}
		currencies = append(currencies, code)
	}
	if err := rows.Err(); err != nil {
		return nil, errors.Wrap(err, "rows.Err")
	}
	return currencies, nil
}

func (c *PgxStorage) Get(ctx context.Context, currencies []string, opt ...cases.Option) ([]en.Rate, error) {
	var cfg cases.Options
	for _, o := range opt {
		if o != nil {
			o(&cfg)
		}
	}

	var rates []en.Rate
	switch cfg.Agg {
	case cases.AggMin, cases.AggMax, cases.AggAvg:
		fn := cfg.Agg.String() // "MIN" | "MAX" | "AVG"
		q := fmt.Sprintf(`
            SELECT base_code, %s(price) AS price
            FROM rates
            WHERE base_code = ANY($1)
            GROUP BY base_code
        `, fn)

		rows, err := c.pool.Query(ctx, q, currencies)
		if err != nil {
			return nil, errors.Wrap(err, "pgx.query")
		}
		defer rows.Close()
		rates, err = scanRates(rows)
		if err != nil {
			return nil, errors.Wrap(err, "scanRates")
		}

	default:
		// последняя котировка
		const q = `
            SELECT DISTINCT ON (base_code) base_code, price, ts
            FROM rates
            WHERE base_code = ANY($1)
            ORDER BY base_code, ts DESC
        `
		rows, err := c.pool.Query(ctx, q, currencies)
		if err != nil {
			return nil, errors.Wrap(err, "pgx.query")
		}
		defer rows.Close()
		rates, err = scanRates(rows)
	}
	return rates, nil
}

func (c *PgxStorage) Save(ctx context.Context, rateChan <-chan en.Rate) error {
	for r := range rateChan {
		const q = `
			INSERT INTO rates (base_code, price, ts)
			VALUES ($1, $2, $3)
			ON CONFLICT DO NOTHING;
			`
		_, err := c.pool.Exec(ctx, q, r.Currency, r.Price, r.Ts)
		if err != nil {
		}
		return errors.Wrap(err, "pgx.exec")
	}
	return nil
}

func (c *PgxStorage) Store(ctx context.Context, currencies []string) error {
	const q = `INSERT INTO symbols (code) VALUES ($1);`
	_, err := c.pool.Exec(ctx, q, currencies)
	if err != nil {
		return errors.Wrap(err, "pgx.exec")
	}
	return nil
}

func scanRates(rows pgx.Rows) ([]en.Rate, error) {
	rates := make([]en.Rate, 0, 8)
	for rows.Next() {
		var (
			base  string
			price float64
			ts    time.Time
		)
		if err := rows.Scan(&base, &price, &ts); err != nil {
			return nil, fmt.Errorf("scan: %w", err)
		}
		rates = append(rates, en.Rate{
			Currency: base,
			Price:    price,
			Ts:       ts,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, errors.Wrap(err, "rows.Err")
	}
	return rates, nil
}