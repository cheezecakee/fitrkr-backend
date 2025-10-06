package user_test

import (
	"reflect"
	"testing"

	"github.com/cheezecakee/fitrkr-backend/internal/core/domain/user"
)

func TestNewRoles_ValidCases(t *testing.T) {
	tests := []struct {
		name  string
		roles []string
		want  user.Roles
	}{
		{
			name:  "empty role slice returns user role",
			roles: []string{},
			want:  user.Roles{user.RoleUser},
		},
		{
			name:  "nil role slice returns user role",
			roles: nil,
			want:  user.Roles{user.RoleUser},
		},
		{
			name:  "return user role",
			roles: []string{"user"},
			want:  user.Roles{user.RoleUser},
		},
		{
			name:  "return admin role",
			roles: []string{"admin"},
			want:  user.Roles{user.RoleAdmin},
		},
		{
			name:  "return moderator role",
			roles: []string{"moderator"},
			want:  user.Roles{user.RoleModerator},
		},
		{
			name:  "return all roles",
			roles: []string{"user", "admin", "moderator"},
			want:  user.Roles{user.RoleUser, user.RoleAdmin, user.RoleModerator},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := user.NewRoles(tt.roles)
			if err != nil {
				t.Fatalf("NewRoles() unexpected error = %v", err)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewRoles() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewRoles_InvalidCases(t *testing.T) {
	tests := []struct {
		name  string
		roles []string
	}{
		{
			name:  "contains invalid roles",
			roles: []string{"invalid"},
		},
		{
			name:  "mixed valid and invalid roles",
			roles: []string{"user", "banana"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := user.NewRoles(tt.roles)
			if err == nil {
				t.Fatalf("expected error, got roles = %v", got)
			}
			if err != user.ErrInvalidRole {
				t.Errorf("expected ErrInvalidRole, got %v", err)
			}
		})
	}
}

func TestRoles_ToStrings(t *testing.T) {
	roles := user.Roles{user.RoleUser, user.RoleAdmin}
	want := []string{"user", "admin"}

	got := roles.ToStrings()

	if !reflect.DeepEqual(got, want) {
		t.Errorf("ToStrings() = %v, want %v", got, want)
	}
}

func TestRoles_FromStrings(t *testing.T) {
	strs := []string{"moderator", "admin"}
	want := user.Roles{user.RoleModerator, user.RoleAdmin}

	got := user.StringsToRoles(strs)

	if !reflect.DeepEqual(got, want) {
		t.Errorf("StringsToRoles() = %v, want %v", got, want)
	}
}
