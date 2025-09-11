package entities

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

type Rate struct {
	Currency    string
	Price       float64 // относительно доаллара
	CurrentTime time.Time
	ID          int
}

func NewRate(currency string, price float64, currentTime time.Time) (Rate, error) {
	c := strings.TrimSpace(currency)
	// проверка длины (3 символа)
	if len(c) != 3 {
		return Rate{}, fmt.Errorf("invalid currency length: %q", c)
	}
	// приведение к верхнему регистру
	upperCur := strings.ToUpper(c)
	// проверка, что только латинские буквы
	if !regexp.MustCompile(`^[A-Z]{3}$`).MatchString(upperCur) {
		return Rate{}, fmt.Errorf("invalid currency format: %q", c)
	}
	if price <= 0 {
		return Rate{}, fmt.Errorf("invalid price: %q", price)
	}
	return Rate{upperCur, price, currentTime, 0}, nil
}
