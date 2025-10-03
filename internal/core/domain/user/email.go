package user

import (
	"errors"
	"strings"
)

var (
	ErrEmptyEmail   = errors.New("empty email supplied")
	ErrInvalidEmail = errors.New("invalid email")
)

type Email string

func NewEmail(email string) (Email, error) {
	email = strings.TrimSpace(email)
	if email == "" {
		return "", ErrEmptyEmail
	}

	// TODO logic to validate email

	return Email(email), nil
}
