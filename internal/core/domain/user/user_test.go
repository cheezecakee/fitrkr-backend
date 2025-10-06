package user_test

import (
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/cheezecakee/fitrkr-backend/internal/core/domain/user"
)

func TestNewUser(t *testing.T) {
	username := user.Username("testuser")
	fullName := "Test User"
	email := user.Email("test@example.com")
	password := user.Password("secret123!")
	roles := user.Roles{user.RoleAdmin}

	u := user.New(username, fullName, email, roles, password)

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

	if time.Since(u.CreatedAt) > 2*time.Second {
		t.Errorf("expected CreatedAt to be recent, got %v", u.CreatedAt)
	}

	if time.Since(u.UpdatedAt) > 2*time.Second {
		t.Errorf("expected UpdatedAt to be recent, got %v", u.UpdatedAt)
	}

	// Roles should be stored and retrievable
	if got := u.Roles(); !equalRoles(got, roles) {
		t.Errorf("expected roles %v, got %v", roles, got)
	}
}

func TestReconstituteUser(t *testing.T) {
	id := uuid.New()
	username := user.Username("restored")
	fullName := "Restored User"
	email := user.Email("restore@example.com")
	roles := user.Roles{user.RoleUser}
	createdAt := time.Now().Add(-10 * time.Hour)
	updatedAt := time.Now().Add(-5 * time.Hour)

	u := user.Reconstitute(id, username, fullName, email, roles, createdAt, updatedAt)

	if u.ID != id {
		t.Errorf("expected ID %v, got %v", id, u.ID)
	}

	if u.CreatedAt != createdAt {
		t.Errorf("expected CreatedAt %v, got %v", createdAt, u.CreatedAt)
	}

	if u.UpdatedAt != updatedAt {
		t.Errorf("expected UpdatedAt %v, got %v", updatedAt, u.UpdatedAt)
	}

	if got := u.Roles(); !equalRoles(got, roles) {
		t.Errorf("expected roles %v, got %v", roles, got)
	}
}

func TestRoleGetter_IsImmutable(t *testing.T) {
	roles := user.Roles{user.RoleAdmin}
	u := user.New("testuser", "Full Name", "test@example.com", roles, "secret123!")

	copyRoles := u.Roles()
	copyRoles[0] = user.RoleModerator

	if equalRoles(copyRoles, u.Roles()) {
		t.Error("expected roles to be immutable, but they were modified")
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
