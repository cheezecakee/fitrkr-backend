package user

import (
	"errors"
	"time"
)

var (
	ErrInvalidTheme      = errors.New("invalid theme")
	ErrInvalidVisibility = errors.New("invalid visibility")
)

type Theme string

const (
	Dark   Theme = "dark"
	Light  Theme = "light"
	System Theme = "system"
)

type Visibility string

const (
	Public  Visibility = "public"
	Private Visibility = "private"
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

func NewTheme(theme string) (Theme, error) {
	switch theme {
	case "dark":
		return Dark, nil
	case "light":
		return Light, nil
	case "system":
		return System, nil
	default:
		return "", ErrInvalidTheme
	}
}

func NewVisibility(visibility string) (Visibility, error) {
	switch visibility {
	case "public":
		return Public, nil
	case "private":
		return Private, nil
	default:
		return "", ErrInvalidVisibility
	}
}
