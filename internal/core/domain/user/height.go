package user

import "errors"

var (
	ErrNegativeHeight    = errors.New("height cannot be negative")
	ErrInvalidHeightUnit = errors.New("invalid height unit")
)

type Height struct {
	value float64
	unit  HeightUnit
}

type HeightUnit string

const (
	Cm   HeightUnit = "cm"
	FtIn HeightUnit = "ft_in"
)

func NewHeight(value float64, unit HeightUnit) (Height, error) {
	if value < 0 {
		return Height{}, ErrNegativeHeight
	}
	if unit == "" {
		unit = Cm
	}
	if unit != Cm && unit != FtIn {
		return Height{}, ErrInvalidHeightUnit
	}
	return Height{
		value: value,
		unit:  unit,
	}, nil
}

func NewHeightFtIn(feet int, inches float64) (Height, error) {
	totalInches := float64(feet)*12 + inches
	return NewHeight(totalInches, FtIn)
}

func (h Height) ToCm() float64 {
	if h.unit == Cm {
		return h.value // already cm
	}

	return h.value * 2.54
}

func (h Height) ToFtIn() (feet int, inches float64) {
	totalInches := h.value
	if h.unit == Cm {
		totalInches = h.value / 2.54
	}

	feet = int(totalInches / 12)
	inches = totalInches - float64(feet*12)
	return feet, inches // ft_in
}

func (h Height) Unit() HeightUnit {
	return h.unit
}

func (h Height) Value() float64 {
	return h.value
}
