package user

import (
	"errors"
	"strings"
)

var ErrInvalidVisibility = errors.New("invalid visibility")

type Visibility string

const (
	Public  Visibility = "public"
	Private Visibility = "private"
)

func NewVisibility(visibility string) (Visibility, error) {
	visibility = strings.TrimSpace(visibility)
	visibility = strings.ToLower(visibility)

	if visibility == "" {
		return Public, nil
	}

	switch visibility {
	case "public":
		return Public, nil
	case "private":
		return Private, nil
	default:
		return "", ErrInvalidVisibility
	}
}
