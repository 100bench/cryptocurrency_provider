package cron

import (
	"context"
	"log"

	en "github.com/100bench/cryptocurrency_provider.git/internal/entities"
	"github.com/pkg/errors"

	"github.com/robfig/cron/v3"
)

func StartEvery5m(
	ctx context.Context,
	currencies []string,
	getRates func(ctx context.Context, currencies []string) ([]en.Rate, error),
) (*cron.Cron, error) {
	c := cron.New(cron.WithChain(
		cron.Recover(cron.DefaultLogger),
	))

	_, err := c.AddFunc("@every 5m", func() {
		// исполняется в своей горутине по расписанию
		rates, err := getRates(ctx, currencies)
		if err != nil {
			log.Printf("cron: getRates failed: %v", errors.WithStack(err))
			return
		}
		log.Printf("cron: fetched %d rates", len(rates))
	})
	if err != nil {
		return nil, errors.Wrap(err, "cron: schedule add")
	}

	c.Start()
	return c, nil
}
