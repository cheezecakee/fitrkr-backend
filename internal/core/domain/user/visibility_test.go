package user_test

import (
	"testing"

	"github.com/cheezecakee/fitrkr-athena/internal/core/domain/user"
)

func TestVisibility(t *testing.T) {
	tests := []struct {
		name       string
		visibility string
		wantErr    bool
	}{
		{
			name:       "valid visibility",
			visibility: "public",
			wantErr:    false,
		},
		{
			name:       "valid visibility - with whitespace trimmed",
			visibility: " private ",
			wantErr:    false,
		},
		{
			name:       "valid visibility - empty",
			visibility: "",
			wantErr:    false,
		},
		{
			name:       "valid visibility - uppercase",
			visibility: "PUBLIC",
			wantErr:    false,
		},
		{
			name:       "invalid visibility - with digits",
			visibility: "private1",
			wantErr:    true,
		},
		{
			name:       "invalid visibility - with underscore",
			visibility: "private_",
			wantErr:    true,
		},
		{
			name:       "invalid visibility",
			visibility: "none",
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.visibility, func(t *testing.T) {
			_, err := user.NewVisibility(tt.visibility)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewVisibility() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewVisibility_ValidCases(t *testing.T) {
	tests := []struct {
		name       string
		visibility string
		want       user.Visibility
	}{
		{
			name:       "public visibility",
			visibility: "public",
			want:       user.Public,
		},
		{
			name:       "private visibility",
			visibility: "private",
			want:       user.Private,
		},
		{
			name:       "default empty visibility",
			visibility: "",
			want:       user.Public,
		},
		{
			name:       "upper case",
			visibility: "PRIVATE",
			want:       user.Private,
		},
		{
			name:       "with white spacing",
			visibility: "   public  ",
			want:       user.Public,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := user.NewVisibility(tt.visibility)
			if err != nil {
				t.Fatalf("NewVisibility() unexpected error = %v", err)
			}
			if got != tt.want {
				t.Errorf("NewVisibility() = %q, want %q", got, tt.want)
			}
		})
	}
}
