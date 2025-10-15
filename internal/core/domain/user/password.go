package user

import (
	"errors"
	"regexp"

	"github.com/cheezecakee/logr"
	"golang.org/x/crypto/bcrypt"
)

type Password string

var (
	ErrEmptyPassword     = errors.New("empty password")
	ErrPasswordTooShort  = errors.New("password too short")
	ErrPasswordTooLong   = errors.New("password too long")
	ErrPasswordNoChar    = errors.New("password must contain a character")
	ErrPasswordNoDigit   = errors.New("password must contain a digit")
	ErrPasswordNoSpecial = errors.New("password must contain a special character")
	ErrPasswordNoUpper   = errors.New("password must contain at least one uppercase letter")
)

var (
	hasUpperLetterRegex = regexp.MustCompile(`[A-Z]`)
	hasLetterRegex      = regexp.MustCompile(`[a-z]`)
	hasDigitRegex       = regexp.MustCompile(`[0-9]`)
	hasSpecialRegex     = regexp.MustCompile(`[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>/?]`)
)

func NewPassword(password string) (Password, error) {
	if password == "" {
		return "", ErrEmptyPassword
	}

	if len(password) < 8 {
		return "", ErrPasswordTooShort
	}

	if len(password) > 64 {
		return "", ErrPasswordTooLong
	}

	if !hasUpperLetterRegex.MatchString(password) {
		return "", ErrPasswordNoUpper
	}

	if !hasLetterRegex.MatchString(password) {
		return "", ErrPasswordNoChar
	}

	if !hasDigitRegex.MatchString(password) {
		return "", ErrPasswordNoDigit
	}

	if !hasSpecialRegex.MatchString(password) {
		return "", ErrPasswordNoSpecial
	}

	hashPass, err := HashPassword(password)
	if err != nil {
		return "", err
	}

	return Password(hashPass), nil
}

func (p Password) Verify(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(p), []byte(password))
	return err == nil
}

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		logr.Get().Errorf("failed to hash password: %v", err)
		return "", err
	}
	return string(hashedPassword), nil
}
