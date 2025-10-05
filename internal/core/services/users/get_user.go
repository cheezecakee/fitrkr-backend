package users

import (
	"context"
	"time"

	"github.com/cheezecakee/logr"
	"github.com/google/uuid"

	"github.com/cheezecakee/fitrkr-backend/internal/core/domain/user"
)

type (
	GetUserByIDReq struct {
		ID uuid.UUID `json:"id"`
	}
	GetUserByUsernameReq struct {
		Username string `json:"username"`
	}
	GetUserByEmailReq struct {
		Email string `json:"email"`
	}
	GetUserResp struct {
		ID        string `json:"id"`
		Username  string `json:"username"`
		Email     string `json:"email"`
		FullName  string `json:"full_name"`
		Roles     []string
		CreatedAt string `json:"created_at"`
		UpdatedAt string `json:"updated_at"`
	}
)

func (s *Service) GetByID(ctx context.Context, req GetUserByIDReq) (*GetUserResp, error) {
	u, err := s.userRepo.GetByID(ctx, req.ID.String())
	if err != nil {
		logr.Get().Errorf("failed to get user by id: %v", err)
		return nil, ErrUserNotFound
	}

	return mapUserToResponse(u), nil
}

func (s *Service) GetByUsername(ctx context.Context, req GetUserByUsernameReq) (*GetUserResp, error) {
	u, err := s.userRepo.GetByUsername(ctx, req.Username)
	if err != nil {
		logr.Get().Errorf("failed to get user by username: %v", err)
		return nil, ErrUserNotFound
	}

	return mapUserToResponse(u), nil
}

func (s *Service) GetByEmail(ctx context.Context, req GetUserByEmailReq) (*GetUserResp, error) {
	u, err := s.userRepo.GetByEmail(ctx, req.Email)
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
		Roles:     user.RolesToStrings(u.Roles()),
		CreatedAt: u.CreatedAt.Format(time.RFC3339),
		UpdatedAt: u.UpdatedAt.Format(time.RFC3339),
	}
}
