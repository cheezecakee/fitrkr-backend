package auth

import (
	"context"
	"fmt"

	"github.com/cheezecakee/logr"
)

type LogoutReq struct {
	Token string
}

func (s *Service) Logout(ctx context.Context, req LogoutReq) error {
	refreshToken, err := s.authRepo.GetByToken(ctx, req.Token)
	if err != nil {
		logr.Get().Errorf("failed to get refresh token: %v", err)
		return fmt.Errorf("failed to get refresh token: %w", err)
	}

	refreshToken.Revoke()

	err = s.authRepo.Update(ctx, refreshToken)
	if err != nil {
		logr.Get().Errorf("failed to update user refresh token: %v", err)
		return fmt.Errorf("failed to update user refresh token: %w", err)
	}
	logr.Get().Info("User logged out")
	return nil
}
