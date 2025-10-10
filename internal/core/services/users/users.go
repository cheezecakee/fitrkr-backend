// Package users
package users

import (
	"context"
	"errors"

	"github.com/cheezecakee/fitrkr-athena/internal/ports"
)

var (
	ErrUserNotFound      = errors.New("user does not exist")
	ErrDuplicateEmail    = errors.New("email already exists")
	ErrDuplicateUsername = errors.New("username already exists")
)

type UserService interface {
	CreateAccount(context.Context, CreateAccountReq) (*CreateAccountResp, error)
	GetByID(ctx context.Context, id string) (*GetUserResp, error)
	GetByUsername(ctx context.Context, username string) (*GetUserResp, error)
	GetByEmail(ctx context.Context, email string) (*GetUserResp, error)
	Update(ctx context.Context, req UpdateUserReq, id string) error
	Delete(ctx context.Context, id string) error
}

type Service struct {
	userRepo ports.UserRepo
}

func NewService(userRepo ports.UserRepo) *Service {
	return &Service{
		userRepo: userRepo,
	}
}
