package user_test

import (
	"testing"

	"github.com/cheezecakee/fitrkr-athena/internal/core/domain/user"
)

func TestNewWeightUnit(t *testing.T) {
	tests := []struct {
		name    string
		unit    string
		want    user.WeightUnit
		wantErr bool
	}{
		{
			name:    "valid unit - kg",
			unit:    "kg",
			want:    user.Kg,
			wantErr: false,
		},
		{
			name:    "valid unit - lb",
			unit:    "lb",
			want:    user.Lb,
			wantErr: false,
		},
		{
			name:    "valid unit - uppercase",
			unit:    "KG",
			want:    user.Kg,
			wantErr: false,
		},
		{
			name:    "valid unit - mixedcase",
			unit:    "lB",
			want:    user.Lb,
			wantErr: false,
		},
		{
			name:    "valid unit - empty defaults to kg",
			unit:    "",
			want:    user.Kg,
			wantErr: false,
		},
		{
			name:    "invalid unit - stone",
			unit:    "st",
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
			unit:    "Kg$",
			want:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := user.NewWeightUnit(tt.unit)
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got none")
				}
				if err != user.ErrInvalidWeightUnit {
					t.Fatalf("expected ErrInvalidWeightUnit, got %v", err)
				}
			} else if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Errorf("NewWeightUnit() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestWeight_Conversion(t *testing.T) {
	tests := []struct {
		name   string
		value  float64
		unit   user.WeightUnit
		wantKg float64
		wantLb float64
	}{
		{
			name:   "100 kg to conversions",
			value:  100,
			unit:   user.Kg,
			wantKg: 100,
			wantLb: 220.462,
		},
		{
			name:   "220.462 lb to conversions",
			value:  220.462,
			unit:   user.Lb,
			wantKg: 100,
			wantLb: 220.462,
		},
		{
			name:   "1 kg to conversions",
			value:  1,
			unit:   user.Kg,
			wantKg: 1,
			wantLb: 2.20462,
		},
		{
			name:   "1 lb to conversions",
			value:  1,
			unit:   user.Lb,
			wantKg: 0.453592,
			wantLb: 1,
		},
		{
			name:   "2.5 kg to conversions",
			value:  2.5,
			unit:   user.Kg,
			wantKg: 2.5,
			wantLb: 5.51155,
		},
		{
			name:   "3.3 lb to conversions",
			value:  3.3,
			unit:   user.Lb,
			wantKg: 1.4968536,
			wantLb: 3.3,
		},
		{
			name:   "very small kg value",
			value:  0.0001,
			unit:   user.Kg,
			wantKg: 0.0001,
			wantLb: 0.000220462,
		},
		{
			name:   "500 kg to conversions",
			value:  500,
			unit:   user.Kg,
			wantKg: 500,
			wantLb: 1102.31,
		},
		{
			name:   "1000 lb to conversions",
			value:  1000,
			unit:   user.Lb,
			wantKg: 453.592,
			wantLb: 1000,
		},
		{
			name:   "75 kg round-trip",
			value:  75,
			unit:   user.Kg,
			wantKg: 75,
			wantLb: 165.3465,
		},
	}

	const tolerance = 0.01

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w, err := user.NewWeight(tt.value, tt.unit)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			gotKg := w.Display(user.Kg)
			gotLb := w.Display(user.Lb)

			if diff := gotKg - tt.wantKg; diff < -tolerance || diff > tolerance {
				t.Errorf("Display(Kg) = %v, want %v", gotKg, tt.wantKg)
			}
			if diff := gotLb - tt.wantLb; diff < -tolerance || diff > tolerance {
				t.Errorf("Display(Lb) = %v, want %v", gotLb, tt.wantLb)
			}
		})
	}
}

func TestNewWeight(t *testing.T) {
	tests := []struct {
		name    string
		value   float64
		unit    user.WeightUnit
		wantErr bool
	}{
		{
			name:    "valid weight - lb",
			value:   140.20,
			unit:    user.Lb,
			wantErr: false,
		},
		{
			name:    "valid weight - kg",
			value:   49.0,
			unit:    user.Kg,
			wantErr: false,
		},
		{
			name:    "invalid weight - negative",
			value:   -39.9,
			unit:    user.Kg,
			wantErr: true,
		},
		{
			name:    "invalid weight - zero",
			value:   0,
			unit:    user.Kg,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := user.NewWeight(tt.value, tt.unit)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewWeight() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
