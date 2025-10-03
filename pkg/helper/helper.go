// Package helper
package helper

import (
	"github.com/cheezecakee/logr"
	"golang.org/x/crypto/bcrypt"
)

func Clamp(value, min, max int) int {
	switch {
	case value < min:
		return min
	case value > max:
		return max
	default:
		return value
	}
}

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		logr.Get().Errorf("failed to hash password: %v", err)
		return "", err
	}
	return string(hashedPassword), nil
}

func ComparePassword(hash, password string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)); err != nil {
		logr.Get().Errorf("failed to compare password: %v", err)
		return err
	}
	return nil
}
