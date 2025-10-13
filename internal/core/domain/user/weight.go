package user

import (
	"errors"
	"strings"
)

var (
	ErrNegativeWeight    = errors.New("weight cannot be negative")
	ErrWeightZero        = errors.New("weight cannot be zero")
	ErrInvalidWeightUnit = errors.New("invalid weight unit")
)

type WeightValue float64

type WeightUnit string

const (
	Kg WeightUnit = "kg"
	Lb WeightUnit = "lb"
)

func NewWeight(value float64, unit WeightUnit) (WeightValue, error) {
	if value < 0 {
		return 0, ErrNegativeWeight
	}
	if value == 0 {
		return 0, ErrWeightZero
	}

	if unit == Lb {
		value *= 0.453592 // Lb to kg for storage
	}

	return WeightValue(value), nil
}

func NewWeightUnit(unit string) (WeightUnit, error) {
	if unit == "" {
		return Kg, nil // default to kg is empty
	}
	unit = strings.ToLower(unit)
	switch unit {
	case "kg":
		return Kg, nil
	case "lb":
		return Lb, nil
	default:
		return "", ErrInvalidWeightUnit
	}
}

func (w WeightValue) Display(unit WeightUnit) WeightValue {
	if unit == Lb {
		return w * 2.20462
	}
	return w
}

// As helper function for something in the future... maybe?
func (w WeightValue) As(unit WeightUnit) WeightValue {
	switch unit {
	case Lb:
		return w * 2.20462
	case Kg:
		return w * 0.453592
	default:
		return w
	}
}

// Todo
// Set a limit to the weight depending on the unit.
