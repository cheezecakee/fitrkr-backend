package users

import (
	"context"
	"fmt"

	"github.com/cheezecakee/logr"

	"github.com/cheezecakee/fitrkr-backend/internal/core/domain/user"
	"github.com/cheezecakee/fitrkr-backend/internal/ports"
)

type CreateAccountReq struct {
	Username  string   `json:"username"`
	Email     string   `json:"email"`
	FirstName string   `json:"first_name"`
	LastName  string   `json:"last_name"`
	Roles     []string `json:"roles"`
	Password  string   `json:"password"`
}

type CreateAccountResp struct {
	UserID string
}

func (s *Service) CreateAccount(ctx context.Context, req CreateAccountReq) (*CreateAccountResp, error) {
	username, err := user.NewUsername(req.Username)
	if err != nil {
		logr.Get().Errorf("invalid username: %v", err)
		return nil, fmt.Errorf("invalid username: %w", err)
	}

	fullName, err := user.NewName(req.FirstName, req.LastName)
	if err != nil {
		logr.Get().Errorf("invalid name: %v", err)
		return nil, fmt.Errorf("invalid name: %w", err)
	}

	email, err := user.NewEmail(req.Email)
	if err != nil {
		logr.Get().Errorf("invalid email: %v", err)
		return nil, fmt.Errorf("invalid email: %w", err)
	}

	roles, err := user.NewRole(req.Roles)
	if err != nil {
		logr.Get().Errorf("invalid role/s %v", err)
		return nil, fmt.Errorf("invalid role/s %w", err)
	}

	password, err := user.NewPassword(req.Password)
	if err != nil {
		logr.Get().Errorf("invalid password %v", err)
		return nil, fmt.Errorf("invalid password %w", err)
	}

	existingUser, err := s.userRepo.GetByUsername(ctx, string(username))
	if err != nil && err != ports.ErrUserNotFound {
		logr.Get().Errorf("failed to check username: %v", err)
		return nil, fmt.Errorf("failed to check username: %w", err)
	}
	if existingUser != nil {
		logr.Get().Errorf("username already exists: %v", err)
		return nil, ErrDuplicateUsername
	}

	existingUser, err = s.userRepo.GetByEmail(ctx, string(email))
	if err != nil && err != ports.ErrUserNotFound {
		logr.Get().Errorf("failed to check email: %v", err)
		return nil, fmt.Errorf("failed to check email: %w", err)
	}
	if existingUser != nil {
		logr.Get().Errorf("email already exists: %v", err)
		return nil, ErrDuplicateEmail
	}

	user := user.New(username, fullName, email, roles, password)

	if err := s.userRepo.Add(ctx, user); err != nil {
		logr.Get().Errorf("failed to add a user: %v", err)
		return nil, fmt.Errorf("failed to add a user: %w", err)
	}

	logr.Get().Info("New user account created")
	return &CreateAccountResp{UserID: user.ID.String()}, nil
}
