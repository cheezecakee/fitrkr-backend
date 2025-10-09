package user_test

import (
	"testing"

	"github.com/cheezecakee/fitrkr-athena/internal/core/domain/user"
)

func TestNewName_ValidCases(t *testing.T) {
	tests := []struct {
		name      string
		firstName string
		lastName  string
		want      string
	}{
		{
			name:      "valid name",
			firstName: "john",
			lastName:  "doe",
			want:      "john doe",
		},
		{
			name:      "valid name - with whitespace trimmed",
			firstName: "  john  ",
			lastName:  "  doe  ",
			want:      "john doe",
		},
		{
			name:      "valid name - minimum length",
			firstName: "Jo",
			lastName:  "Do",
			want:      "Jo Do",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := user.NewName(tt.firstName, tt.lastName)
			if err != nil {
				t.Fatalf("NewName() unexpected error = %v", err)
			}
			if got != tt.want {
				t.Errorf("NewName() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestNewName_InvalidCases(t *testing.T) {
	tests := []struct {
		name      string
		firstName string
		lastName  string
		wantErr   bool
	}{
		{
			name:      "invalid name - both empty",
			firstName: "",
			lastName:  "",
			wantErr:   true,
		},
		{
			name:      "invalid name - first name empty",
			firstName: "",
			lastName:  "doe",
			wantErr:   true,
		},
		{
			name:      "invalid name - last name empty",
			firstName: "john",
			lastName:  "",
			wantErr:   true,
		},
		{
			name:      "invalid name - first name too short",
			firstName: "j",
			lastName:  "doe",
			wantErr:   true,
		},
		{
			name:      "invalid name - last name too short",
			firstName: "john",
			lastName:  "d",
			wantErr:   true,
		},
		{
			name:      "invalid name - first name contains digit",
			firstName: "jo1n",
			lastName:  "doe",
			wantErr:   true,
		},
		{
			name:      "invalid name - last name contains digit",
			firstName: "john",
			lastName:  "d0e",
			wantErr:   true,
		},
		{
			name:      "invalid name - first name contains special character",
			firstName: "john!",
			lastName:  "doe",
			wantErr:   true,
		},
		{
			name:      "invalid name - last name contains special character",
			firstName: "john",
			lastName:  "doe?",
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := user.NewName(tt.firstName, tt.lastName)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewName() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
