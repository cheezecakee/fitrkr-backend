package auth

import (
	"context"
	"fmt"

	"github.com/cheezecakee/logr"
)

type RevokeTokenReq struct {
	Token string
}

func (s *Service) Revoke(ctx context.Context, req RevokeTokenReq) error {
	token, err := s.authRepo.GetByToken(ctx, req.Token)
	if err != nil {
		logr.Get().Errorf("failed to revoke token: %v", err)
		return fmt.Errorf("failed to revoke token: %w", err)
	}

	token.Revoke()

	err = s.authRepo.Update(ctx, *token)
	if err != nil {
		logr.Get().Errorf("failed to update token: %v", err)
		return fmt.Errorf("failed to update token: %w", err)
	}

	return nil
}
