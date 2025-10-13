package user

import (
	"errors"
	"time"
)

var (
	ErrUpgradeNotAvailable    = errors.New("upgrade not available: must be on Basic plan")
	ErrInvalidUpgradeTarget   = errors.New("can only upgrade to Premium plan")
	ErrDowngradeNotAvailable  = errors.New("downgrade not available: must be on Premium plan")
	ErrInvalidDowngradeTarget = errors.New("can only downgrade to Basic plan")
	ErrAlreadyOnBasic         = errors.New("already on Basic plan")
)

const (
	MonthlyDuration = 30 * 24 * time.Hour  // 30 days
	YearlyDuration  = 365 * 24 * time.Hour // 365 days
	TrialDays       = 14
)

type Subscription struct {
	Plan          Plan
	BillingPeriod *Period
	StartedAt     time.Time
	ExpiresAt     *time.Time
	AutoRenew     bool
	CancelledAt   *time.Time

	LastPaymentAt       *time.Time
	LastPaymentAmount   *float64
	LastPaymentCurrency *Currency

	TrialEndsAt *time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func NewSubscription() Subscription {
	return Subscription{
		Plan:      Basic, // Free plan
		StartedAt: time.Now(),
		AutoRenew: false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func (s *Subscription) Upgrade(plan Plan, period Period) error {
	if s.Plan != Basic {
		return ErrUpgradeNotAvailable
	}
	if plan != Premium {
		return ErrInvalidUpgradeTarget
	}

	s.Plan = plan
	s.BillingPeriod = &period
	s.AutoRenew = true

	expiry := time.Now()
	if period == Monthly {
		expiry = expiry.Add(MonthlyDuration) // Monthly
	} else {
		expiry = expiry.Add(YearlyDuration) // Yearly
	}
	s.ExpiresAt = &expiry

	s.Touch()
	return nil
}

func (s *Subscription) ProcessExpiry() {
	if s.ExpiresAt == nil {
		return // Basic plan, no expiry
	}
	now := time.Now()

	if now.Before(*s.ExpiresAt) {
		return // Not expired yet
	}

	if s.AutoRenew {
		s.Renew() // Attemp renewal
	} else {
		// Downgrade to Basic after expiry
		s.Plan = Basic
		s.BillingPeriod = nil
		s.ExpiresAt = nil
		s.Touch()
	}
}

func (s *Subscription) ProcessPayment(amount float64, currency Currency) {
	now := time.Now()
	s.LastPaymentAt = &now
	s.LastPaymentAmount = &amount
	s.LastPaymentCurrency = &currency
	s.UpdatedAt = now

	if s.BillingPeriod != nil {
		s.ExtendExpiry(s.getNextBillingDuration())
	}
}

func (s *Subscription) ExtendExpiry(duration time.Duration) {
	if s.ExpiresAt == nil {
		return // Can't extend Basic plan
	}

	extended := s.ExpiresAt.Add(duration)
	s.ExpiresAt = &extended
	s.Touch()
}

func (s *Subscription) Cancel() error {
	if s.Plan == Basic {
		return ErrAlreadyOnBasic
	}

	// Don't change plan yet - keep Premium until expiry
	s.AutoRenew = false
	now := time.Now()
	s.CancelledAt = &now
	s.Touch()

	return nil
}

func (s *Subscription) Renew() {
	if s.BillingPeriod == nil {
		return
	}

	s.ExtendExpiry(s.getNextBillingDuration())
}

func (s *Subscription) StartTrial() {
	now := time.Now()
	hours := time.Duration(TrialDays) * 24
	trialPeriod := now.Add(hours * time.Hour)

	s.TrialEndsAt = &trialPeriod
}

func (s *Subscription) DaysUntilExpiry() int {
	if s.ExpiresAt == nil {
		return 0
	}

	duration := time.Until(*s.ExpiresAt)
	return int(duration.Hours() / 24)
}

func (s *Subscription) TimeRemaining() time.Duration {
	if s.ExpiresAt == nil {
		return 0
	}

	return time.Until(*s.ExpiresAt)
}

func (s *Subscription) HasExpired() bool {
	if s.ExpiresAt == nil {
		return false
	}

	return time.Now().After(*s.ExpiresAt)
}

func (s *Subscription) getNextBillingDuration() time.Duration {
	if s.BillingPeriod == nil {
		return 0
	}

	switch *s.BillingPeriod {
	case Monthly:
		return MonthlyDuration
	case Yearly:
		return YearlyDuration
	default:
		return 0
	}
}

func (s *Subscription) Touch() {
	s.UpdatedAt = time.Now()
}
