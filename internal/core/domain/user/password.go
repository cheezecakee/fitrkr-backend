package user

import (
	"errors"
	"regexp"

	"github.com/cheezecakee/fitrkr-backend/pkg/helper"
)

type Password string

var (
	ErrEmptyPassword     = errors.New("empty password")
	ErrPasswordTooShort  = errors.New("password too short")
	ErrPasswordTooLong   = errors.New("password too long")
	ErrPasswordNoChar    = errors.New("password must contain a character")
	ErrPasswordNoDigit   = errors.New("password must contain a digit")
	ErrPasswordNoSpecial = errors.New("password must contain a special character")
)

var (
	hasLetterRegex = regexp.MustCompile(`[a-zA-Z]`)
	hasDigitRegex  = regexp.MustCompile(`[0-9]`)
	specialRegex   = regexp.MustCompile(`[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>/?]`)
)

func NewPassword(password string) (Password, error) {
	if password == "" {
		return "", ErrEmptyPassword
	}

	if len(password) < 8 {
		return "", ErrPasswordTooShort
	}

	if len(password) > 16 {
		return "", ErrPasswordTooLong
	}

	if !hasLetterRegex.MatchString(password) {
		return "", ErrPasswordNoChar
	}

	if !hasDigitRegex.MatchString(password) {
		return "", ErrPasswordNoDigit
	}

	if !specialRegex.MatchString(password) {
		return "", ErrPasswordNoSpecial
	}

	hashPass, err := helper.HashPassword(password)
	if err != nil {
		return "", err
	}

	return Password(hashPass), nil
}
