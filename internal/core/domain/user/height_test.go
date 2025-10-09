package user_test

import (
	"testing"

	"github.com/cheezecakee/fitrkr-athena/internal/core/domain/user"
)

func TestNewHeightUnit(t *testing.T) {
	tests := []struct {
		name    string
		unit    string
		want    user.HeightUnit
		wantErr bool
	}{
		{
			name:    "valid unit - cm",
			unit:    "cm",
			want:    user.Cm,
			wantErr: false,
		},
		{
			name:    "valid unit - ft_in",
			unit:    "ft_in",
			want:    user.FtIn,
			wantErr: false,
		},
		{
			name:    "valid unit - uppercase",
			unit:    "CM",
			want:    user.Cm,
			wantErr: false,
		},
		{
			name:    "valid unit - mixedcase",
			unit:    "Ft_iN",
			want:    user.FtIn,
			wantErr: false,
		},
		{
			name:    "valid unit - empty (defaults to cm)",
			unit:    "",
			want:    user.Cm,
			wantErr: false,
		},
		{
			name:    "invalid unit - meters",
			unit:    "m",
			want:    "",
			wantErr: true,
		},
		{
			name:    "invalid unit - random",
			unit:    "foo",
			want:    "",
			wantErr: true,
		},
		{
			name:    "invalid unit - number",
			unit:    "123",
			want:    "",
			wantErr: true,
		},
		{
			name:    "invalid unit - special chars",
			unit:    "Cm$",
			want:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := user.NewHeightUnit(tt.unit)
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got none")
				}
				if err != user.ErrInvalidHeightUnit {
					t.Fatalf("expected ErrInvalidHeightUnit, got %v", err)
				}
			} else if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Errorf("NewHeightUnit() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestHeight_Conversion(t *testing.T) {
	tests := []struct {
		name   string
		value  float64
		unit   user.HeightUnit
		wantCm float64
	}{
		// --- Basic conversions ---
		{
			name:   "100 cm to ft/in",
			value:  100,
			unit:   user.Cm,
			wantCm: 100,
		},
		{
			name:   "180 cm to ft/in",
			value:  180,
			unit:   user.Cm,
			wantCm: 180,
		},
		{
			name:   "60 inches (5 ft) to cm",
			value:  60,
			unit:   user.FtIn,
			wantCm: 152.4,
		},
		{
			name:   "72 inches (6 ft) to cm",
			value:  72,
			unit:   user.FtIn,
			wantCm: 182.88,
		},

		// --- Edge cases ---
		{
			name:   "very small height",
			value:  1,
			unit:   user.Cm,
			wantCm: 1,
		},
		{
			name:   "large height",
			value:  250,
			unit:   user.Cm,
			wantCm: 250,
		},
	}

	const tolerance = 0.01 // floating-point tolerance (cm)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h, err := user.NewHeight(tt.value, tt.unit)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			gotCm := h.ToCm()
			gotFt, gotIn := h.ToFtIn()
			backToCm := float64(gotFt*12)*2.54 + gotIn*2.54

			if diff := gotCm - tt.wantCm; diff < -tolerance || diff > tolerance {
				t.Errorf("ToCm() = %v, want %v", gotCm, tt.wantCm)
			}
			if diff := backToCm - tt.wantCm; diff < -tolerance || diff > tolerance {
				t.Errorf("round-trip cm mismatch: got %v cm, want %v cm", backToCm, tt.wantCm)
			}
		})
	}
}

func TestNewHeight(t *testing.T) {
	tests := []struct {
		name    string
		value   float64
		unit    user.HeightUnit
		wantErr bool
	}{
		{
			name:    "valid height - cm",
			value:   175,
			unit:    user.Cm,
			wantErr: false,
		},
		{
			name:    "valid height - ft/in",
			value:   70,
			unit:    user.FtIn,
			wantErr: false,
		},
		{
			name:    "invalid height - negative number",
			value:   -39.9,
			unit:    user.Cm,
			wantErr: true,
		},
		{
			name:    "invalid height - zero",
			value:   0,
			unit:    user.Cm,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := user.NewHeight(tt.value, tt.unit)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewHeight() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
