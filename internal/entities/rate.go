package entities

import (
	"time"
)

type Rate struct {
	Currency    string
	Price       float64 // относительно доаллара
	CurrentTime time.Time
}

func NewRate(currency string, price float64, currentTime time.Time) (*Rate, error) {
	return &Rate{currency, price, currentTime}, nil
}
