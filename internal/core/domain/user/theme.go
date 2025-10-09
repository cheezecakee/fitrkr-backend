package user

import (
	"errors"
	"strings"
)

var ErrInvalidTheme = errors.New("invalid theme")

type Theme string

const (
	Dark   Theme = "dark"
	Light  Theme = "light"
	System Theme = "system"
)

func NewTheme(theme string) (Theme, error) {
	theme = strings.TrimSpace(theme)
	theme = strings.ToLower(theme)

	if theme == "" {
		return System, nil
	}

	switch theme {
	case "dark":
		return Dark, nil
	case "light":
		return Light, nil
	case "system":
		return System, nil
	default:
		return "", ErrInvalidTheme
	}
}
