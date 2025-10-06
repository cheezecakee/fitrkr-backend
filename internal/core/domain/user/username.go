package user

import (
	"errors"
	"regexp"
	"strings"
)

var (
	ErrEmptyUsername        = errors.New("empty username supplied")
	ErrUsernameTooShort     = errors.New("username too short")
	ErrUsernameTooLong      = errors.New("username too long")
	ErrUsernameInvalidChars = errors.New("username contains invalid characters")
)

type Username string

func NewUsername(username string) (Username, error) {
	username = strings.TrimSpace(username)
	username = strings.ToLower(username)

	if username == "" {
		return "", ErrEmptyUsername
	}

	if len(username) < 3 {
		return "", ErrUsernameTooShort
	}

	if len(username) > 20 {
		return "", ErrUsernameTooLong
	}
	// Only allow alphanumeric and underscore
	validUsername := regexp.MustCompile(`^[a-z0-9_]+$`)
	if !validUsername.MatchString(username) {
		return "", ErrUsernameInvalidChars
	}

	return Username(username), nil
}
