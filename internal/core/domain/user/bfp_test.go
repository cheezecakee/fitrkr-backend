package user_test

import (
	"testing"

	"github.com/cheezecakee/fitrkr-backend/internal/core/domain/user"
)

func TestBFP(t *testing.T) {
	tests := []struct {
		name    string
		bfp     float64
		wantErr bool
	}{
		{
			name:    "valid bfp",
			bfp:     20,
			wantErr: false,
		},
		{
			name:    "valid bfp - exactly 0",
			bfp:     0,
			wantErr: false,
		},
		{
			name:    "valid bfp - exactly 100",
			bfp:     100,
			wantErr: false,
		},
		{
			name:    "invalid bfp - exactly over 100",
			bfp:     100.1,
			wantErr: true,
		},
		{
			name:    "invalid bfp - negative value",
			bfp:     -10,
			wantErr: true,
		},
		{
			name:    "ivalid bfp - over 100",
			bfp:     140,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := user.NewBFP(tt.bfp)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewBFP() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
