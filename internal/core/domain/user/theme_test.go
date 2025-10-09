package user_test

import (
	"testing"

	"github.com/cheezecakee/fitrkr-athena/internal/core/domain/user"
)

func TestTheme(t *testing.T) {
	tests := []struct {
		name    string
		theme   string
		wantErr bool
	}{
		{
			name:    "valid theme",
			theme:   "dark",
			wantErr: false,
		},
		{
			name:    "valid theme - with whitespace trimmed",
			theme:   " system ",
			wantErr: false,
		},
		{
			name:    "valid theme - empty",
			theme:   "",
			wantErr: false,
		},
		{
			name:    "valid theme - uppercase",
			theme:   "LIGHT",
			wantErr: false,
		},
		{
			name:    "invalid theme - with digits",
			theme:   "light1",
			wantErr: true,
		},
		{
			name:    "invalid theme - with underscore",
			theme:   "da_rk",
			wantErr: true,
		},
		{
			name:    "invalid theme",
			theme:   "gray",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.theme, func(t *testing.T) {
			_, err := user.NewTheme(tt.theme)
			if (err != nil) != tt.wantErr {
				t.Errorf("Newtheme() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewTheme_ValidCases(t *testing.T) {
	tests := []struct {
		name  string
		theme string
		want  user.Theme
	}{
		{
			name:  "dark theme",
			theme: "dark",
			want:  user.Dark,
		},
		{
			name:  "system theme",
			theme: "system",
			want:  user.System,
		},
		{
			name:  "light theme",
			theme: "light",
			want:  user.Light,
		},
		{
			name:  "default empty theme",
			theme: "",
			want:  user.System,
		},
		{
			name:  "upper case",
			theme: "DarK",
			want:  user.Dark,
		},
		{
			name:  "with white spacing",
			theme: "   light  ",
			want:  user.Light,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := user.NewTheme(tt.theme)
			if err != nil {
				t.Fatalf("NewTheme() unexpected error = %v", err)
			}
			if got != tt.want {
				t.Errorf("NewTheme() = %q, want %q", got, tt.want)
			}
		})
	}
}
