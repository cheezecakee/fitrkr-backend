package user

import (
	"errors"
	"strings"
)

var (
	ErrEmptyUsername    = errors.New("empty username supplied")
	ErrUsernameTooShort = errors.New("username too short")
)

type Username string

func NewUsername(username string) (Username, error) {
	username = strings.TrimSpace(username)
	if username == "" {
		return "", ErrEmptyUsername
	}

	if len(username) < 3 {
		return "", ErrUsernameTooShort
	}

	return Username(username), nil
}
