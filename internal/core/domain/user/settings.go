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

func NewSettings() Settings {
	return Settings{
		WeightUnit:      Kg,
		HeightUnit:      Cm,
		Theme:           System,
		Visibility:      Public,
		EmailNotif:      true,
		PushNotif:       true,
		WorkoutReminder: true,
		StreakReminder:  true,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}
}

func (s *Settings) Touch() {
	s.UpdatedAt = time.Now()
}
