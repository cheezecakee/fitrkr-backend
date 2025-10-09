package user_test

import (
	"testing"
	"time"

	"github.com/cheezecakee/fitrkr-athena/internal/core/domain/user"
)

func TestNewSubscription(t *testing.T) {
	sub := user.NewSubscription()

	// Verify Basic plan is set
	if sub.Plan != user.Basic {
		t.Errorf("expected Plan to be Basic, got %v", sub.Plan)
	}

	// Verify optional fields are nil for Basic plan
	if sub.BillingPeriod != nil {
		t.Error("expected BillingPeriod to be nil for Basic plan")
	}
	if sub.ExpiresAt != nil {
		t.Error("expected ExpiresAt to be nil for Basic plan")
	}
	if sub.CancelledAt != nil {
		t.Error("expected CancelledAt to be nil")
	}
	if sub.LastPaymentAt != nil {
		t.Error("expected LastPaymentAt to be nil")
	}
	if sub.LastPaymentAmount != nil {
		t.Error("expected LastPaymentAmount to be nil")
	}
	if sub.LastPaymentCurrency != nil {
		t.Error("expected LastPaymentCurrency to be nil")
	}
	if sub.TrialEndsAt != nil {
		t.Error("expected TrialEndsAt to be nil")
	}

	// Verify defaults
	if sub.AutoRenew {
		t.Error("expected AutoRenew to be false for Basic plan")
	}

	// Verify timestamps are set
	if sub.StartedAt.IsZero() {
		t.Error("expected StartedAt to be set")
	}
	if sub.CreatedAt.IsZero() {
		t.Error("expected CreatedAt to be set")
	}
	if sub.UpdatedAt.IsZero() {
		t.Error("expected UpdatedAt to be set")
	}
}

func TestSubscription_Upgrade(t *testing.T) {
	tests := []struct {
		name        string
		initialPlan user.Plan
		targetPlan  user.Plan
		period      user.Period
		wantErr     error
	}{
		{
			name:        "successful upgrade from Basic to Premium monthly",
			initialPlan: user.Basic,
			targetPlan:  user.Premium,
			period:      user.Monthly,
			wantErr:     nil,
		},
		{
			name:        "successful upgrade from Basic to Premium yearly",
			initialPlan: user.Basic,
			targetPlan:  user.Premium,
			period:      user.Yearly,
			wantErr:     nil,
		},
		{
			name:        "cannot upgrade from Premium",
			initialPlan: user.Premium,
			targetPlan:  user.Premium,
			period:      user.Monthly,
			wantErr:     user.ErrUpgradeNotAvailable,
		},
		{
			name:        "cannot upgrade to Basic",
			initialPlan: user.Basic,
			targetPlan:  user.Basic,
			period:      user.Monthly,
			wantErr:     user.ErrInvalidUpgradeTarget,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sub := user.NewSubscription()
			sub.Plan = tt.initialPlan

			err := sub.Upgrade(tt.targetPlan, tt.period)

			if err != tt.wantErr {
				t.Errorf("Upgrade() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr == nil {
				// Verify successful upgrade
				if sub.Plan != user.Premium {
					t.Errorf("expected Plan to be Premium, got %v", sub.Plan)
				}
				if sub.BillingPeriod == nil {
					t.Fatal("expected BillingPeriod to be set")
				}
				if *sub.BillingPeriod != tt.period {
					t.Errorf("expected BillingPeriod %v, got %v", tt.period, *sub.BillingPeriod)
				}
				if sub.ExpiresAt == nil {
					t.Fatal("expected ExpiresAt to be set")
				}
				if !sub.AutoRenew {
					t.Error("expected AutoRenew to be true after upgrade")
				}

				// Verify expiry is set correctly based on period
				expectedDuration := user.MonthlyDuration
				if tt.period == user.Yearly {
					expectedDuration = user.YearlyDuration
				}

				timeDiff := time.Until(*sub.ExpiresAt)
				tolerance := 5 * time.Second
				if timeDiff < expectedDuration-tolerance || timeDiff > expectedDuration+tolerance {
					t.Errorf("expected expiry around %v from now, got %v", expectedDuration, timeDiff)
				}
			}
		})
	}
}

func TestSubscription_Cancel(t *testing.T) {
	tests := []struct {
		name    string
		setup   func() user.Subscription
		wantErr error
	}{
		{
			name: "cancel Premium subscription",
			setup: func() user.Subscription {
				sub := user.NewSubscription()
				sub.Upgrade(user.Premium, user.Monthly)
				return sub
			},
			wantErr: nil,
		},
		{
			name: "cannot cancel Basic plan",
			setup: func() user.Subscription {
				return user.NewSubscription()
			},
			wantErr: user.ErrAlreadyOnBasic,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sub := tt.setup()
			err := sub.Cancel()

			if err != tt.wantErr {
				t.Errorf("Cancel() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr == nil {
				// Verify cancellation
				if sub.AutoRenew {
					t.Error("expected AutoRenew to be false after cancel")
				}
				if sub.CancelledAt == nil {
					t.Error("expected CancelledAt to be set")
				}
				// Plan should still be Premium until expiry
				if sub.Plan != user.Premium {
					t.Error("expected Plan to remain Premium until expiry")
				}
			}
		})
	}
}

func TestSubscription_ProcessPayment(t *testing.T) {
	sub := user.NewSubscription()
	sub.Upgrade(user.Premium, user.Monthly)

	initialExpiry := *sub.ExpiresAt
	amount := 9.99
	currency := user.USD

	sub.ProcessPayment(amount, currency)

	// Verify payment details are recorded
	if sub.LastPaymentAt == nil {
		t.Fatal("expected LastPaymentAt to be set")
	}
	if sub.LastPaymentAmount == nil {
		t.Fatal("expected LastPaymentAmount to be set")
	}
	if *sub.LastPaymentAmount != amount {
		t.Errorf("expected LastPaymentAmount %v, got %v", amount, *sub.LastPaymentAmount)
	}
	if sub.LastPaymentCurrency == nil {
		t.Fatal("expected LastPaymentCurrency to be set")
	}
	if *sub.LastPaymentCurrency != currency {
		t.Errorf("expected LastPaymentCurrency %v, got %v", currency, *sub.LastPaymentCurrency)
	}

	// Verify expiry was extended
	if sub.ExpiresAt.Before(initialExpiry) || sub.ExpiresAt.Equal(initialExpiry) {
		t.Error("expected ExpiresAt to be extended after payment")
	}
}

func TestSubscription_ProcessExpiry(t *testing.T) {
	tests := []struct {
		name     string
		setup    func() user.Subscription
		wantPlan user.Plan
	}{
		{
			name: "Basic plan - no expiry processing",
			setup: func() user.Subscription {
				return user.NewSubscription()
			},
			wantPlan: user.Basic,
		},
		{
			name: "Premium not expired - no change",
			setup: func() user.Subscription {
				sub := user.NewSubscription()
				sub.Upgrade(user.Premium, user.Monthly)
				return sub
			},
			wantPlan: user.Premium,
		},
		{
			name: "Premium expired with AutoRenew - should renew",
			setup: func() user.Subscription {
				sub := user.NewSubscription()
				sub.Upgrade(user.Premium, user.Monthly)
				// Set expiry to past
				past := time.Now().Add(-24 * time.Hour)
				sub.ExpiresAt = &past
				sub.AutoRenew = true
				return sub
			},
			wantPlan: user.Premium,
		},
		{
			name: "Premium expired without AutoRenew - downgrade to Basic",
			setup: func() user.Subscription {
				sub := user.NewSubscription()
				sub.Upgrade(user.Premium, user.Monthly)
				// Set expiry to past
				past := time.Now().Add(-24 * time.Hour)
				sub.ExpiresAt = &past
				sub.AutoRenew = false
				return sub
			},
			wantPlan: user.Basic,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sub := tt.setup()
			sub.ProcessExpiry()

			if sub.Plan != tt.wantPlan {
				t.Errorf("expected Plan %v, got %v", tt.wantPlan, sub.Plan)
			}

			// If downgraded to Basic, verify cleanup
			if tt.wantPlan == user.Basic && sub.BillingPeriod == nil {
				if sub.BillingPeriod != nil {
					t.Error("expected BillingPeriod to be nil after downgrade")
				}
				if sub.ExpiresAt != nil {
					t.Error("expected ExpiresAt to be nil after downgrade")
				}
			}
		})
	}
}

func TestSubscription_StartTrial(t *testing.T) {
	sub := user.NewSubscription()
	trialDays := 7

	sub.StartTrial(trialDays)

	if sub.TrialEndsAt == nil {
		t.Fatal("expected TrialEndsAt to be set")
	}

	expectedEnd := time.Now().Add(time.Duration(trialDays) * 24 * time.Hour)
	tolerance := 5 * time.Second

	if sub.TrialEndsAt.Before(expectedEnd.Add(-tolerance)) || sub.TrialEndsAt.After(expectedEnd.Add(tolerance)) {
		t.Errorf("expected TrialEndsAt around %v, got %v", expectedEnd, sub.TrialEndsAt)
	}
}

func TestSubscription_DaysUntilExpiry(t *testing.T) {
	tests := []struct {
		name     string
		setup    func() user.Subscription
		wantDays int
	}{
		{
			name: "Basic plan - no expiry",
			setup: func() user.Subscription {
				return user.NewSubscription()
			},
			wantDays: 0,
		},
		{
			name: "Premium with 30 days remaining",
			setup: func() user.Subscription {
				sub := user.NewSubscription()
				sub.Upgrade(user.Premium, user.Monthly)
				return sub
			},
			wantDays: 30,
		},
		{
			name: "Premium expired",
			setup: func() user.Subscription {
				sub := user.NewSubscription()
				sub.Upgrade(user.Premium, user.Monthly)
				past := time.Now().Add(-24 * time.Hour)
				sub.ExpiresAt = &past
				return sub
			},
			wantDays: -1, // negative = expired
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sub := tt.setup()
			days := sub.DaysUntilExpiry()

			// Allow 1 day tolerance for timing
			tolerance := 1
			if days < tt.wantDays-tolerance || days > tt.wantDays+tolerance {
				t.Errorf("expected around %d days, got %d", tt.wantDays, days)
			}
		})
	}
}

func TestSubscription_HasExpired(t *testing.T) {
	tests := []struct {
		name    string
		setup   func() user.Subscription
		wantExp bool
	}{
		{
			name: "Basic plan never expires",
			setup: func() user.Subscription {
				return user.NewSubscription()
			},
			wantExp: false,
		},
		{
			name: "Premium not expired",
			setup: func() user.Subscription {
				sub := user.NewSubscription()
				sub.Upgrade(user.Premium, user.Monthly)
				return sub
			},
			wantExp: false,
		},
		{
			name: "Premium expired",
			setup: func() user.Subscription {
				sub := user.NewSubscription()
				sub.Upgrade(user.Premium, user.Monthly)
				past := time.Now().Add(-24 * time.Hour)
				sub.ExpiresAt = &past
				return sub
			},
			wantExp: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sub := tt.setup()
			if sub.HasExpired() != tt.wantExp {
				t.Errorf("HasExpired() = %v, want %v", sub.HasExpired(), tt.wantExp)
			}
		})
	}
}

func TestSubscription_ExtendExpiry(t *testing.T) {
	t.Run("cannot extend Basic plan", func(t *testing.T) {
		sub := user.NewSubscription()
		sub.ExtendExpiry(30 * 24 * time.Hour)

		if sub.ExpiresAt != nil {
			t.Error("expected ExpiresAt to remain nil for Basic plan")
		}
	})

	t.Run("extend Premium subscription", func(t *testing.T) {
		sub := user.NewSubscription()
		sub.Upgrade(user.Premium, user.Monthly)

		initialExpiry := *sub.ExpiresAt
		extension := 30 * 24 * time.Hour

		sub.ExtendExpiry(extension)

		expectedExpiry := initialExpiry.Add(extension)
		tolerance := 5 * time.Second

		if sub.ExpiresAt.Before(expectedExpiry.Add(-tolerance)) || sub.ExpiresAt.After(expectedExpiry.Add(tolerance)) {
			t.Errorf("expected expiry around %v, got %v", expectedExpiry, sub.ExpiresAt)
		}
	})
}

func TestSubscription_Renew(t *testing.T) {
	tests := []struct {
		name   string
		period user.Period
	}{
		{
			name:   "renew monthly subscription",
			period: user.Monthly,
		},
		{
			name:   "renew yearly subscription",
			period: user.Yearly,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sub := user.NewSubscription()
			sub.Upgrade(user.Premium, tt.period)

			initialExpiry := *sub.ExpiresAt

			sub.Renew()

			// Verify expiry was extended
			if !sub.ExpiresAt.After(initialExpiry) {
				t.Error("expected ExpiresAt to be extended after renewal")
			}

			expectedDuration := user.MonthlyDuration
			if tt.period == user.Yearly {
				expectedDuration = user.YearlyDuration
			}

			actualExtension := sub.ExpiresAt.Sub(initialExpiry)
			tolerance := 5 * time.Second

			if actualExtension < expectedDuration-tolerance || actualExtension > expectedDuration+tolerance {
				t.Errorf("expected extension of %v, got %v", expectedDuration, actualExtension)
			}
		})
	}
}

func TestSubscription_TimeRemaining(t *testing.T) {
	t.Run("Basic plan has no time remaining", func(t *testing.T) {
		sub := user.NewSubscription()
		if sub.TimeRemaining() != 0 {
			t.Errorf("expected TimeRemaining to be 0 for Basic, got %v", sub.TimeRemaining())
		}
	})

	t.Run("Premium subscription has time remaining", func(t *testing.T) {
		sub := user.NewSubscription()
		sub.Upgrade(user.Premium, user.Monthly)

		remaining := sub.TimeRemaining()
		if remaining <= 0 {
			t.Error("expected positive time remaining for Premium")
		}

		expectedRemaining := user.MonthlyDuration
		tolerance := 10 * time.Second

		if remaining < expectedRemaining-tolerance || remaining > expectedRemaining+tolerance {
			t.Errorf("expected around %v remaining, got %v", expectedRemaining, remaining)
		}
	})
}

func TestSubscription_OptionalFieldsNilForBasic(t *testing.T) {
	sub := user.NewSubscription()

	// All optional fields should be nil for Basic plan
	if sub.BillingPeriod != nil {
		t.Error("BillingPeriod should be nil for Basic")
	}
	if sub.ExpiresAt != nil {
		t.Error("ExpiresAt should be nil for Basic")
	}
	if sub.CancelledAt != nil {
		t.Error("CancelledAt should be nil initially")
	}
	if sub.LastPaymentAt != nil {
		t.Error("LastPaymentAt should be nil initially")
	}
	if sub.LastPaymentAmount != nil {
		t.Error("LastPaymentAmount should be nil initially")
	}
	if sub.LastPaymentCurrency != nil {
		t.Error("LastPaymentCurrency should be nil initially")
	}
	if sub.TrialEndsAt != nil {
		t.Error("TrialEndsAt should be nil initially")
	}
}

func TestSubscription_OptionalFieldsSetForPremium(t *testing.T) {
	sub := user.NewSubscription()
	sub.Upgrade(user.Premium, user.Monthly)

	// Required fields for Premium should be set
	if sub.BillingPeriod == nil {
		t.Error("BillingPeriod should be set for Premium")
	}
	if sub.ExpiresAt == nil {
		t.Error("ExpiresAt should be set for Premium")
	}

	// Payment fields set after payment
	sub.ProcessPayment(9.99, user.USD)

	if sub.LastPaymentAt == nil {
		t.Error("LastPaymentAt should be set after payment")
	}
	if sub.LastPaymentAmount == nil {
		t.Error("LastPaymentAmount should be set after payment")
	}
	if sub.LastPaymentCurrency == nil {
		t.Error("LastPaymentCurrency should be set after payment")
	}
}
