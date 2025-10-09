package user

import (
	"errors"
	"strings"
)

var ErrInvalidBillingPeriod = errors.New("invalid billing period")

type Period string

const (
	Monthly Period = "monthly"
	Yearly  Period = "yearly"
)

func NewBillingPeriod(period string) (Period, error) {
	period = strings.TrimSpace(period)
	period = strings.ToLower(period)

	switch period {
	case "monthly":
		return Monthly, nil
	case "yearly":
		return Yearly, nil
	default:
		return "", ErrInvalidBillingPeriod
	}
}
