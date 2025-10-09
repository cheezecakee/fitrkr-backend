package users

import (
	"context"
	"fmt"

	"github.com/cheezecakee/logr"

	"github.com/cheezecakee/fitrkr-athena/internal/ports"
)

func (s *Service) Delete(ctx context.Context, req string) error {
	err := s.userRepo.Delete(ctx, req)
	if err != nil {
		if err == ports.ErrUserNotFound {
			logr.Get().Error("user not found")
			return ErrUserNotFound
		}
		logr.Get().Errorf("failed to delete user: %v", err)
		return fmt.Errorf("failed to delete user: %w", err)
	}

	logr.Get().Info("User deleted successfully")
	return nil
}
