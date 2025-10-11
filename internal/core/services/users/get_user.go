package users

import (
	"context"
	"time"

	"github.com/cheezecakee/logr"
	"github.com/google/uuid"

	"github.com/cheezecakee/fitrkr-athena/internal/core/domain/user"
	"github.com/cheezecakee/fitrkr-athena/internal/core/ports"
)

type GetUserResp struct {
	ID        uuid.UUID     `json:"id"`
	Username  user.Username `json:"username"`
	Email     user.Email    `json:"email"`
	FullName  string        `json:"full_name"`
	Roles     user.Roles    `json:"roles"`
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
}

type GetUserByIDReq struct {
	ID string
}

func (s *Service) GetByID(ctx context.Context, req GetUserByIDReq) (*GetUserResp, error) {
	u, err := s.userRepo.GetByID(ctx, req.ID)
	if err != nil {
		logr.Get().Errorf("failed to get user by id: %v", err)
		return nil, ErrUserNotFound
	}

	return mapUserToResponse(u), nil
}

type GetUserByUsernameReq struct {
	Username string
}

func (s *Service) GetByUsername(ctx context.Context, req GetUserByUsernameReq) (*GetUserResp, error) {
	u, err := s.userRepo.GetByUsername(ctx, req.Username)
	if err != nil {
		logr.Get().Errorf("failed to get user by username: %v", err)
		return nil, ErrUserNotFound
	}

	return mapUserToResponse(u), nil
}

type GetUserByEmailReq struct {
	Email string
}

func (s *Service) GetByEmail(ctx context.Context, req GetUserByEmailReq) (*GetUserResp, error) {
	u, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		logr.Get().Errorf("failed to get user by email: %v", err)
		return nil, ErrUserNotFound
	}

	return mapUserToResponse(u), nil
}

func mapUserToResponse(u *ports.User) *GetUserResp {
	return &GetUserResp{
		ID:        u.ID,
		Username:  user.Username(u.Username),
		Email:     user.Email(u.Email),
		FullName:  u.FullName,
		Roles:     user.StringsToRoles(u.Roles),
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}
