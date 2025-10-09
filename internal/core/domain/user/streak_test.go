package user_test

import (
	"testing"
	"time"

	"github.com/cheezecakee/fitrkr-backend/internal/core/domain/user"
)

func TestStreak(t *testing.T) {
	tests := []struct {
		name     string
		restDays int
		wantErr  bool
	}{
		{
			name:     "valid streak",
			restDays: 2,
			wantErr:  false,
		},
		{
			name:     "invalid streak - negative value",
			restDays: -5,
			wantErr:  true,
		},
		{
			name:     "ivalid streak - over 6",
			restDays: 8,
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := user.NewStreak(tt.restDays)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewStreak() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestStreak_ValidCases(t *testing.T) {
	tests := []struct {
		name     string
		restDays int
		want     int
	}{
		{
			name:     "valid streak",
			restDays: 3,
			want:     3,
		},
		{
			name:     "valid streak - default 2",
			restDays: 0,
			want:     2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := user.NewStreak(tt.restDays)
			restDays := got.RestDays()
			if err != nil {
				t.Errorf("NewStreak() unexpected error = %v", err)
			}
			if restDays != tt.want {
				t.Errorf("NewStreak() = %v, want %v", restDays, tt.want)
			}
		})
	}
}

func TestStreak_RecordWorkout(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name        string
		setup       func() user.Streak
		workoutTime time.Time
		wantCurrent int
		wantLongest int
	}{
		{
			name: "first workout initializes streak",
			setup: func() user.Streak {
				s, _ := user.NewStreak(2)
				return s
			},
			workoutTime: now,
			wantCurrent: 1,
			wantLongest: 1,
		},
		{
			name: "second workout within restDays continues streak",
			setup: func() user.Streak {
				s, _ := user.NewStreak(2)
				// record a workout 1 day ago
				return s.RecordWorkout(now.Add(-24 * time.Hour))
			},
			workoutTime: now,
			wantCurrent: 2,
			wantLongest: 2,
		},
		{
			name: "missed beyond restDays resets streak",
			setup: func() user.Streak {
				s, _ := user.NewStreak(2)
				// simulate older workouts
				s = s.RecordWorkout(now.Add(-72 * time.Hour)) // 3 days ago
				return s
			},
			workoutTime: now,
			wantCurrent: 1,
			wantLongest: 1, // since streak resets
		},
		{
			name: "streak resets but keeps longest record",
			setup: func() user.Streak {
				s, _ := user.NewStreak(2)
				s = s.RecordWorkout(time.Now().Add(-48 * time.Hour)) // 2 days ago
				s = s.RecordWorkout(time.Now())                      // continue
				s = s.RecordWorkout(time.Now().Add(-72 * time.Hour)) // break (3 days ago)
				return s
			},
			workoutTime: time.Now(),
			wantCurrent: 1,
			wantLongest: 2,
		},
		{
			name: "workout exactly on restDays boundary continues streak",
			setup: func() user.Streak {
				s, _ := user.NewStreak(3)
				s = s.RecordWorkout(time.Now().Add(-72 * time.Hour)) // exactly 3 days ago
				return s
			},
			workoutTime: time.Now(),
			wantCurrent: 2,
			wantLongest: 2,
		},
		{
			name: "workout just over restDays resets streak",
			setup: func() user.Streak {
				s, _ := user.NewStreak(3)
				s = s.RecordWorkout(time.Now().Add(-73 * time.Hour)) // slightly over 3 days
				return s
			},
			workoutTime: time.Now(),
			wantCurrent: 1,
			wantLongest: 1,
		},
		{
			name: "10 consecutive workouts within restDays",
			setup: func() user.Streak {
				s, _ := user.NewStreak(2)
				now := time.Now()
				for i := 9; i >= 0; i-- {
					s = s.RecordWorkout(now.Add(time.Duration(-i*24) * time.Hour))
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
			initial := tt.setup()
			got := initial.RecordWorkout(tt.workoutTime)

			if got.Current() != tt.wantCurrent {
				t.Errorf("Current() = %v, want %v", got.Current(), tt.wantCurrent)
			}
			if got.Longest() != tt.wantLongest {
				t.Errorf("Longest() = %v, want %v", got.Longest(), tt.wantLongest)
			}
		})
	}
}

func TestStreak_IsActive(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name       string
		setup      func() user.Streak
		wantActive bool
	}{
		{
			name: "no workouts yet",
			setup: func() user.Streak {
				s, _ := user.NewStreak(2)
				return s
			},
			wantActive: false,
		},
		{
			name: "recent workout within restDays",
			setup: func() user.Streak {
				s, _ := user.NewStreak(3)
				return s.RecordWorkout(now.Add(-48 * time.Hour)) // 2 days ago
			},
			wantActive: true,
		},
		{
			name: "workout too long ago",
			setup: func() user.Streak {
				s, _ := user.NewStreak(2)
				return s.RecordWorkout(now.Add(-96 * time.Hour)) // 4 days ago
			},
			wantActive: false,
		},
		{
			name: "new streak with no workouts is inactive",
			setup: func() user.Streak {
				s, _ := user.NewStreak(3)
				return s
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
		setup   func() user.Streak
		wantRem int
	}{
		{
			name: "no workouts yet",
			setup: func() user.Streak {
				s, _ := user.NewStreak(2)
				return s
			},
			wantRem: 0,
		},
		{
			name: "recent workout - within restDays",
			setup: func() user.Streak {
				s, _ := user.NewStreak(3)
				return s.RecordWorkout(now.Add(-24 * time.Hour)) // 1 day ago
			},
			wantRem: 2,
		},
		{
			name: "expired streak",
			setup: func() user.Streak {
				s, _ := user.NewStreak(2)
				return s.RecordWorkout(now.Add(-96 * time.Hour)) // 4 days ago
			},
			wantRem: 0,
		},
		{
			name: "half day since last workout",
			setup: func() user.Streak {
				s, _ := user.NewStreak(2)
				return s.RecordWorkout(time.Now().Add(-12 * time.Hour))
			},
			wantRem: 2, // should round down correctly
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
	s, _ := user.NewStreak(2)
	s = s.RecordWorkout(time.Now())
	s = s.Break()

	if s.Current() != 0 {
		t.Errorf("expected Current() = 0 after Break(), got %v", s.Current())
	}
	if !s.LastWorkout().IsZero() {
		t.Errorf("expected LastWorkout() to be zero after Break()")
	}
}

func TestStreak_Progress(t *testing.T) {
	s, _ := user.NewStreak(3)
	now := time.Now()
	s = s.RecordWorkout(now.Add(-36 * time.Hour)) // 1.5 days ago

	progress := s.Progress()
	if progress <= 0 || progress >= 1 {
		t.Errorf("Progress() = %v, want between 0 and 1", progress)
	}
}
