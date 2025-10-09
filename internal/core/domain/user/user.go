// Package user
package user

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID `json:"id"`
	Username     Username  `json:"username"`
	FullName     string    `json:"full_name"`
	Email        Email     `json:"email"`
	Password     Password  `json:"-"`
	roles        []Role
	CreatedAt    time.Time    `json:"created_at"`
	UpdatedAt    time.Time    `json:"updated_at"`
	Stats        Stats        `json:"stats"`
	Subscription Subscription `json:"subscription"`
	Settings     Settings     `json:"settings"`
}

func New(username Username, fullName string, email Email, roles []Role, password Password, stats Stats, subscription Subscription, settings Settings) User {
	return User{
		ID:           uuid.New(),
		Username:     username,
		FullName:     fullName,
		Email:        email,
		roles:        roles,
		Password:     password,
		Stats:        stats,
		Subscription: subscription,
		Settings:     settings,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
}

// Roles getter
func (u *User) Roles() Roles {
	copyRoles := make(Roles, len(u.roles))
	copy(copyRoles, u.roles)
	return copyRoles
}

type UserSnapshot struct {
	ID           uuid.UUID
	Username     Username
	FullName     string
	Email        Email
	Roles        Roles
	Stats        Stats
	Subscription Subscription
	Settings     Settings
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// Reconstitute rebuilds an existing user from persistence (for DB loading)
func (s UserSnapshot) Reconstitute() *User {
	return &User{
		ID:           s.ID,
		Username:     s.Username,
		FullName:     s.FullName,
		Email:        s.Email,
		roles:        []Role(s.Roles),
		Stats:        s.Stats,
		Subscription: s.Subscription,
		Settings:     s.Settings,
		CreatedAt:    s.CreatedAt,
		UpdatedAt:    s.UpdatedAt,
	}
}
