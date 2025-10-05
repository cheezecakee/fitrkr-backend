// Package users
package users

import (
	"context"
	"errors"

	"github.com/cheezecakee/fitrkr-backend/internal/ports"
)

var (
	ErrUserNotFound      = errors.New("user does not exist")
	ErrDuplicateEmail    = errors.New("email already exists")
	ErrDuplicateUsername = errors.New("username already exists")
)

type UserService interface {
	CreateAccount(context.Context, CreateAccountReq) (*CreateAccountResp, error)
	GetByID(ctx context.Context, req GetUserByIDReq) (*GetUserResp, error)
	GetByUsername(ctx context.Context, req GetUserByUsernameReq) (*GetUserResp, error)
	GetByEmail(ctx context.Context, req GetUserByEmailReq) (*GetUserResp, error)
}

type Service struct {
	userRepo ports.UserRepo
}

func NewService(userRepo ports.UserRepo) *Service {
	return &Service{
		userRepo: userRepo,
	}
}
