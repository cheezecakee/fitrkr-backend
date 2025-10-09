package user_test

import (
	"testing"

	"github.com/cheezecakee/fitrkr-athena/internal/core/domain/user"
)

func TestNewBillingPeriod(t *testing.T) {
	tests := []struct {
		name    string
		period  string
		wantErr bool
	}{
		{
			name:    "valid period - monthly",
			period:  "monthly",
			wantErr: false,
		},
		{
			name:    "valid period - yearly",
			period:  "yearly",
			wantErr: false,
		},
		{
			name:    "valid period - with whitespace trimmed",
			period:  "  monthly  ",
			wantErr: false,
		},
		{
			name:    "valid period - uppercase",
			period:  "MONTHLY",
			wantErr: false,
		},
		{
			name:    "valid period - mixed case",
			period:  "YeArLy",
			wantErr: false,
		},
		{
			name:    "invalid period - empty",
			period:  "",
			wantErr: true,
		},
		{
			name:    "invalid period - whitespace only",
			period:  "   ",
			wantErr: true,
		},
		{
			name:    "invalid period - weekly",
			period:  "weekly",
			wantErr: true,
		},
		{
			name:    "invalid period - annual",
			period:  "annual",
			wantErr: true,
		},
		{
			name:    "invalid period - quarterly",
			period:  "quarterly",
			wantErr: true,
		},
		{
			name:    "invalid period - with digits",
			period:  "monthly1",
			wantErr: true,
		},
		{
			name:    "invalid period - with special chars",
			period:  "yearly!",
			wantErr: true,
		},
		{
			name:    "invalid period - with underscore",
			period:  "month_ly",
			wantErr: true,
		},
		{
			name:    "invalid period - misspelled",
			period:  "monthyl",
			wantErr: true,
		},
		{
			name:    "invalid period - with spaces",
			period:  "per month",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := user.NewBillingPeriod(tt.period)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewBillingPeriod() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewBillingPeriod_ValidCases(t *testing.T) {
	tests := []struct {
		name   string
		period string
		want   user.Period
	}{
		{
			name:   "monthly period",
			period: "monthly",
			want:   user.Monthly,
		},
		{
			name:   "yearly period",
			period: "yearly",
			want:   user.Yearly,
		},
		{
			name:   "uppercase converted to monthly",
			period: "MONTHLY",
			want:   user.Monthly,
		},
		{
			name:   "mixed case converted to yearly",
			period: "YeArLy",
			want:   user.Yearly,
		},
		{
			name:   "whitespace trimmed - monthly",
			period: "   monthly   ",
			want:   user.Monthly,
		},
		{
			name:   "whitespace trimmed - yearly",
			period: "  yearly  ",
			want:   user.Yearly,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := user.NewBillingPeriod(tt.period)
			if err != nil {
				t.Fatalf("NewBillingPeriod() unexpected error = %v", err)
			}
			if got != tt.want {
				t.Errorf("NewBillingPeriod() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestNewBillingPeriod_InvalidCases(t *testing.T) {
	tests := []struct {
		name   string
		period string
	}{
		{
			name:   "invalid - empty string",
			period: "",
		},
		{
			name:   "invalid - weekly",
			period: "weekly",
		},
		{
			name:   "invalid - quarterly",
			period: "quarterly",
		},
		{
			name:   "invalid - annual",
			period: "annual",
		},
		{
			name:   "invalid - daily",
			period: "daily",
		},
		{
			name:   "invalid - biannual",
			period: "biannual",
		},
		{
			name:   "invalid - random string",
			period: "notaperiod",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := user.NewBillingPeriod(tt.period)
			if err == nil {
				t.Fatalf("expected error, got period = %v", got)
			}
			if err != user.ErrInvalidBillingPeriod {
				t.Errorf("expected ErrInvalidBillingPeriod, got %v", err)
			}
		})
	}
}

func TestPeriod_StringValue(t *testing.T) {
	tests := []struct {
		name   string
		period user.Period
		want   string
	}{
		{
			name:   "monthly period string",
			period: user.Monthly,
			want:   "monthly",
		},
		{
			name:   "yearly period string",
			period: user.Yearly,
			want:   "yearly",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.period) != tt.want {
				t.Errorf("Period string = %q, want %q", string(tt.period), tt.want)
			}
		})
	}
}
