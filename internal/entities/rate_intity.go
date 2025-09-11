package entities

import (
	"time"
)

type Rate struct {
	Currency    string    `json:"currency"`
	Price       int64     `json:"price"` // относительно доаллара
	CurrentTime time.Time `json:"current_time"`
	Id          int       `json:"id"`
}
