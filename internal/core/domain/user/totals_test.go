package user_test

import (
	"testing"

	"github.com/cheezecakee/fitrkr-athena/internal/core/domain/user"
)

func TestNewTotals(t *testing.T) {
	totals := user.NewTotals()

	if totals.Workouts != 0 {
		t.Errorf("expected workouts to be 0, got %d", totals.Workouts)
	}
	if totals.Lifted != 0 {
		t.Errorf("expected lifted to be 0, got %f", totals.Lifted)
	}
	if totals.Time != 0 {
		t.Errorf("expected time to be 0, got %d", totals.Time)
	}
}

func TestTotals_RecordWorkout(t *testing.T) {
	tests := []struct {
		name             string
		lifted           user.Weight
		duration         user.Duration
		expectedWorkouts int
		expectedLifted   float64
		expectedTime     int
	}{
		{
			name:             "record complete workout",
			lifted:           mustNewWeightKg(150),
			duration:         mustNewDuration(45),
			expectedWorkouts: 1,
			expectedLifted:   150,
			expectedTime:     45,
		},
		{
			name:             "record workout with lb weight",
			lifted:           mustNewWeight(100, user.Lb),
			duration:         mustNewDuration(30),
			expectedWorkouts: 1,
			expectedLifted:   45.3592, // 100 * 0.453592
			expectedTime:     30,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			totals := user.NewTotals()
			totals.RecordWorkout(tt.lifted, tt.duration)

			if totals.Workouts != tt.expectedWorkouts {
				t.Errorf("expected workouts %d, got %d", tt.expectedWorkouts, totals.Workouts)
			}
			if totals.Lifted != tt.expectedLifted {
				t.Errorf("expected lifted %f, got %f", tt.expectedLifted, totals.Lifted)
			}
			if totals.Time != tt.expectedTime {
				t.Errorf("expected time %d, got %d", tt.expectedTime, totals.Time)
			}
		})
	}
}

func TestTotals_RecordWorkout_Multiple(t *testing.T) {
	totals := user.NewTotals()

	totals.RecordWorkout(mustNewWeightKg(100), mustNewDuration(30))
	totals.RecordWorkout(mustNewWeightKg(150), mustNewDuration(45))
	totals.RecordWorkout(mustNewWeightKg(75), mustNewDuration(20))

	if totals.Workouts != 3 {
		t.Errorf("expected 3 workouts, got %d", totals.Workouts)
	}
	if totals.Lifted != 325 {
		t.Errorf("expected lifted 325, got %f", totals.Lifted)
	}
	if totals.Time != 95 {
		t.Errorf("expected time 95, got %d", totals.Time)
	}
}

// Helper functions for tests
func mustNewWeight(value float64, unit user.WeightUnit) user.Weight {
	w, err := user.NewWeight(value, unit)
	if err != nil {
		panic(err)
	}
	return w
}

func mustNewWeightKg(kg float64) user.Weight {
	return mustNewWeight(kg, user.Kg)
}
