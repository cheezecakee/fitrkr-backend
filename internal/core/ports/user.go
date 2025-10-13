// Package ports
package ports

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"

	"github.com/cheezecakee/fitrkr-athena/internal/core/domain/user"
)

var (
	ErrUserNotFound      = errors.New("user does not exist")
	ErrDuplicateEmail    = errors.New("email already exists")
	ErrDuplicateUsername = errors.New("username already exists")
)

type User struct {
	ID           uuid.UUID
	Username     user.Username
	FullName     string
	PasswordHash string
	Email        user.Email
	Roles        []string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type UserRepo interface {
	Add(ctx context.Context, user user.User) error
	GetByUsername(ctx context.Context, username string) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	GetByID(ctx context.Context, id string) (*User, error)
	Update(ctx context.Context, user user.User) error
	Delete(ctx context.Context, id string) error

	AddStats(ctx context.Context, stats user.Stats, userID string) error
	GetStatsByID(ctx context.Context, userID string) (*user.Stats, error)

	AddSubscription(ctx context.Context, sub user.Subscription, userID string) error
	GetSubscriptionByID(ctx context.Context, userID string) (*user.Subscription, error)
	UpdateSubscription(ctx context.Context, sub user.Subscription, userID string) error

	AddSettings(ctx context.Context, settings user.Settings, userID string) error
	GetSettingsByID(ctx context.Context, userID string) (*user.Settings, error)
	UpdateSettings(ctx context.Context, settings user.Settings, userID string) error
}
