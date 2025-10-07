package user

import "errors"

var (
	ErrNegativeWeight    = errors.New("weight cannot be negative")
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
	if unit == "" {
		unit = Kg
	}
	if unit != Kg && unit != Lb {
		return Weight{}, ErrInvalidWeightUnit
	}
	return Weight{
		value: value,
		unit:  unit,
	}, nil
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
