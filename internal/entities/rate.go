package entities

import (
	"github.com/pkg/errors"
	"time"
)

type Rate struct {
	Currency string
	Price    float64 // относительно доаллара
	Ts       time.Time
}

func NewRate(currency string, price float64, currentTime time.Time) (*Rate, error) {
	if currency == "" {
		return nil, errors.Wrap(ErrInvalidParams, "NewRate: currency is empty")
	}
	if price <= 0 {
		return nil, errors.Wrap(ErrInvalidParams, "NewRate: price is zero")
	}
	return &Rate{currency, price, currentTime}, nil
}
