package user

import (
	"errors"
)

var ErrInvalidBFP = errors.New("invalid bodyfat percentage")

type BFP float64

func NewBFP(value float64) (BFP, error) {
	if value < 0 || value > 100 {
		return 0, ErrInvalidBFP
	}
	return BFP(value), nil
}
