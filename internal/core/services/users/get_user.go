package users

import (
	"context"
	"time"

	"github.com/cheezecakee/logr"

	"github.com/cheezecakee/fitrkr-athena/internal/core/domain/user"
)

type GetUserResp struct {
	ID        string `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	FullName  string `json:"full_name"`
	Roles     []string
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func (s *Service) GetByID(ctx context.Context, req string) (*GetUserResp, error) {
	u, err := s.userRepo.GetByID(ctx, req)
	if err != nil {
		logr.Get().Errorf("failed to get user by id: %v", err)
		return nil, ErrUserNotFound
	}

	return mapUserToResponse(u), nil
}

func (s *Service) GetByUsername(ctx context.Context, req string) (*GetUserResp, error) {
	u, err := s.userRepo.GetByUsername(ctx, req)
	if err != nil {
		logr.Get().Errorf("failed to get user by username: %v", err)
		return nil, ErrUserNotFound
	}

	return mapUserToResponse(u), nil
}

func (s *Service) GetByEmail(ctx context.Context, req string) (*GetUserResp, error) {
	u, err := s.userRepo.GetByEmail(ctx, req)
	if err != nil {
		logr.Get().Errorf("failed to get user by email: %v", err)
		return nil, ErrUserNotFound
	}

	return mapUserToResponse(u), nil
}

func mapUserToResponse(u *user.User) *GetUserResp {
	return &GetUserResp{
		ID:        u.ID.String(),
		Username:  string(u.Username),
		Email:     string(u.Email),
		FullName:  u.FullName,
		Roles:     u.Roles().ToStrings(),
		CreatedAt: u.CreatedAt.Format(time.RFC3339),
		UpdatedAt: u.UpdatedAt.Format(time.RFC3339),
	}
}
