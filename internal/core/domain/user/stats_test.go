package user_test

import (
	"testing"
	"time"

	"github.com/cheezecakee/fitrkr-athena/internal/core/domain/user"
)

func TestNewStats(t *testing.T) {
	streak := mustNewStreak(2) // restDays = 2

	stats := user.NewStats(streak)

	// Verify streak is set correctly
	if stats.Streak.Current() != 0 {
		t.Errorf("expected streak current to be 0, got %d", stats.Streak.Current())
	}
	if stats.Streak.RestDays() != 2 {
		t.Errorf("expected streak rest days to be 2, got %d", stats.Streak.RestDays())
	}

	// Verify totals is initialized
	if stats.Totals.Workouts() != 0 {
		t.Errorf("expected totals workouts to be 0, got %d", stats.Totals.Workouts())
	}
	if stats.Totals.Volume() != 0 {
		t.Errorf("expected totals volume to be 0, got %f", stats.Totals.Volume())
	}
	if stats.Totals.TimeMinutes() != 0 {
		t.Errorf("expected totals time to be 0, got %d", stats.Totals.TimeMinutes())
	}

	// Verify optional fields are nil
	if stats.Weight != nil {
		t.Error("expected Weight to be nil")
	}
	if stats.Height != nil {
		t.Error("expected Height to be nil")
	}
	if stats.BFP != nil {
		t.Error("expected BFP to be nil")
	}

	// Verify timestamps are set and reasonable
	if stats.CreatedAt.IsZero() {
		t.Error("expected CreatedAt to be set")
	}
	if stats.UpdatedAt.IsZero() {
		t.Error("expected UpdatedAt to be set")
	}

	// Timestamps should be very close (within 1 second)
	if stats.UpdatedAt.Before(stats.CreatedAt) {
		t.Error("expected UpdatedAt to be after or equal to CreatedAt")
	}
	timeDiff := stats.UpdatedAt.Sub(stats.CreatedAt)
	if timeDiff > time.Second {
		t.Errorf("expected CreatedAt and UpdatedAt to be within 1 second, got %v", timeDiff)
	}
}

func TestNewStats_TimestampsAreRecent(t *testing.T) {
	before := time.Now()
	streak := mustNewStreak(2)
	stats := user.NewStats(streak)
	after := time.Now()

	// Verify CreatedAt is between before and after
	if stats.CreatedAt.Before(before) || stats.CreatedAt.After(after) {
		t.Errorf("expected CreatedAt to be between %v and %v, got %v", before, after, stats.CreatedAt)
	}

	// Verify UpdatedAt is between before and after
	if stats.UpdatedAt.Before(before) || stats.UpdatedAt.After(after) {
		t.Errorf("expected UpdatedAt to be between %v and %v, got %v", before, after, stats.UpdatedAt)
	}
}

func TestStats_OptionalFieldsCanBeSet(t *testing.T) {
	streak := mustNewStreak(2)
	stats := user.NewStats(streak)

	// Set optional fields
	weight := mustNewWeightKg(75.5)
	height := mustNewHeight(180, user.Cm)
	bfp := mustNewBFP(15.5)

	stats.Weight = &weight
	stats.Height = &height
	stats.BFP = &bfp

	// Verify they're set correctly
	if stats.Weight == nil || stats.Weight.Value() != 75.5 {
		t.Error("Weight not set correctly")
	}
	if stats.Height == nil || stats.Height.ToCm() != 180 {
		t.Error("Height not set correctly")
	}
	if stats.BFP == nil || float64(*stats.BFP) != 15.5 {
		t.Error("BFP not set correctly")
	}
}

func TestStats_TotalsCanBeUpdated(t *testing.T) {
	streak := mustNewStreak(2)
	stats := user.NewStats(streak)

	// Update totals
	stats.Totals = stats.Totals.RecordWorkout(
		mustNewWeightKg(100),
		mustNewDuration(45),
	)

	// Verify update
	if stats.Totals.Workouts() != 1 {
		t.Errorf("expected 1 workout, got %d", stats.Totals.Workouts())
	}
	if stats.Totals.Volume() != 100 {
		t.Errorf("expected volume 100, got %f", stats.Totals.Volume())
	}
	if stats.Totals.TimeMinutes() != 45 {
		t.Errorf("expected 45 minutes, got %d", stats.Totals.TimeMinutes())
	}
}

// Helper functions
func mustNewStreak(restDays int) user.Streak {
	s, err := user.NewStreak(restDays)
	if err != nil {
		panic(err)
	}
	return s
}

func mustNewHeight(value float64, unit user.HeightUnit) user.Height {
	h, err := user.NewHeight(value, unit)
	if err != nil {
		panic(err)
	}
	return h
}

func mustNewBFP(percentage float64) user.BFP {
	b, err := user.NewBFP(percentage)
	if err != nil {
		panic(err)
	}
	return b
}
