package user

import "errors"

var (
	ErrNegativeDuration = errors.New("duration cannot be negative")
	ErrZeroDuration     = errors.New("duration cannot be zero")
)

type Duration struct {
	minutes int
}

func NewDuration(minutes int) (Duration, error) {
	if minutes < 0 {
		return Duration{}, ErrNegativeDuration
	}
	if minutes == 0 {
		return Duration{}, ErrZeroDuration
	}

	return Duration{minutes: minutes}, nil
}

func (d Duration) Minutes() int {
	return d.minutes
}
