// Package auth
package auth

import (
	"context"

	"github.com/cheezecakee/fitrkr-athena/internal/core/ports"
)

var ErrPasswordIncorrect = "incorrect password"

type AuthService interface {
	Login(ctx context.Context, req LoginReq) (LoginResp, error)
	Logout(ctx context.Context, req LogoutReq) error
	Revoke(ctx context.Context, req RevokeTokenReq) error
}

type Service struct {
	authRepo ports.AuthRepo
	userRepo ports.UserRepo
}

func NewService(authRepo ports.AuthRepo, userRepo ports.UserRepo) *Service {
	return &Service{
		authRepo: authRepo,
		userRepo: userRepo,
	}
}
