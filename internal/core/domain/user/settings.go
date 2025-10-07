package user

import (
	"time"
)

type Settings struct {
	WeightUnit WeightUnit
	HeightUnit HeightUnit
	Theme      Theme

	Visibility Visibility

	EmailNotif      bool
	PushNotif       bool
	WorkoutReminder bool
	StreakReminder  bool

	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewSettings(weightUnit WeightUnit, heightUnit HeightUnit, theme Theme, visibility Visibility) Settings {
	return Settings{
		WeightUnit:      weightUnit,
		HeightUnit:      heightUnit,
		Theme:           theme,
		Visibility:      visibility,
		EmailNotif:      true,
		PushNotif:       true,
		WorkoutReminder: true,
		StreakReminder:  true,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}
}
