// Package ports
package ports

import (
	"context"
	"errors"

	"github.com/cheezecakee/fitrkr-backend/internal/core/domain/user"
)

var (
	ErrUserNotFound      = errors.New("user does not exist")
	ErrDuplicateEmail    = errors.New("email already exists")
	ErrDuplicateUsername = errors.New("username already exists")
)

type UserRepo interface {
	Add(ctx context.Context, u user.User) error
	GetByUsername(ctx context.Context, username string) (*user.User, error)
	GetByEmail(ctx context.Context, email string) (*user.User, error)
	GetByID(ctx context.Context, id string) (*user.User, error)
	Update(ctx context.Context, u user.User) error
}
