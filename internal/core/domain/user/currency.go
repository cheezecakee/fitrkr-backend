package user

import (
	"errors"
	"strings"
)

var ErrInvalidCurrency = errors.New("invalid currency")

type Currency string

const (
	USD Currency = "USD"
	EUR Currency = "EUR"
	GBP Currency = "GBP"
	CAD Currency = "CAD"
	AUD Currency = "AUD"
	BRL Currency = "BRL"
	RMB Currency = "RMB"
)

func NewCurrency(code string) (Currency, error) {
	code = strings.TrimSpace(code)

	if code == "" {
		return USD, nil
	}

	if len(code) != 3 {
		return "", ErrInvalidCurrency
	}

	code = strings.ToUpper(code)

	switch code {
	case "USD", "EUR", "GBP", "CAD", "AUD", "BRL", "RMB":
		return Currency(code), nil
	default:
		return "", ErrInvalidCurrency
	}
}
