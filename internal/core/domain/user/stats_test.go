package user_test

import (
	"testing"
	"time"

	"github.com/cheezecakee/fitrkr-athena/internal/core/domain/user"
)

func TestNewStats(t *testing.T) {
	stats := user.NewStats()

	// Verify streak is initialized with defaults
	if stats.Streak.RestDays != 2 {
		t.Errorf("expected streak RestDays to be 2, got %d", stats.Streak.RestDays)
	}
	if stats.Streak.Current != 0 {
		t.Errorf("expected streak Current to be 0, got %d", stats.Streak.Current)
	}
	if stats.Streak.Longest != 0 {
		t.Errorf("expected streak Longest to be 0, got %d", stats.Streak.Longest)
	}

	// Verify totals is initialized
	if stats.Totals.Workouts != 0 {
		t.Errorf("expected totals workouts to be 0, got %d", stats.Totals.Workouts)
	}
	if stats.Totals.Lifted != 0 {
		t.Errorf("expected totals lifted to be 0, got %f", stats.Totals.Lifted)
	}
	if stats.Totals.Time != 0 {
		t.Errorf("expected totals time to be 0, got %d", stats.Totals.Time)
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
	stats := user.NewStats()
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
	stats := user.NewStats()

	// Set optional fields
	weight, err := user.NewWeight(75.5, user.Kg)
	if err != nil {
		t.Fatalf("failed to create weight: %v", err)
	}

	height, err := user.NewHeight(180, user.Cm)
	if err != nil {
		t.Fatalf("failed to create height: %v", err)
	}

	bfp, err := user.NewBFP(15.5)
	if err != nil {
		t.Fatalf("failed to create BFP: %v", err)
	}

	stats.Weight = &weight
	stats.Height = &height
	stats.BFP = &bfp

	// Verify they're set correctly
	if stats.Weight == nil {
		t.Error("Weight should not be nil after setting")
	}
	if stats.Height == nil {
		t.Error("Height should not be nil after setting")
	}
	if stats.BFP == nil {
		t.Error("BFP should not be nil after setting")
	}
}

func TestStats_TotalsCanBeUpdated(t *testing.T) {
	stats := user.NewStats()

	// Create weight and duration for recording
	weight, err := user.NewWeight(100, user.Kg)
	if err != nil {
		t.Fatalf("failed to create weight: %v", err)
	}

	duration, err := user.NewDuration(45)
	if err != nil {
		t.Fatalf("failed to create duration: %v", err)
	}

	// Update totals via RecordWorkout (mutates in place)
	stats.Totals.RecordWorkout(weight, duration)

	// Verify update
	if stats.Totals.Workouts != 1 {
		t.Errorf("expected 1 workout, got %d", stats.Totals.Workouts)
	}
	if stats.Totals.Lifted != 100 {
		t.Errorf("expected lifted 100, got %f", stats.Totals.Lifted)
	}
	if stats.Totals.Time != 45 {
		t.Errorf("expected 45 minutes, got %d", stats.Totals.Time)
	}
}
