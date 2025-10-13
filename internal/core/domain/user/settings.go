package user

import (
	"time"
)

type Settings struct {
	WeightUnit WeightUnit `json:"weight_unit"`
	HeightUnit HeightUnit `json:"height_unit"`
	Theme      Theme      `json:"theme"`

	Visibility Visibility `json:"visibility"`

	EmailNotif      bool `json:"email_notif"`
	PushNotif       bool `json:"push_notif"`
	WorkoutReminder bool `json:"workout_reminder"`
	StreakReminder  bool `json:"streak_reminder"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
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
