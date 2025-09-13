package cases

import (
	"context"
	"fmt"
	"github.com/100bench/cryptocurrency_provider.git/internal/cases/port"
	en "github.com/100bench/cryptocurrency_provider.git/internal/entities"
	"github.com/pkg/errors"
)

type Publisher struct {
	broker port.Broker
}

func NewPublisher(pub port.Broker) (*Publisher, error) {
	if pub == nil {
		return nil, errors.Wrap(en.ErrNilDependency, "broker")
	}
	return &Publisher{pub}, nil
}

func (p *Publisher) Publish(ctx context.Context, rates []en.Rate) error {
	err := p.broker.Publish(ctx, rates)
	if err != nil {
		return fmt.Errorf("usecase Publisher.broker.Publish: publish rates=%v: %w", rates, err)
	}
	return nil
}
