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
			name:    "valid unit",
			unit:    "kg",
			want:    user.Kg,
			wantErr: false,
		},
		{
			name:    "valid unit",
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
			name:    "valid unit - empty",
			unit:    "",
			want:    user.Kg, // Should default to kg
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
					t.Fatal("expected error, not none")
				}
				if err != user.ErrInvalidWeightUnit {
					t.Fatalf("expected ErrInvalidWeightUnit, got %v", err)
				}
			} else if err != nil {
				t.Fatalf("unexpected err: %v", err)
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
		// --- Basic conversions ---
		{
			name:   "100 kg to lb and back",
			value:  100,
			unit:   user.Kg,
			wantKg: 100,
			wantLb: 220.462,
		},
		{
			name:   "220.462 lb to kg and back",
			value:  220.462,
			unit:   user.Lb,
			wantKg: 100,
			wantLb: 220.462,
		},

		// --- Small values ---
		{
			name:   "1 kg to lb",
			value:  1,
			unit:   user.Kg,
			wantKg: 1,
			wantLb: 2.20462,
		},
		{
			name:   "1 lb to kg",
			value:  1,
			unit:   user.Lb,
			wantKg: 0.453592,
			wantLb: 1,
		},

		// --- Decimal precision ---
		{
			name:   "2.5 kg to lb",
			value:  2.5,
			unit:   user.Kg,
			wantKg: 2.5,
			wantLb: 5.51155,
		},
		{
			name:   "3.3 lb to kg",
			value:  3.3,
			unit:   user.Lb,
			wantKg: 1.4968536,
			wantLb: 3.3,
		},

		// --- Zero and near-zero ---
		{
			name:   "very small kg value",
			value:  0.0001,
			unit:   user.Kg,
			wantKg: 0.0001,
			wantLb: 0.000220462,
		},

		// --- Large realistic values ---
		{
			name:   "500 kg to lb",
			value:  500,
			unit:   user.Kg,
			wantKg: 500,
			wantLb: 1102.31,
		},
		{
			name:   "1000 lb to kg",
			value:  1000,
			unit:   user.Lb,
			wantKg: 453.592,
			wantLb: 1000,
		},

		// --- Round-trip accuracy test case ---
		{
			name:   "round-trip kg → lb → kg",
			value:  75,
			unit:   user.Kg,
			wantKg: 75,
			wantLb: 165.3465,
		},
	}

	const tolerance = 0.001 // floating-point tolerance

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w, err := user.NewWeight(tt.value, tt.unit)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			gotKg := w.ToKg()
			gotLb := w.ToLbs()

			if diff := gotKg - tt.wantKg; diff < -tolerance || diff > tolerance {
				t.Errorf("ToKg() = %v, want %v", gotKg, tt.wantKg)
			}
			if diff := gotLb - tt.wantLb; diff < -tolerance || diff > tolerance {
				t.Errorf("ToLbs() = %v, want %v", gotLb, tt.wantLb)
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
			name:    "invalid weight - negative number",
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
