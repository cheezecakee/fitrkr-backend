package user_test

import (
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/cheezecakee/fitrkr-athena/internal/core/domain/user"
)

func TestNewUser(t *testing.T) {
	username := user.Username("testuser")
	fullName := "Test User"
	email := user.Email("test@example.com")
	password := user.Password("Secure123!")
	roles := user.Roles{user.RoleAdmin}
	sub := user.NewSubscription()
	settings := user.NewSettings()
	stats := user.NewStats()

	u := user.New(username, fullName, email, roles, password, stats, sub, settings)

	if u.ID == uuid.Nil {
		t.Error("expected generated UUID, got nil")
	}

	if u.Username != username {
		t.Errorf("expected username %v, got %v", username, u.Username)
	}

	if u.FullName != fullName {
		t.Errorf("expected full name %v, got %v", fullName, u.FullName)
	}

	if u.Email != email {
		t.Errorf("expected email %v, got %v", email, u.Email)
	}

	// Roles are directly accessible now
	if !equalRoles(u.Roles, roles) {
		t.Errorf("expected roles %v, got %v", roles, u.Roles)
	}

	if time.Since(u.CreatedAt) > 2*time.Second {
		t.Errorf("expected CreatedAt to be recent, got %v", u.CreatedAt)
	}

	if time.Since(u.UpdatedAt) > 2*time.Second {
		t.Errorf("expected UpdatedAt to be recent, got %v", u.UpdatedAt)
	}
}

func TestNewRolesValidation(t *testing.T) {
	// Valid roles
	roles, err := user.NewRoles([]string{"user", "admin"})
	if err != nil {
		t.Errorf("expected no error for valid roles, got %v", err)
	}

	if !equalRoles(roles, user.Roles{user.RoleUser, user.RoleAdmin}) {
		t.Errorf("expected roles to be user and admin")
	}

	// Invalid role
	_, err = user.NewRoles([]string{"invalid"})
	if err == nil {
		t.Error("expected error for invalid role")
	}

	// Empty defaults to RoleUser
	roles, _ = user.NewRoles(nil)
	if !equalRoles(roles, user.Roles{user.RoleUser}) {
		t.Errorf("expected default role to be user")
	}
}

func equalRoles(a, b user.Roles) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
