package users

import (
	"context"
	"fmt"
	"time"

	"github.com/cheezecakee/logr"

	"github.com/cheezecakee/fitrkr-athena/internal/core/domain/user"
	"github.com/cheezecakee/fitrkr-athena/internal/core/ports"
)

type UpdateUserReq struct {
	ID        string `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

func (s *Service) Update(ctx context.Context, req UpdateUserReq) error {
	existingUser, err := s.userRepo.GetByID(ctx, req.ID)
	if err != nil {
		logr.Get().Errorf("failed to get user: %v", err)
		return fmt.Errorf("failed to get user: %w", err)
	}

	if req.Username != "" {
		username, err := user.NewUsername(req.Username)
		if err != nil {
			logr.Get().Errorf("invalid username: %v", err)
			return fmt.Errorf("invalid username: %w", err)
		}
		if string(username) != string(existingUser.Username) {
			userWithUsername, err := s.userRepo.GetByUsername(ctx, string(username))
			if err != nil && err != ports.ErrUserNotFound {
				logr.Get().Errorf("failed to check username: %v", err)
				return fmt.Errorf("failed to check username: %w", err)
			}
			if userWithUsername != nil && userWithUsername.ID != existingUser.ID {
				logr.Get().Errorf("username already exists: %v", err)
				return ErrDuplicateUsername
			}
		}

		existingUser.Username = username
	}

	if req.Email != "" {
		email, err := user.NewEmail(req.Email)
		if err != nil {
			logr.Get().Errorf("invalid email: %v", err)
			return fmt.Errorf("invalid email: %w", err)
		}

		if string(email) != string(existingUser.Email) {
			userWithEmail, err := s.userRepo.GetByEmail(ctx, string(email))
			if err != nil && err != ports.ErrUserNotFound {
				logr.Get().Errorf("failed to check email: %v", err)
				return fmt.Errorf("failed to check email: %w", err)
			}
			if userWithEmail != nil && userWithEmail.ID != existingUser.ID {
				logr.Get().Errorf("email already exists: %v", err)
				return ErrDuplicateEmail
			}
		}

		existingUser.Email = email
	}

	if req.FirstName != "" || req.LastName != "" {
		firstName := req.LastName
		lastName := req.LastName

		if firstName == "" || lastName == "" {
			// Parse existing full name to get parts
			// For now, just require both or neither
			if firstName == "" && lastName == "" {
				// Keep existing
			} else {
				return fmt.Errorf("both first and last name must be provided together")
			}
		} else {
			fullName, err := user.NewName(req.FirstName, req.LastName)
			if err != nil {
				logr.Get().Errorf("invalid name: %v", err)
				return fmt.Errorf("invalid name: %w", err)
			}

			existingUser.FullName = fullName
		}
	}

	existingUser.UpdatedAt = time.Now()
	user := &user.User{
		ID:        existingUser.ID,
		Username:  existingUser.Username,
		FullName:  existingUser.FullName,
		Email:     existingUser.Email,
		CreatedAt: existingUser.CreatedAt,
		UpdatedAt: existingUser.UpdatedAt,
	}
	err = s.userRepo.Update(ctx, *user)
	if err != nil {
		logr.Get().Errorf("error update user: %v", err)
		return fmt.Errorf("error update user: %w", err)
	}

	logr.Get().Info("User updated successfully")
	return nil
}
