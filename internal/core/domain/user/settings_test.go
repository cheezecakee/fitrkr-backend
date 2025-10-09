package user_test

import (
	"testing"
	"time"

	"github.com/cheezecakee/fitrkr-athena/internal/core/domain/user"
)

func TestNewSettings(t *testing.T) {
	weightUnit := user.Kg
	heightUnit := user.Cm
	theme := user.Dark
	visibility := user.Public

	settings := user.NewSettings(weightUnit, heightUnit, theme, visibility)

	// Verify provided values are set correctly
	if settings.WeightUnit != user.Kg {
		t.Errorf("expected WeightUnit to be Kg, got %v", settings.WeightUnit)
	}
	if settings.HeightUnit != user.Cm {
		t.Errorf("expected HeightUnit to be Cm, got %v", settings.HeightUnit)
	}
	if settings.Theme != user.Dark {
		t.Errorf("expected Theme to be Dark, got %v", settings.Theme)
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

func TestNewSettings_DifferentCombinations(t *testing.T) {
	tests := []struct {
		name       string
		weightUnit user.WeightUnit
		heightUnit user.HeightUnit
		theme      user.Theme
		visibility user.Visibility
	}{
		{
			name:       "metric units with dark theme and public",
			weightUnit: user.Kg,
			heightUnit: user.Cm,
			theme:      user.Dark,
			visibility: user.Public,
		},
		{
			name:       "imperial units with light theme and private",
			weightUnit: user.Lb,
			heightUnit: user.FtIn,
			theme:      user.Light,
			visibility: user.Private,
		},
		{
			name:       "mixed units with system theme and public",
			weightUnit: user.Kg,
			heightUnit: user.FtIn,
			theme:      user.System,
			visibility: user.Public,
		},
		{
			name:       "imperial weight metric height",
			weightUnit: user.Lb,
			heightUnit: user.Cm,
			theme:      user.Dark,
			visibility: user.Private,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			settings := user.NewSettings(tt.weightUnit, tt.heightUnit, tt.theme, tt.visibility)

			if settings.WeightUnit != tt.weightUnit {
				t.Errorf("expected WeightUnit %v, got %v", tt.weightUnit, settings.WeightUnit)
			}
			if settings.HeightUnit != tt.heightUnit {
				t.Errorf("expected HeightUnit %v, got %v", tt.heightUnit, settings.HeightUnit)
			}
			if settings.Theme != tt.theme {
				t.Errorf("expected Theme %v, got %v", tt.theme, settings.Theme)
			}
			if settings.Visibility != tt.visibility {
				t.Errorf("expected Visibility %v, got %v", tt.visibility, settings.Visibility)
			}
		})
	}
}

func TestNewSettings_TimestampsAreRecent(t *testing.T) {
	before := time.Now()
	settings := user.NewSettings(user.Kg, user.Cm, user.System, user.Public)
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
	settings := user.NewSettings(user.Kg, user.Cm, user.Dark, user.Public)

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
	settings := user.NewSettings(user.Kg, user.Cm, user.Light, user.Private)

	// Verify initial state
	if settings.WeightUnit != user.Kg {
		t.Error("expected initial WeightUnit to be Kg")
	}
	if settings.Theme != user.Light {
		t.Error("expected initial Theme to be Light")
	}

	// Update preferences
	settings.WeightUnit = user.Lb
	settings.HeightUnit = user.FtIn
	settings.Theme = user.Dark
	settings.Visibility = user.Public

	// Verify updates
	if settings.WeightUnit != user.Lb {
		t.Errorf("expected WeightUnit to be Lb, got %v", settings.WeightUnit)
	}
	if settings.HeightUnit != user.FtIn {
		t.Errorf("expected HeightUnit to be FtIn, got %v", settings.HeightUnit)
	}
	if settings.Theme != user.Dark {
		t.Errorf("expected Theme to be Dark, got %v", settings.Theme)
	}
	if settings.Visibility != user.Public {
		t.Errorf("expected Visibility to be Public, got %v", settings.Visibility)
	}
}
