package user

import (
	"errors"
	"regexp"
	"strings"
)

var (
	ErrEmptyEmail   = errors.New("empty email supplied")
	ErrInvalidEmail = errors.New("invalid email")
)

type Email string

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

func NewEmail(email string) (Email, error) {
	email = strings.TrimSpace(email)
	if email == "" {
		return "", ErrEmptyEmail
	}

	if !emailRegex.MatchString(email) {
		return "", ErrInvalidEmail
	}

	return Email(email), nil
}
