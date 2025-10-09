package user_test

import (
	"testing"

	"github.com/cheezecakee/fitrkr-athena/internal/core/domain/user"
)

func TestNewTotals(t *testing.T) {
	totals := user.NewTotals()

	if totals.Workouts() != 0 {
		t.Errorf("expected workouts to be 0, got %d", totals.Workouts())
	}
	if totals.Volume() != 0 {
		t.Errorf("expected volume to be 0, got %f", totals.Volume())
	}
	if totals.TimeMinutes() != 0 {
		t.Errorf("expected time to be 0, got %d", totals.TimeMinutes())
	}
}

func TestAddWorkout(t *testing.T) {
	totals := user.NewTotals()

	updated := totals.AddWorkout()

	if updated.Workouts() != 1 {
		t.Errorf("expected workouts to be 1, got %d", updated.Workouts())
	}
	if updated.Volume() != 0 {
		t.Errorf("expected volume to be unchanged at 0, got %f", updated.Volume())
	}
	if updated.TimeMinutes() != 0 {
		t.Errorf("expected time to be unchanged at 0, got %d", updated.TimeMinutes())
	}
}

func TestAddWorkout_Multiple(t *testing.T) {
	totals := user.NewTotals()

	totals = totals.AddWorkout().AddWorkout().AddWorkout()

	if totals.Workouts() != 3 {
		t.Errorf("expected workouts to be 3, got %d", totals.Workouts())
	}
}

func TestAddWorkout_Immutability(t *testing.T) {
	original := user.NewTotals()
	updated := original.AddWorkout()

	if original.Workouts() != 0 {
		t.Errorf("original should be unchanged, expected 0 workouts, got %d", original.Workouts())
	}
	if updated.Workouts() != 1 {
		t.Errorf("updated should have 1 workout, got %d", updated.Workouts())
	}
}

func TestAddVolume(t *testing.T) {
	tests := []struct {
		name           string
		weight         user.Weight
		expectedVolume float64
	}{
		{
			name:           "add 100kg",
			weight:         mustNewWeightKg(100),
			expectedVolume: 100,
		},
		{
			name:           "add 50.5kg",
			weight:         mustNewWeightKg(50.5),
			expectedVolume: 50.5,
		},
		{
			name:           "add 220lb should convert to kg",
			weight:         mustNewWeight(220, user.Lb),
			expectedVolume: 99.790240, // 220 * 0.453592
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			totals := user.NewTotals()
			updated := totals.AddVolume(tt.weight)

			if updated.Volume() != tt.expectedVolume {
				t.Errorf("expected volume %f, got %f", tt.expectedVolume, updated.Volume())
			}
			if updated.Workouts() != 0 {
				t.Errorf("expected workouts to be unchanged at 0, got %d", updated.Workouts())
			}
			if updated.TimeMinutes() != 0 {
				t.Errorf("expected time to be unchanged at 0, got %d", updated.TimeMinutes())
			}
		})
	}
}

func TestAddVolume_Accumulation(t *testing.T) {
	totals := user.NewTotals()

	totals = totals.AddVolume(mustNewWeightKg(100))
	totals = totals.AddVolume(mustNewWeightKg(50))
	totals = totals.AddVolume(mustNewWeightKg(25.5))

	expected := 175.5
	if totals.Volume() != expected {
		t.Errorf("expected accumulated volume %f, got %f", expected, totals.Volume())
	}
}

func TestAddVolume_Immutability(t *testing.T) {
	original := user.NewTotals()
	weight := mustNewWeightKg(100)
	updated := original.AddVolume(weight)

	if original.Volume() != 0 {
		t.Errorf("original should be unchanged, expected 0 volume, got %f", original.Volume())
	}
	if updated.Volume() != 100 {
		t.Errorf("updated should have 100 volume, got %f", updated.Volume())
	}
}

func TestAddDuration(t *testing.T) {
	tests := []struct {
		name         string
		duration     user.Duration
		expectedTime int
	}{
		{
			name:         "add 60 minutes",
			duration:     mustNewDuration(60),
			expectedTime: 60,
		},
		{
			name:         "add 1 minute",
			duration:     mustNewDuration(1),
			expectedTime: 1,
		},
		{
			name:         "add 120 minutes",
			duration:     mustNewDuration(120),
			expectedTime: 120,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			totals := user.NewTotals()
			updated := totals.AddDuration(tt.duration)

			if updated.TimeMinutes() != tt.expectedTime {
				t.Errorf("expected time %d, got %d", tt.expectedTime, updated.TimeMinutes())
			}
			if updated.Workouts() != 0 {
				t.Errorf("expected workouts to be unchanged at 0, got %d", updated.Workouts())
			}
			if updated.Volume() != 0 {
				t.Errorf("expected volume to be unchanged at 0, got %f", updated.Volume())
			}
		})
	}
}

func TestAddDuration_Accumulation(t *testing.T) {
	totals := user.NewTotals()

	totals = totals.AddDuration(mustNewDuration(30))
	totals = totals.AddDuration(mustNewDuration(45))
	totals = totals.AddDuration(mustNewDuration(15))

	expected := 90
	if totals.TimeMinutes() != expected {
		t.Errorf("expected accumulated time %d, got %d", expected, totals.TimeMinutes())
	}
}

func TestAddDuration_Immutability(t *testing.T) {
	original := user.NewTotals()
	updated := original.AddDuration(mustNewDuration(60))

	if original.TimeMinutes() != 0 {
		t.Errorf("original should be unchanged, expected 0 time, got %d", original.TimeMinutes())
	}
	if updated.TimeMinutes() != 60 {
		t.Errorf("updated should have 60 time, got %d", updated.TimeMinutes())
	}
}

func TestRecordWorkout(t *testing.T) {
	tests := []struct {
		name             string
		volume           user.Weight
		duration         user.Duration
		expectedWorkouts int
		expectedVolume   float64
		expectedTime     int
	}{
		{
			name:             "record complete workout",
			volume:           mustNewWeightKg(150),
			duration:         mustNewDuration(45),
			expectedWorkouts: 1,
			expectedVolume:   150,
			expectedTime:     45,
		},
		{
			name:             "record workout with lb weight",
			volume:           mustNewWeight(100, user.Lb),
			duration:         mustNewDuration(30),
			expectedWorkouts: 1,
			expectedVolume:   45.3592, // 100 * 0.453592
			expectedTime:     30,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			totals := user.NewTotals()
			updated := totals.RecordWorkout(tt.volume, tt.duration)

			if updated.Workouts() != tt.expectedWorkouts {
				t.Errorf("expected workouts %d, got %d", tt.expectedWorkouts, updated.Workouts())
			}
			if updated.Volume() != tt.expectedVolume {
				t.Errorf("expected volume %f, got %f", tt.expectedVolume, updated.Volume())
			}
			if updated.TimeMinutes() != tt.expectedTime {
				t.Errorf("expected time %d, got %d", tt.expectedTime, updated.TimeMinutes())
			}
		})
	}
}

func TestRecordWorkout_Multiple(t *testing.T) {
	totals := user.NewTotals()

	totals = totals.RecordWorkout(mustNewWeightKg(100), mustNewDuration(30))
	totals = totals.RecordWorkout(mustNewWeightKg(150), mustNewDuration(45))
	totals = totals.RecordWorkout(mustNewWeightKg(75), mustNewDuration(20))

	if totals.Workouts() != 3 {
		t.Errorf("expected 3 workouts, got %d", totals.Workouts())
	}
	if totals.Volume() != 325 {
		t.Errorf("expected volume 325, got %f", totals.Volume())
	}
	if totals.TimeMinutes() != 95 {
		t.Errorf("expected time 95, got %d", totals.TimeMinutes())
	}
}

func TestRecordWorkout_Immutability(t *testing.T) {
	original := user.NewTotals()
	weight := mustNewWeightKg(100)
	duration := mustNewDuration(45)
	updated := original.RecordWorkout(weight, duration)

	if original.Workouts() != 0 || original.Volume() != 0 || original.TimeMinutes() != 0 {
		t.Error("original totals should be unchanged")
	}
	if updated.Workouts() != 1 || updated.Volume() != 100 || updated.TimeMinutes() != 45 {
		t.Error("updated totals should reflect recorded workout")
	}
}

func TestMethodChaining(t *testing.T) {
	totals := user.NewTotals().
		AddWorkout().
		AddVolume(mustNewWeightKg(50)).
		AddDuration(mustNewDuration(30)).
		AddWorkout().
		AddVolume(mustNewWeightKg(75))

	if totals.Workouts() != 2 {
		t.Errorf("expected 2 workouts, got %d", totals.Workouts())
	}
	if totals.Volume() != 125 {
		t.Errorf("expected volume 125, got %f", totals.Volume())
	}
	if totals.TimeMinutes() != 30 {
		t.Errorf("expected time 30, got %d", totals.TimeMinutes())
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
