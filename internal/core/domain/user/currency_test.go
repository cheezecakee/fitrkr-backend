package user_test

import (
	"testing"

	"github.com/cheezecakee/fitrkr-athena/internal/core/domain/user"
)

func TestNewCurrency(t *testing.T) {
	tests := []struct {
		name    string
		code    string
		wantErr bool
	}{
		{
			name:    "valid currency - USD",
			code:    "USD",
			wantErr: false,
		},
		{
			name:    "valid currency - EUR",
			code:    "EUR",
			wantErr: false,
		},
		{
			name:    "valid currency - GBP",
			code:    "GBP",
			wantErr: false,
		},
		{
			name:    "valid currency - CAD",
			code:    "CAD",
			wantErr: false,
		},
		{
			name:    "valid currency - AUD",
			code:    "AUD",
			wantErr: false,
		},
		{
			name:    "valid currency - BRL",
			code:    "BRL",
			wantErr: false,
		},
		{
			name:    "valid currency - RMB",
			code:    "RMB",
			wantErr: false,
		},
		{
			name:    "valid currency - lowercase",
			code:    "usd",
			wantErr: false,
		},
		{
			name:    "valid currency - mixed case",
			code:    "EuR",
			wantErr: false,
		},
		{
			name:    "valid currency - with whitespace",
			code:    "  GBP  ",
			wantErr: false,
		},
		{
			name:    "valid currency - empty defaults to USD",
			code:    "",
			wantErr: false,
		},
		{
			name:    "valid currency - whitespace only defaults to USD",
			code:    "   ",
			wantErr: false,
		},
		{
			name:    "invalid currency - too short",
			code:    "US",
			wantErr: true,
		},
		{
			name:    "invalid currency - too long",
			code:    "USDD",
			wantErr: true,
		},
		{
			name:    "invalid currency - not supported",
			code:    "JPY",
			wantErr: true,
		},
		{
			name:    "invalid currency - not supported",
			code:    "MXN",
			wantErr: true,
		},
		{
			name:    "invalid currency - invalid code",
			code:    "XXX",
			wantErr: true,
		},
		{
			name:    "invalid currency - with numbers",
			code:    "US1",
			wantErr: true,
		},
		{
			name:    "invalid currency - with special chars",
			code:    "US$",
			wantErr: true,
		},
		{
			name:    "invalid currency - single char",
			code:    "$",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := user.NewCurrency(tt.code)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewCurrency() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewCurrency_ValidCases(t *testing.T) {
	tests := []struct {
		name string
		code string
		want user.Currency
	}{
		{
			name: "USD currency",
			code: "USD",
			want: user.USD,
		},
		{
			name: "EUR currency",
			code: "EUR",
			want: user.EUR,
		},
		{
			name: "GBP currency",
			code: "GBP",
			want: user.GBP,
		},
		{
			name: "CAD currency",
			code: "CAD",
			want: user.CAD,
		},
		{
			name: "AUD currency",
			code: "AUD",
			want: user.AUD,
		},
		{
			name: "BRL currency",
			code: "BRL",
			want: user.BRL,
		},
		{
			name: "RMB currency",
			code: "RMB",
			want: user.RMB,
		},
		{
			name: "empty defaults to USD",
			code: "",
			want: user.USD,
		},
		{
			name: "whitespace defaults to USD",
			code: "   ",
			want: user.USD,
		},
		{
			name: "lowercase converted to USD",
			code: "usd",
			want: user.USD,
		},
		{
			name: "mixed case converted to EUR",
			code: "EuR",
			want: user.EUR,
		},
		{
			name: "whitespace trimmed - GBP",
			code: "  gbp  ",
			want: user.GBP,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := user.NewCurrency(tt.code)
			if err != nil {
				t.Fatalf("NewCurrency() unexpected error = %v", err)
			}
			if got != tt.want {
				t.Errorf("NewCurrency() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestNewCurrency_InvalidCases(t *testing.T) {
	tests := []struct {
		name string
		code string
	}{
		{
			name: "invalid - JPY (not supported)",
			code: "JPY",
		},
		{
			name: "invalid - CNY (not supported)",
			code: "CNY",
		},
		{
			name: "invalid - MXN (not supported)",
			code: "MXN",
		},
		{
			name: "invalid - INR (not supported)",
			code: "INR",
		},
		{
			name: "invalid - CHF (not supported)",
			code: "CHF",
		},
		{
			name: "invalid - too short",
			code: "US",
		},
		{
			name: "invalid - too long",
			code: "USDD",
		},
		{
			name: "invalid - single character",
			code: "$",
		},
		{
			name: "invalid - random string",
			code: "ABC",
		},
		{
			name: "invalid - with numbers",
			code: "US1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := user.NewCurrency(tt.code)
			if err == nil {
				t.Fatalf("expected error, got currency = %v", got)
			}
			if err != user.ErrInvalidCurrency {
				t.Errorf("expected ErrInvalidCurrency, got %v", err)
			}
		})
	}
}

func TestCurrency_StringValue(t *testing.T) {
	tests := []struct {
		name     string
		currency user.Currency
		want     string
	}{
		{
			name:     "USD string",
			currency: user.USD,
			want:     "USD",
		},
		{
			name:     "EUR string",
			currency: user.EUR,
			want:     "EUR",
		},
		{
			name:     "GBP string",
			currency: user.GBP,
			want:     "GBP",
		},
		{
			name:     "CAD string",
			currency: user.CAD,
			want:     "CAD",
		},
		{
			name:     "AUD string",
			currency: user.AUD,
			want:     "AUD",
		},
		{
			name:     "BRL string",
			currency: user.BRL,
			want:     "BRL",
		},
		{
			name:     "RMB string",
			currency: user.RMB,
			want:     "RMB",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.currency) != tt.want {
				t.Errorf("Currency string = %q, want %q", string(tt.currency), tt.want)
			}
		})
	}
}

func TestCurrency_AllSupportedCurrencies(t *testing.T) {
	// This test ensures all constants are valid
	currencies := []user.Currency{
		user.USD,
		user.EUR,
		user.GBP,
		user.CAD,
		user.AUD,
		user.BRL,
		user.RMB,
	}

	for _, curr := range currencies {
		t.Run(string(curr), func(t *testing.T) {
			// Try to create currency from the constant's string value
			got, err := user.NewCurrency(string(curr))
			if err != nil {
				t.Errorf("NewCurrency(%q) unexpected error = %v", curr, err)
			}
			if got != curr {
				t.Errorf("NewCurrency(%q) = %q, want %q", curr, got, curr)
			}
		})
	}
}
