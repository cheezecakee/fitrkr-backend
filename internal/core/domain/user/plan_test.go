package user_test

import (
	"testing"

	"github.com/cheezecakee/fitrkr-athena/internal/core/domain/user"
)

func TestNewPlan(t *testing.T) {
	tests := []struct {
		name    string
		plan    string
		wantErr bool
	}{
		{
			name:    "valid plan - basic",
			plan:    "basic",
			wantErr: false,
		},
		{
			name:    "valid plan - premium",
			plan:    "premium",
			wantErr: false,
		},
		{
			name:    "valid plan - with whitespace trimmed",
			plan:    "  premium  ",
			wantErr: false,
		},
		{
			name:    "valid plan - empty defaults to basic",
			plan:    "",
			wantErr: false,
		},
		{
			name:    "valid plan - uppercase",
			plan:    "BASIC",
			wantErr: false,
		},
		{
			name:    "valid plan - mixed case",
			plan:    "PreMiuM",
			wantErr: false,
		},
		{
			name:    "invalid plan - unknown tier",
			plan:    "pro",
			wantErr: true,
		},
		{
			name:    "invalid plan - with digits",
			plan:    "basic1",
			wantErr: true,
		},
		{
			name:    "invalid plan - with special chars",
			plan:    "premium!",
			wantErr: true,
		},
		{
			name:    "invalid plan - with underscore",
			plan:    "pre_mium",
			wantErr: true,
		},
		{
			name:    "invalid plan - misspelled",
			plan:    "premum",
			wantErr: true,
		},
		{
			name:    "invalid plan - with spaces",
			plan:    "basic plan",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := user.NewPlan(tt.plan)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewPlan() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewPlan_ValidCases(t *testing.T) {
	tests := []struct {
		name string
		plan string
		want user.Plan
	}{
		{
			name: "basic plan",
			plan: "basic",
			want: user.Basic,
		},
		{
			name: "premium plan",
			plan: "premium",
			want: user.Premium,
		},
		{
			name: "empty defaults to basic",
			plan: "",
			want: user.Basic,
		},
		{
			name: "uppercase converted to basic",
			plan: "BASIC",
			want: user.Basic,
		},
		{
			name: "mixed case converted to premium",
			plan: "PrEmIuM",
			want: user.Premium,
		},
		{
			name: "whitespace trimmed - basic",
			plan: "   basic   ",
			want: user.Basic,
		},
		{
			name: "whitespace trimmed - premium",
			plan: "  premium  ",
			want: user.Premium,
		},
		{
			name: "whitespace only defaults to basic",
			plan: "   ",
			want: user.Basic,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := user.NewPlan(tt.plan)
			if err != nil {
				t.Fatalf("NewPlan() unexpected error = %v", err)
			}
			if got != tt.want {
				t.Errorf("NewPlan() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestNewPlan_InvalidCases(t *testing.T) {
	tests := []struct {
		name string
		plan string
	}{
		{
			name: "invalid - enterprise",
			plan: "enterprise",
		},
		{
			name: "invalid - free",
			plan: "free",
		},
		{
			name: "invalid - pro",
			plan: "pro",
		},
		{
			name: "invalid - starter",
			plan: "starter",
		},
		{
			name: "invalid - random string",
			plan: "notaplan",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := user.NewPlan(tt.plan)
			if err == nil {
				t.Fatalf("expected error, got plan = %v", got)
			}
			if err != user.ErrInvalidPlan {
				t.Errorf("expected ErrInvalidPlan, got %v", err)
			}
		})
	}
}

func TestPlan_StringValue(t *testing.T) {
	tests := []struct {
		name string
		plan user.Plan
		want string
	}{
		{
			name: "basic plan string",
			plan: user.Basic,
			want: "basic",
		},
		{
			name: "premium plan string",
			plan: user.Premium,
			want: "premium",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.plan) != tt.want {
				t.Errorf("Plan string = %q, want %q", string(tt.plan), tt.want)
			}
		})
	}
}
