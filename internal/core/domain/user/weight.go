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
	Value float64
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

	if unit == Lb {
		value *= 0.453592 // Lb to kg for storage
	}

	return Weight{
		Value: value,
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

func (w Weight) Display(unit WeightUnit) float64 {
	if unit == Lb {
		return w.Value * 2.20462
	}
	return w.Value
}

// As helper function for something in the future... maybe?
func (w Weight) As(unit WeightUnit) Weight {
	switch unit {
	case Lb:
		return Weight{Value: w.Value * 2.20462}
	case Kg:
		return Weight{Value: w.Value * 0.453592}
	default:
		return w
	}
}

// Todo
// Set a limit to the weight depending on the unit.
