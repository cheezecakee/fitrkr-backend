// Package user
package user

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID    `json:"id"`
	Username     Username     `json:"username"`
	FullName     string       `json:"full_name"`
	Email        Email        `json:"email"`
	Password     Password     `json:"-"`
	Roles        Roles        `json:"roles"`
	CreatedAt    time.Time    `json:"created_at"`
	UpdatedAt    time.Time    `json:"updated_at"`
	Stats        Stats        `json:"stats"`
	Subscription Subscription `json:"subscription"`
	Settings     Settings     `json:"settings"`
}

func New(username Username, fullName string, email Email, roles Roles, password Password, stats Stats, subscription Subscription, settings Settings) User {
	return User{
		ID:           uuid.New(),
		Username:     username,
		FullName:     fullName,
		Email:        email,
		Roles:        roles,
		Password:     password,
		Stats:        stats,
		Subscription: subscription,
		Settings:     settings,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
}
