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

func New(username Username, fullName string, email Email, roles []Role, password Password) User {
	return User{
		ID:        uuid.New(),
		Username:  username,
		FullName:  fullName,
		Email:     email,
		roles:     roles,
		Password:  password,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

// Reconstitute rebuilds an existing user from persistence (for DB loading)
func Reconstitute(id uuid.UUID, username Username, fullName string, email Email, roles Roles, createdAt, updatedAt time.Time) *User {
	return &User{
		ID:        id,
		Username:  username,
		FullName:  fullName,
		Email:     email,
		roles:     []Role(roles),
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}
}

// Roles getter
func (u *User) Roles() Roles {
	copyRoles := make(Roles, len(u.roles))
	copy(copyRoles, u.roles)
	return copyRoles
}
