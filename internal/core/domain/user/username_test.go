package user_test

import (
	"testing"

	"github.com/cheezecakee/fitrkr-backend/internal/core/domain/user"
)

func TestUsername(t *testing.T) {
	tests := []struct {
		name     string
		username string
		wantErr  bool
	}{
		{
			name:     "valid username - with whitespace trimmed",
			username: " cheezecake ",
			wantErr:  false,
		},
		{
			name:     "valid username",
			username: "cheezecake",
			wantErr:  false,
		},
		{
			name:     "valid username - with digits",
			username: "cheezecak33",
			wantErr:  false,
		},
		{
			name:     "valid username - with underscore",
			username: "cheeze_cake",
			wantErr:  false,
		},
		{
			name:     "valid username - at minimum length", // Min 3 char
			username: "che",
			wantErr:  false,
		},
		{
			name:     "valid username - at maximum length", // Max 20 char
			username: "cheezecake1234567890",
			wantErr:  false,
		},
		{
			name:     "invalid username - empty",
			username: "",
			wantErr:  true,
		},
		{
			name:     "invalid username - too short",
			username: "ch",
			wantErr:  true,
		},
		{
			name:     "invalid username - too long",
			username: "ThisIsAVeryLongUsername",
			wantErr:  true,
		},
		{
			name:     "invalid username - contains special",
			username: "cheeze@cake",
			wantErr:  true,
		},
		{
			name:     "invalid username - contains spaces",
			username: "cheeze cake",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.username, func(t *testing.T) {
			_, err := user.NewUsername(tt.username)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewUsername() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewUsername_ValidCases(t *testing.T) {
	tests := []struct {
		name     string
		username string
		want     string
	}{
		{
			name:     "converts uppercase to lower",
			username: "ChEeZeCaKe",
			want:     "cheezecake",
		},
		{
			name:     "no change",
			username: "cheezecake",
			want:     "cheezecake",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := user.NewUsername(tt.username)
			if err != nil {
				t.Fatalf("NewUsername() unexpected error = %v", err)
			}
			if string(got) != tt.want {
				t.Errorf("NewUsername() = %q, want %q", got, tt.want)
			}
		})
	}
}
