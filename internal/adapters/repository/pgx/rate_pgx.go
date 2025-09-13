package pgx

import (
	"context"
	en "github.com/100bench/cryptocurrency_provider.git/internal/entities"
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

func (c *PgxClient) GetList(ctx context.Context, currencies []string) ([]en.Rate, error) {
	//Сформировать SQL-запрос с учётом фильтров и сортировки.
	//query := `SELECT r.base_code,
	//	   r.price,
	//	   r.ts
	//FROM rates r
	//JOIN (
	//	SELECT base_code, MAX(ts) AS max_ts
	//	FROM rates
	//	GROUP BY base_code
	//) last_rates
	//  ON r.base_code = last_rates.base_code
	// AND r.ts = last_rates.max_ts
	//WHERE r.base_code IN (SELECT code FROM symbols);`
	//Выполнить запрос через клиент pgx, передав ctx и параметры.
	//
	//Обработать результат: пройти по строкам результата (rows), в каждой строке вызвать Scan.
	//
	//Сконвертировать данные из базы в доменные сущности (Rate).
	//
	//Сложить сущности в слайс для возврата.
	//
	//Проверить ошибки при выполнении запроса и после итерации по строкам.
	//
	//Закрыть rows после завершения работы.
	return nil, nil
}
