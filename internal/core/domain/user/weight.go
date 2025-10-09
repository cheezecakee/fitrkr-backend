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

type Weight struct {
	value float64
	unit  WeightUnit
}

type WeightUnit string

const (
	Kg WeightUnit = "kg"
	Lb WeightUnit = "lb"
)

func NewWeight(value float64, unit WeightUnit) (Weight, error) {
	if value < 0 {
		return Weight{}, ErrNegativeWeight
	}
	if value == 0 {
		return Weight{}, ErrWeightZero
	}

	return Weight{
		value: value,
		unit:  unit,
	}, nil
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

func (w Weight) ToKg() float64 {
	if w.unit == Kg {
		return w.value
	}
	return w.value * 0.453592 // lb to kg
}

func (w Weight) ToLbs() float64 {
	if w.unit == Lb {
		return w.value
	}
	return w.value * 2.20462 // kg to lb
}

func (w Weight) Unit() WeightUnit {
	return w.unit
}

func (w Weight) Value() float64 {
	return w.value
}

// Todo
// Set a limit to the weight depending on the unit.
