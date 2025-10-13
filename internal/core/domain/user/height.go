package user

import (
	"errors"
	"strings"
)

var (
	ErrNegativeHeight    = errors.New("height cannot be negative")
	ErrHeightZero        = errors.New("height cannot be zero")
	ErrInvalidHeightUnit = errors.New("invalid height unit")
)

type HeightValue float64

type HeightUnit string

const (
	Cm HeightUnit = "cm"
	Ft HeightUnit = "ft"
)

func NewHeight(value float64, unit HeightUnit) (HeightValue, error) {
	if value < 0 {
		return 0, ErrNegativeHeight
	}
	if value == 0 {
		return 0, ErrHeightZero
	}

	if unit == Ft {
		value *= 30.48
	}

	return HeightValue(value), nil
}

func NewHeightUnit(unit string) (HeightUnit, error) {
	if unit == "" {
		return Cm, nil // default to kg is empty
	}
	unit = strings.ToLower(unit)
	switch unit {
	case "cm":
		return Cm, nil
	case "ft":
		return Ft, nil
	default:
		return "", ErrInvalidHeightUnit
	}
}

func (h HeightValue) Display(unit HeightUnit) HeightValue {
	if unit == Ft {
		return h / 30.48
	}

	return h
}
