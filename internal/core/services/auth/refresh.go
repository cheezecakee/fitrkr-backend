package auth

import (
	"context"
	"fmt"

	"github.com/cheezecakee/logr"
	"github.com/google/uuid"

	"github.com/cheezecakee/fitrkr-athena/internal/core/domain/auth"
)

type RefreshReq struct {
	Token string
}

type RefreshResp struct {
	UserID uuid.UUID
	Roles  []string
	Token  string
}

func (s *Service) Refresh(ctx context.Context, req RefreshReq) (RefreshResp, error) {
	currentToken, err := s.authRepo.GetByToken(ctx, req.Token)
	if err != nil {
		logr.Get().Errorf("failed to get refresh token: %v", err)
		return RefreshResp{}, fmt.Errorf("failed to get refresh token: %w", err)
	}

	if currentToken.IsExpired() {
		return RefreshResp{}, ErrRefreshTokenExpired
	}
	if currentToken.IsRevoked {
		return RefreshResp{}, ErrRefreshTokenRevoked
	}

	currentToken.Revoke()
	err = s.authRepo.Update(ctx, *currentToken)
	if err != nil {
		logr.Get().Errorf("failed to update refresh token: %v", err)
		return RefreshResp{}, fmt.Errorf("failed to update refresh token: %w", err)
	}

	// Create a new token
	token, err := auth.NewRefreshToken(currentToken.UserID)
	if err != nil {
		logr.Get().Errorf("failed to generate refresh token: %v", err)
		return RefreshResp{}, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	err = s.authRepo.Add(ctx, token)
	if err != nil {
		logr.Get().Errorf("failed to create refresh token: %v", err)
		return RefreshResp{}, fmt.Errorf("failed to create refresh token: %w", err)
	}

	user, err := s.userRepo.GetByID(ctx, currentToken.UserID.String())
	if err != nil {
		logr.Get().Errorf("failed to get user: %v", err)
		return RefreshResp{}, fmt.Errorf("failed to get user: %w", err)
	}

	logr.Get().Info("user token refreshed")
	return RefreshResp{Token: token.Token, UserID: user.ID, Roles: user.Roles}, nil
}
