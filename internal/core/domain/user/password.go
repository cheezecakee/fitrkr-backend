package user

import (
	"errors"

	"github.com/cheezecakee/fitrkr-backend/pkg/helper"
)

type Password string

var (
	ErrEmptyPassword    = errors.New("empty password")
	ErrPasswordTooShort = errors.New("password too short")
	ErrPasswordTooLong  = errors.New("password too long")
)

func NewPassword(password string) (Password, error) {
	if len(password) < 8 {
		return "", ErrPasswordTooShort
	}
	if len(password) > 16 {
		return "", ErrPasswordTooLong
	}

	hashPass, err := helper.HashPassword(password)
	if err != nil {
		return "", err
	}

	return Password(hashPass), nil
}
