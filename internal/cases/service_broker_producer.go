package cases

import (
	"context"
	"fmt"
	en "github.com/100bench/cryptocurrency_provider.git/internal/entities"
	"github.com/pkg/errors"
)

type Producer struct {
	broker Broker
}

func NewProducer(prod Broker) (*Producer, error) {
	if prod == nil {
		return nil, errors.Wrap(en.ErrNilDependency, "broker")
	}
	return &Producer{prod}, nil
}

func (p *Producer) Produce(ctx context.Context, rates []en.Rate) error {
	err := p.broker.Produce(ctx, rates)
	if err != nil {
		return fmt.Errorf("usecase Producer.broker.Produce: produce rates=%v: %w", rates, err)
	}
	return nil
}
