// Package users
package users

import (
	"context"

	"github.com/cheezecakee/fitrkr-backend/internal/ports"
)

type UserService interface {
	CreateAccount(context.Context, CreateAccountReq) (*CreateAccountResp, error)
}

type Service struct {
	userRepo ports.UserRepo
}

func NewService(userRepo ports.UserRepo) *Service {
	return &Service{
		userRepo: userRepo,
	}
}
