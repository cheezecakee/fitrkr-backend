package user

import (
	"errors"
	"strings"
)

var ErrInvalidPlan = errors.New("invalid plan")

type Plan string

const (
	Basic   Plan = "basic"
	Premium Plan = "premium"
)

func NewPlan(plan string) (Plan, error) {
	plan = strings.TrimSpace(plan)
	plan = strings.ToLower(plan)
	if plan == "" {
		return Basic, nil
	}

	switch plan {
	case "basic": // free tier
		return Basic, nil
	case "premium":
		return Premium, nil
	default:
		return "", ErrInvalidPlan
	}
}
