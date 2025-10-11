package user_test

import (
	"testing"
	"time"

	"github.com/cheezecakee/fitrkr-athena/internal/core/domain/user"
)

func TestNewSettings(t *testing.T) {
	settings := user.NewSettings()

	// Verify defaults are set
	if settings.WeightUnit != user.Kg {
		t.Errorf("expected WeightUnit to be Kg, got %v", settings.WeightUnit)
	}
	if settings.HeightUnit != user.Cm {
		t.Errorf("expected HeightUnit to be Cm, got %v", settings.HeightUnit)
	}
	if settings.Theme != user.System {
		t.Errorf("expected Theme to be System, got %v", settings.Theme)
	}
	if settings.Visibility != user.Public {
		t.Errorf("expected Visibility to be Public, got %v", settings.Visibility)
	}

	// Verify notification defaults are all true
	if !settings.EmailNotif {
		t.Error("expected EmailNotif to default to true")
	}
	if !settings.PushNotif {
		t.Error("expected PushNotif to default to true")
	}
	if !settings.WorkoutReminder {
		t.Error("expected WorkoutReminder to default to true")
	}
	if !settings.StreakReminder {
		t.Error("expected StreakReminder to default to true")
	}

	// Verify timestamps are set
	if settings.CreatedAt.IsZero() {
		t.Error("expected CreatedAt to be set")
	}
	if settings.UpdatedAt.IsZero() {
		t.Error("expected UpdatedAt to be set")
	}

	// Timestamps should be very close (within 1 second)
	if settings.UpdatedAt.Before(settings.CreatedAt) {
		t.Error("expected UpdatedAt to be after or equal to CreatedAt")
	}
	timeDiff := settings.UpdatedAt.Sub(settings.CreatedAt)
	if timeDiff > time.Second {
		t.Errorf("expected CreatedAt and UpdatedAt to be within 1 second, got %v", timeDiff)
	}
}

func TestNewSettings_TimestampsAreRecent(t *testing.T) {
	before := time.Now()
	settings := user.NewSettings()
	after := time.Now()

	// Verify CreatedAt is between before and after
	if settings.CreatedAt.Before(before) || settings.CreatedAt.After(after) {
		t.Errorf("expected CreatedAt to be between %v and %v, got %v", before, after, settings.CreatedAt)
	}

	// Verify UpdatedAt is between before and after
	if settings.UpdatedAt.Before(before) || settings.UpdatedAt.After(after) {
		t.Errorf("expected UpdatedAt to be between %v and %v, got %v", before, after, settings.UpdatedAt)
	}
}

func TestSettings_NotificationDefaultsCanBeChanged(t *testing.T) {
	settings := user.NewSettings()

	// All notifications should start as true
	if !settings.EmailNotif || !settings.PushNotif || !settings.WorkoutReminder || !settings.StreakReminder {
		t.Error("expected all notifications to default to true")
	}

	// User can disable notifications
	settings.EmailNotif = false
	settings.WorkoutReminder = false

	if settings.EmailNotif {
		t.Error("expected EmailNotif to be false after change")
	}
	if settings.WorkoutReminder {
		t.Error("expected WorkoutReminder to be false after change")
	}

	// Others should still be true
	if !settings.PushNotif {
		t.Error("expected PushNotif to remain true")
	}
	if !settings.StreakReminder {
		t.Error("expected StreakReminder to remain true")
	}
}

func TestSettings_CanUpdateUnitsAndPreferences(t *testing.T) {
	settings := user.NewSettings()

	// Verify initial state
	if settings.WeightUnit != user.Kg {
		t.Error("expected initial WeightUnit to be Kg")
	}
	if settings.Theme != user.System {
		t.Error("expected initial Theme to be System")
	}

	// Update preferences
	settings.WeightUnit = user.Lb
	settings.HeightUnit = user.Ft
	settings.Theme = user.Dark
	settings.Visibility = user.Private

	// Verify updates
	if settings.WeightUnit != user.Lb {
		t.Errorf("expected WeightUnit to be Lb, got %v", settings.WeightUnit)
	}
	if settings.HeightUnit != user.Ft {
		t.Errorf("expected HeightUnit to be Ft, got %v", settings.HeightUnit)
	}
	if settings.Theme != user.Dark {
		t.Errorf("expected Theme to be Dark, got %v", settings.Theme)
	}
	if settings.Visibility != user.Private {
		t.Errorf("expected Visibility to be Private, got %v", settings.Visibility)
	}
}

func TestSettings_TouchUpdatesTimestamp(t *testing.T) {
	settings := user.NewSettings()
	originalUpdatedAt := settings.UpdatedAt

	// Small delay to ensure time difference is noticeable
	time.Sleep(10 * time.Millisecond)

	settings.Touch()

	if !settings.UpdatedAt.After(originalUpdatedAt) {
		t.Error("expected UpdatedAt to be updated after Touch()")
	}
}
