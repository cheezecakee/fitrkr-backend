package user

import (
	"errors"
	"strings"
)

type Plan string

var ErrInvalidPlan = errors.New("invalid plan")

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
