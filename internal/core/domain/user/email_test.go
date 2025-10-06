package user_test

import (
	"testing"

	"github.com/cheezecakee/fitrkr-backend/internal/core/domain/user"
)

func TestNewEmail(t *testing.T) {
	tests := []struct {
		name    string
		email   string
		wantErr bool
	}{
		{
			name:    "valid email",
			email:   "test@example.com",
			wantErr: false,
		},
		{
			name:    "valid email - with subdomain",
			email:   "test@email.example.com",
			wantErr: false,
		},
		{
			name:    "valid email - with plus",
			email:   "test+tag@example.com",
			wantErr: false,
		},
		{
			name:    "invalid email - no @",
			email:   "testexample.com",
			wantErr: true,
		},
		{
			name:    "valid email - empty",
			email:   "",
			wantErr: true,
		},
		{
			name:    "invalid email - no domain",
			email:   "test@",
			wantErr: true,
		},
		{
			name:    "invalid email - multiple @",
			email:   "test@@example.com",
			wantErr: true,
		},
		{
			name:    "invalid email - no local part",
			email:   "@example.com",
			wantErr: true,
		},
		{
			name:    "invalid email - spaces",
			email:   "test @example.com",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := user.NewEmail(tt.email)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewEmail() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
