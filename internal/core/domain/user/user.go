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
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Stats        `json:"stats"`
	Subscription `json:"subscription"`
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

func (u *User) Roles() []Role {
	copyRoles := make([]Role, len(u.roles))
	copy(copyRoles, u.roles)
	return copyRoles
}
