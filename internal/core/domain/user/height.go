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

type Height struct {
	Value float64
}

type HeightUnit string

const (
	Cm HeightUnit = "cm"
	Ft HeightUnit = "ft"
)

func NewHeight(value float64, unit HeightUnit) (Height, error) {
	if value < 0 {
		return Height{}, ErrNegativeHeight
	}
	if value == 0 {
		return Height{}, ErrHeightZero
	}

	if unit == Ft {
		value *= 30.48
	}

	return Height{
		Value: value,
	}, nil
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

func (h Height) Display(unit HeightUnit) float64 {
	if unit == Ft {
		return h.Value / 30.48
	}

	return h.Value
}
