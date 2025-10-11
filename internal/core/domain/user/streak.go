package user

import (
	"errors"
	"time"
)

var ErrInvalidRestDays = errors.New("rest days must be between 1-6")

type Streak struct {
	RestDays    int
	Current     int
	Longest     int
	LastWorkout time.Time
}

func NewStreak() Streak {
	return Streak{
		RestDays: 2,
		Current:  0,
		Longest:  0,
	}
}

func (s *Streak) UpdateRestDays(restDays int) error {
	if restDays == 0 {
		restDays = 2
	}

	if restDays < 1 || restDays > 6 {
		return ErrInvalidRestDays
	}

	s.RestDays = restDays

	return nil
}

func (s *Streak) RecordWorkout(workoutDate time.Time) {
	if s.LastWorkout.IsZero() {
		s.Current = 1
		s.Longest = 1
		s.LastWorkout = workoutDate
		return
	}

	// Calculate days since last workout
	daysSince := workoutDate.Sub(s.LastWorkout).Hours() / 24
	if daysSince > float64(s.RestDays) {
		s.Current = 1
	} else {
		s.Current++
	}

	if s.Current > s.Longest {
		s.Longest = s.Current
	}

	s.LastWorkout = workoutDate
}

func (s Streak) IsActive() bool {
	if s.LastWorkout.IsZero() {
		return false
	}

	daysSince := time.Since(s.LastWorkout).Hours() / 24
	return daysSince <= float64(s.RestDays)
}

func (s Streak) DaysUntilExpiry() int {
	if s.LastWorkout.IsZero() {
		return 0
	}
	daysSince := time.Since(s.LastWorkout).Hours() / 24
	remaining := s.RestDays - int(daysSince)
	if remaining < 0 {
		return 0
	}
	return remaining
}

// Break manually reset the streak
func (s *Streak) Break() {
	s.Current = 0
	s.LastWorkout = time.Time{}
}

func (s Streak) Progress() float64 {
	if s.LastWorkout.IsZero() {
		return 0
	}

	daysSince := time.Since(s.LastWorkout).Hours() / 24
	if daysSince > float64(s.RestDays) {
		return 1
	}
	return daysSince / float64(s.RestDays)
}
