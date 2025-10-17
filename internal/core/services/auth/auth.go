// Package auth
package auth

import (
	"context"
	"errors"

	"github.com/cheezecakee/fitrkr-athena/internal/core/ports"
)

var (
	ErrPasswordIncorrect   = errors.New("incorrect password")
	ErrRefreshTokenExpired = errors.New("refresh token expired")
	ErrRefreshTokenRevoked = errors.New("refresh token revoked")
)

type AuthService interface {
	Login(ctx context.Context, req LoginReq) (LoginResp, error)
	Logout(ctx context.Context, req LogoutReq) error
	Revoke(ctx context.Context, req RevokeTokenReq) error
	Refresh(ctx context.Context, req RefreshReq) (RefreshResp, error)
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
