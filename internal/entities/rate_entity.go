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
	// code(string) check
	if len(currency) != 3 {
		return Rate{}, fmt.Errorf("invalid currency length: %q", currency)
	}
	upperCur := strings.ToUpper(currency)
	if !regexp.MustCompile(`^[A-Z]{3}$`).MatchString(upperCur) {
		return Rate{}, fmt.Errorf("invalid currency format: %q", currency)
	}
	// price check
	if price <= 0 {
		return Rate{}, fmt.Errorf("invalid price: %q", price)
	}
	// if time in future
	if currentTime.After(time.Now()) {
		return Rate{}, fmt.Errorf("invalid current time: %q", currentTime)
	}
	return Rate{upperCur, price, currentTime, 0}, nil
}
