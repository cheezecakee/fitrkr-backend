package auth

import (
	"context"
	"fmt"

	"github.com/cheezecakee/logr"
	"github.com/google/uuid"

	"github.com/cheezecakee/fitrkr-athena/internal/core/domain/auth"
	"github.com/cheezecakee/fitrkr-athena/internal/core/domain/user"
)

type LoginReq struct {
	Username string        `json:"username"`
	Password user.Password `json:"password"`
}

type LoginResp struct {
	RefreshToken string
	UserID       uuid.UUID
	Roles        []string
}

func (s *Service) Login(ctx context.Context, req LoginReq) (LoginResp, error) {
	user, err := s.userRepo.GetByUsername(ctx, req.Username)
	if err != nil {
		logr.Get().Errorf("failed to get user: %v", err)
		return LoginResp{}, fmt.Errorf("failed to get user: %w", err)
	}

	if !req.Password.Verify(user.PasswordHash) {
		logr.Get().Errorf("password: %v, hash: %v", req.Password, user.PasswordHash)
		logr.Get().Error("password incorrect")
		return LoginResp{}, fmt.Errorf("%v", ErrPasswordIncorrect)
	}

	token, err := auth.NewRefreshToken(user.ID)
	if err != nil {
		logr.Get().Errorf("failed to generate refresh token: %v", err)
		return LoginResp{}, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	err = s.authRepo.Add(ctx, token)
	if err != nil {
		logr.Get().Errorf("failed to create refresh token: %v", err)
		return LoginResp{}, fmt.Errorf("failed to create refresh token: %w", err)
	}

	return LoginResp{
		RefreshToken: token.Token,
		UserID:       user.ID,
		Roles:        user.Roles,
	}, nil
}
