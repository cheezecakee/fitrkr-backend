package user_test

import (
	"errors"
	"testing"

	"github.com/cheezecakee/fitrkr-athena/internal/core/domain/user"
)

func TestNewDuration(t *testing.T) {
	tests := []struct {
		name        string
		minutes     int
		expectError error
	}{
		{
			name:        "valid duration - 30 minutes",
			minutes:     30,
			expectError: nil,
		},
		{
			name:        "valid duration - 1 minute",
			minutes:     1,
			expectError: nil,
		},
		{
			name:        "valid duration - large value",
			minutes:     1440, // 24 hours
			expectError: nil,
		},
		{
			name:        "zero duration",
			minutes:     0,
			expectError: user.ErrZeroDuration,
		},
		{
			name:        "negative duration",
			minutes:     -1,
			expectError: user.ErrNegativeDuration,
		},
		{
			name:        "large negative duration",
			minutes:     -100,
			expectError: user.ErrNegativeDuration,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			duration, err := user.NewDuration(tt.minutes)

			if tt.expectError != nil {
				if err == nil {
					t.Errorf("expected error %v, got nil", tt.expectError)
					return
				}
				if !errors.Is(err, tt.expectError) {
					t.Errorf("expected error %v, got %v", tt.expectError, err)
				}
				return
			}

			if err != nil {
				t.Errorf("expected no error, got %v", err)
				return
			}

			if duration.Minutes() != tt.minutes {
				t.Errorf("expected minutes %d, got %d", tt.minutes, duration.Minutes())
			}
		})
	}
}

func TestDuration_Minutes(t *testing.T) {
	tests := []struct {
		name            string
		minutes         int
		expectedMinutes int
	}{
		{
			name:            "30 minutes",
			minutes:         30,
			expectedMinutes: 30,
		},
		{
			name:            "1 minute",
			minutes:         1,
			expectedMinutes: 1,
		},
		{
			name:            "120 minutes",
			minutes:         120,
			expectedMinutes: 120,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			duration := mustNewDuration(tt.minutes)

			if duration.Minutes() != tt.expectedMinutes {
				t.Errorf("expected %d minutes, got %d", tt.expectedMinutes, duration.Minutes())
			}
		})
	}
}

// Helper function for tests
func mustNewDuration(minutes int) user.Duration {
	d, err := user.NewDuration(minutes)
	if err != nil {
		panic(err)
	}
	return d
}
