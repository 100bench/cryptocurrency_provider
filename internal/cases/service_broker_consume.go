package cases

import (
	"context"
	en "github.com/100bench/cryptocurrency_provider.git/internal/entities"
	"github.com/pkg/errors"
)

type Consumer struct {
	broker Broker
	store  Storage
}

func NewConsumer(broker Broker) (*Consumer, error) {
	if broker == nil {
		return nil, errors.Wrap(en.ErrNilDependency, "broker.Publish")
	}
	return &Consumer{broker: broker}, nil
}

func (c *Consumer) Comsume(ctx context.Context) error {
	ch, err := c.broker.Consume(ctx)
	if err != nil {
		return errors.Wrap(err, "usecase Consumer.broker.Consume")
	}
	err = c.store.Save(ctx, ch)
	if err != nil {
		return errors.Wrap(err, "usecase Consumer.store.Save")
	}
	return nil

}
