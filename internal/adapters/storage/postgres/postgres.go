package postgres

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/100bench/cryptocurrency_provider.git/internal/cases"
	en "github.com/100bench/cryptocurrency_provider.git/internal/entities"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"
)

type PgxStorage struct {
	pool *pgxpool.Pool
}

func NewPgxClient(ctx context.Context, dsn string) (*PgxStorage, error) {
	pool, err := pgxpool.Connect(ctx, dsn)
	if err != nil {
		return nil, errors.Wrap(err, "pgxpool.ConnectConfig")
	}

	return &PgxStorage{pool: pool}, nil
}

func (c *PgxStorage) Close() {
	c.pool.Close()
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
		log.Printf("query: %v", q)

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
		if err != nil {
			return nil, errors.Wrap(err, "scanRates")
		}
		log.Printf("query: %v", q)
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
			return errors.Wrap(err, "pgx.exec")
		}
	}
	return nil
}

func (c *PgxStorage) GetSymbols(ctx context.Context) ([]string, error) {
	q := `SELECT code FROM symbols WHERE kind = 'crypto';`
	rows, err := c.pool.Query(ctx, q)
	if err != nil {
		return nil, errors.Wrap(err, "pgx.query")
	}
	defer rows.Close()

	var symbols []string
	for rows.Next() {
		var code string
		if err = rows.Scan(&code); err != nil {
			return nil, errors.Wrap(err, "rows.Scan")
		}
		symbols = append(symbols, code)
	}
	if err = rows.Err(); err != nil {
		return nil, errors.Wrap(err, "rows.Err")
	}
	return symbols, nil
}

func scanRates(rows pgx.Rows) ([]en.Rate, error) {
	rates := make([]en.Rate, 0, 8)
	for rows.Next() {
		var (
			base  string
			price float64
			ts    time.Time
		)

		// Проверяем количество колонок в результате
		fieldDescriptions := rows.FieldDescriptions()
		if len(fieldDescriptions) == 2 {
			// Для агрегатных функций (MIN, MAX, AVG) - только base_code и price
			if err := rows.Scan(&base, &price); err != nil {
				return nil, fmt.Errorf("scan: %w", err)
			}
			ts = time.Now() // Устанавливаем текущее время для агрегатных данных
		} else {
			// Для обычных запросов - base_code, price, ts
			if err := rows.Scan(&base, &price, &ts); err != nil {
				return nil, fmt.Errorf("scan: %w", err)
			}
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
