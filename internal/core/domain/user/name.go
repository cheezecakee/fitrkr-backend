package user

import (
	"errors"
	"strings"
)

var (
	ErrEmptyName    = errors.New("empty name supplied")
	ErrNameTooShort = errors.New("name too short")
)

func NewName(firstName, lastName string) (string, error) {
	var name []string
	firstName = strings.TrimSpace(firstName)
	lastName = strings.TrimSpace(lastName)
	if firstName == "" || lastName == "" {
		return "", ErrEmptyName
	}

	if len(firstName) < 2 || len(lastName) < 2 {
		return "", ErrNameTooShort
	}

	name = append(name, firstName, lastName)

	fullname := strings.Join(name, " ")

	return fullname, nil
}
