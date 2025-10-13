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
			name:    "valid unit - ft",
			unit:    "ft",
			want:    user.Ft,
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
			unit:    "Ft",
			want:    user.Ft,
			wantErr: false,
		},
		{
			name:    "valid unit - empty defaults to cm",
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
		wantCm user.HeightValue
		wantFt user.HeightValue
	}{
		{
			name:   "100 cm to conversions",
			value:  100,
			unit:   user.Cm,
			wantCm: 100,
			wantFt: 3.28084,
		},
		{
			name:   "180 cm to conversions",
			value:  180,
			unit:   user.Cm,
			wantCm: 180,
			wantFt: 5.90551,
		},
		{
			name:   "5 ft to cm",
			value:  5,
			unit:   user.Ft,
			wantCm: 152.4,
			wantFt: 5,
		},
		{
			name:   "6 ft to cm",
			value:  6,
			unit:   user.Ft,
			wantCm: 182.88,
			wantFt: 6,
		},
		{
			name:   "very small height",
			value:  1,
			unit:   user.Cm,
			wantCm: 1,
			wantFt: 0.0328084,
		},
		{
			name:   "large height",
			value:  250,
			unit:   user.Cm,
			wantCm: 250,
			wantFt: 8.2021,
		},
	}

	const tolerance = 0.01

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h, err := user.NewHeight(tt.value, tt.unit)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			gotCm := h.Display(user.Cm)
			gotFt := h.Display(user.Ft)

			if diff := gotCm - tt.wantCm; diff < -tolerance || diff > tolerance {
				t.Errorf("Display(Cm) = %v, want %v", gotCm, tt.wantCm)
			}
			if diff := gotFt - tt.wantFt; diff < -tolerance || diff > tolerance {
				t.Errorf("Display(Ft) = %v, want %v", gotFt, tt.wantFt)
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
			name:    "valid height - ft",
			value:   70,
			unit:    user.Ft,
			wantErr: false,
		},
		{
			name:    "invalid height - negative",
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
