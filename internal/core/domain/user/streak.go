package user

import (
	"errors"
	"time"
)

var ErrInvalidRestDays = errors.New("rest days must be between 1-6")

type Streak struct {
	restDays    int
	current     int
	longest     int
	lastWorkout time.Time
}

func NewStreak(restDays int) (Streak, error) {
	if restDays == 0 {
		restDays = 2
	}

	if restDays < 1 || restDays > 6 {
		return Streak{}, ErrInvalidRestDays
	}

	return Streak{
		restDays: restDays,
		current:  0,
		longest:  0,
	}, nil
}

func (s Streak) RecordWorkout(workoutDate time.Time) Streak {
	if s.lastWorkout.IsZero() {
		return Streak{
			restDays:    s.restDays,
			current:     1,
			longest:     1,
			lastWorkout: workoutDate,
		}
	}

	// Calculate days since last workout
	hoursSince := workoutDate.Sub(s.lastWorkout).Hours()
	daysSince := hoursSince / 24

	newCurrent := s.current
	if daysSince > float64(s.restDays) {
		newCurrent = 1 // Streak broken
	} else {
		newCurrent++ // Streak continues
	}

	newLongest := s.longest
	newLongest = max(newLongest, newCurrent)

	return Streak{
		restDays:    s.restDays,
		current:     newCurrent,
		longest:     newLongest,
		lastWorkout: workoutDate,
	}
}

func (s Streak) RestDays() int {
	return s.restDays
}

func (s Streak) Current() int {
	return s.current
}

func (s Streak) Longest() int {
	return s.longest
}

func (s Streak) LastWorkout() time.Time {
	return s.lastWorkout
}

func (s Streak) IsActive() bool {
	if s.lastWorkout.IsZero() {
		return false
	}
	daysSince := time.Since(s.lastWorkout).Hours() / 24
	return daysSince <= float64(s.restDays)
}

func (s Streak) DaysUntilExpiry() int {
	if s.lastWorkout.IsZero() {
		return 0
	}
	daysSince := time.Since(s.lastWorkout).Hours() / 24
	remaining := s.restDays - int(daysSince)
	if remaining < 0 {
		return 0
	}
	return remaining
}
