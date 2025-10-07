package user

import (
	"errors"
)

var ErrInvalidVisibility = errors.New("invalid visibility")

type Visibility string

const (
	Public  Visibility = "public"
	Private Visibility = "private"
)

func NewVisibility(visibility string) (Visibility, error) {
	switch visibility {
	case "public":
		return Public, nil
	case "private":
		return Private, nil
	default:
		return "", ErrInvalidVisibility
	}
}
