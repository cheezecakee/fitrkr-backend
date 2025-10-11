package user_test

import (
	"testing"
	"time"

	"github.com/cheezecakee/fitrkr-athena/internal/core/domain/user"
)

func TestNewStreak(t *testing.T) {
	s := user.NewStreak()

	if s.RestDays != 2 {
		t.Errorf("expected RestDays to default to 2, got %v", s.RestDays)
	}
	if s.Current != 0 {
		t.Errorf("expected Current to be 0, got %v", s.Current)
	}
	if s.Longest != 0 {
		t.Errorf("expected Longest to be 0, got %v", s.Longest)
	}
	if !s.LastWorkout.IsZero() {
		t.Error("expected LastWorkout to be zero")
	}
}

func TestStreak_UpdateRestDays(t *testing.T) {
	tests := []struct {
		name     string
		restDays int
		wantErr  bool
	}{
		{
			name:     "valid rest days",
			restDays: 3,
			wantErr:  false,
		},
		{
			name:     "valid rest days - minimum",
			restDays: 1,
			wantErr:  false,
		},
		{
			name:     "valid rest days - maximum",
			restDays: 6,
			wantErr:  false,
		},
		{
			name:     "default when zero",
			restDays: 0,
			wantErr:  false,
		},
		{
			name:     "invalid - negative",
			restDays: -5,
			wantErr:  true,
		},
		{
			name:     "invalid - over maximum",
			restDays: 8,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := user.NewStreak()
			err := s.UpdateRestDays(tt.restDays)

			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateRestDays() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr {
				expected := tt.restDays
				if expected == 0 {
					expected = 2 // defaults to 2
				}
				if s.RestDays != expected {
					t.Errorf("expected RestDays to be %d, got %d", expected, s.RestDays)
				}
			}
		})
	}
}

func TestStreak_RecordWorkout(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name        string
		setup       func() *user.Streak
		workoutTime time.Time
		wantCurrent int
		wantLongest int
	}{
		{
			name: "first workout initializes streak",
			setup: func() *user.Streak {
				return &user.Streak{RestDays: 2, Current: 0, Longest: 0}
			},
			workoutTime: now,
			wantCurrent: 1,
			wantLongest: 1,
		},
		{
			name: "second workout within restDays continues streak",
			setup: func() *user.Streak {
				s := &user.Streak{RestDays: 2, Current: 0, Longest: 0}
				s.RecordWorkout(now.Add(-24 * time.Hour))
				return s
			},
			workoutTime: now,
			wantCurrent: 2,
			wantLongest: 2,
		},
		{
			name: "missed beyond restDays resets streak",
			setup: func() *user.Streak {
				s := &user.Streak{RestDays: 2, Current: 0, Longest: 0}
				s.RecordWorkout(now.Add(-72 * time.Hour))
				return s
			},
			workoutTime: now,
			wantCurrent: 1,
			wantLongest: 1,
		},
		{
			name: "streak resets but keeps longest record",
			setup: func() *user.Streak {
				s := &user.Streak{RestDays: 2, Current: 0, Longest: 0}
				s.RecordWorkout(time.Now().Add(-48 * time.Hour))
				s.RecordWorkout(time.Now())
				s.RecordWorkout(time.Now().Add(-72 * time.Hour))
				return s
			},
			workoutTime: time.Now(),
			wantCurrent: 1,
			wantLongest: 2,
		},
		{
			name: "workout exactly on restDays boundary continues streak",
			setup: func() *user.Streak {
				s := &user.Streak{RestDays: 3, Current: 0, Longest: 0}
				s.RecordWorkout(time.Now().Add(-72 * time.Hour))
				return s
			},
			workoutTime: time.Now(),
			wantCurrent: 2,
			wantLongest: 2,
		},
		{
			name: "workout just over restDays resets streak",
			setup: func() *user.Streak {
				s := &user.Streak{RestDays: 3, Current: 0, Longest: 0}
				s.RecordWorkout(time.Now().Add(-73 * time.Hour))
				return s
			},
			workoutTime: time.Now(),
			wantCurrent: 1,
			wantLongest: 1,
		},
		{
			name: "10 consecutive workouts within restDays",
			setup: func() *user.Streak {
				s := &user.Streak{RestDays: 2, Current: 0, Longest: 0}
				now := time.Now()
				for i := 9; i >= 0; i-- {
					s.RecordWorkout(now.Add(time.Duration(-i*24) * time.Hour))
				}
				return s
			},
			workoutTime: time.Now(),
			wantCurrent: 11,
			wantLongest: 11,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := tt.setup()
			s.RecordWorkout(tt.workoutTime)

			if s.Current != tt.wantCurrent {
				t.Errorf("Current = %v, want %v", s.Current, tt.wantCurrent)
			}
			if s.Longest != tt.wantLongest {
				t.Errorf("Longest = %v, want %v", s.Longest, tt.wantLongest)
			}
		})
	}
}

func TestStreak_IsActive(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name       string
		setup      func() *user.Streak
		wantActive bool
	}{
		{
			name: "no workouts yet",
			setup: func() *user.Streak {
				return &user.Streak{RestDays: 2, Current: 0, Longest: 0}
			},
			wantActive: false,
		},
		{
			name: "recent workout within restDays",
			setup: func() *user.Streak {
				s := &user.Streak{RestDays: 3, Current: 0, Longest: 0}
				s.RecordWorkout(now.Add(-48 * time.Hour))
				return s
			},
			wantActive: true,
		},
		{
			name: "workout too long ago",
			setup: func() *user.Streak {
				s := &user.Streak{RestDays: 2, Current: 0, Longest: 0}
				s.RecordWorkout(now.Add(-96 * time.Hour))
				return s
			},
			wantActive: false,
		},
		{
			name: "new streak with no workouts is inactive",
			setup: func() *user.Streak {
				return &user.Streak{RestDays: 3, Current: 0, Longest: 0}
			},
			wantActive: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := tt.setup()
			if got := s.IsActive(); got != tt.wantActive {
				t.Errorf("IsActive() = %v, want %v", got, tt.wantActive)
			}
		})
	}
}

func TestStreak_DaysUntilExpiry(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name    string
		setup   func() *user.Streak
		wantRem int
	}{
		{
			name: "no workouts yet",
			setup: func() *user.Streak {
				return &user.Streak{RestDays: 2, Current: 0, Longest: 0}
			},
			wantRem: 0,
		},
		{
			name: "recent workout - within restDays",
			setup: func() *user.Streak {
				s := &user.Streak{RestDays: 3, Current: 0, Longest: 0}
				s.RecordWorkout(now.Add(-24 * time.Hour))
				return s
			},
			wantRem: 2,
		},
		{
			name: "expired streak",
			setup: func() *user.Streak {
				s := &user.Streak{RestDays: 2, Current: 0, Longest: 0}
				s.RecordWorkout(now.Add(-96 * time.Hour))
				return s
			},
			wantRem: 0,
		},
		{
			name: "half day since last workout",
			setup: func() *user.Streak {
				s := &user.Streak{RestDays: 2, Current: 0, Longest: 0}
				s.RecordWorkout(time.Now().Add(-12 * time.Hour))
				return s
			},
			wantRem: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := tt.setup()
			if got := s.DaysUntilExpiry(); got != tt.wantRem {
				t.Errorf("DaysUntilExpiry() = %v, want %v", got, tt.wantRem)
			}
		})
	}
}

func TestStreak_Break(t *testing.T) {
	s := &user.Streak{RestDays: 2, Current: 0, Longest: 0}
	s.RecordWorkout(time.Now())

	if s.Current == 0 {
		t.Error("expected Current to be set before Break()")
	}

	s.Break()

	if s.Current != 0 {
		t.Errorf("expected Current = 0 after Break(), got %v", s.Current)
	}
	if !s.LastWorkout.IsZero() {
		t.Error("expected LastWorkout to be zero after Break()")
	}
}

func TestStreak_Progress(t *testing.T) {
	s := &user.Streak{RestDays: 3, Current: 0, Longest: 0}
	now := time.Now()
	s.RecordWorkout(now.Add(-36 * time.Hour))

	progress := s.Progress()
	if progress <= 0 || progress >= 1 {
		t.Errorf("Progress() = %v, want between 0 and 1", progress)
	}
}
