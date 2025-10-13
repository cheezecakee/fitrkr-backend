package users

import (
	"context"
	"fmt"

	"github.com/cheezecakee/logr"

	"github.com/cheezecakee/fitrkr-athena/internal/core/domain/user"
)

type UpdateSettingsReq struct {
	UserID     string
	WeightUnit *string
	HeightUnit *string
	Theme      *string

	Visibility *string

	EmailNotif      *bool
	PushNotif       *bool
	WorkoutReminder *bool
	StreakReminder  *bool
}

func (s *Service) UpdateSettings(ctx context.Context, req UpdateSettingsReq) error {
	settings, err := s.userRepo.GetSettingsByID(ctx, req.UserID)
	if err != nil {
		logr.Get().Errorf("failed to get user settings: %v", err)
		return fmt.Errorf("failed to get user settings: %w", err)
	}

	if req.WeightUnit != nil {
		weightUnit, err := user.NewWeightUnit(*req.WeightUnit)
		if err != nil {
			logr.Get().Errorf("failed to create new weight unit: %v", err)
			return fmt.Errorf("failed to create new weight unit: %w", err)
		}
		settings.WeightUnit = weightUnit
	}

	if req.HeightUnit != nil {
		heightUnit, err := user.NewHeightUnit(*req.HeightUnit)
		if err != nil {
			logr.Get().Errorf("failed to create new height unit: %v", err)
			return fmt.Errorf("failed to create new height unit: %w", err)
		}
		settings.HeightUnit = heightUnit
	}

	if req.Theme != nil {
		theme, err := user.NewTheme(*req.Theme)
		if err != nil {
			logr.Get().Errorf("failed to create new theme unit: %v", err)
			return fmt.Errorf("failed to create new theme unit: %w", err)
		}
		settings.Theme = theme
	}

	if req.Visibility != nil {
		visibility, err := user.NewVisibility(*req.Visibility)
		if err != nil {
			logr.Get().Errorf("failed to create new visibility: %v", err)
			return fmt.Errorf("failed to create new visibility: %w", err)
		}
		settings.Visibility = visibility
	}

	if req.EmailNotif != nil {
		settings.EmailNotif = *req.EmailNotif
	}

	if req.PushNotif != nil {
		settings.PushNotif = *req.PushNotif
	}

	if req.WorkoutReminder != nil {
		settings.WorkoutReminder = *req.WorkoutReminder
	}

	if req.StreakReminder != nil {
		settings.StreakReminder = *req.StreakReminder
	}

	settings.Touch()

	err = s.userRepo.UpdateSettings(ctx, *settings, req.UserID)
	if err != nil {
		logr.Get().Errorf("failed to update user settings: %v", err)
		return fmt.Errorf("failed to update user settings: %w", err)
	}

	return nil
}
