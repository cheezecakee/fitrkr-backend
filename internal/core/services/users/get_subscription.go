package users

import (
	"context"
	"fmt"
	"time"

	"github.com/cheezecakee/logr"

	"github.com/cheezecakee/fitrkr-athena/internal/core/domain/user"
)

type GetSubscriptionReq struct {
	ID string
}

type GetSubscriptionResp struct {
	Plan        user.Plan
	Period      *user.Period
	StartAt     time.Time
	ExpiresAt   *time.Time
	AutoRenew   bool
	CancelledAt *time.Time

	LastPaymentAt       *time.Time
	LastPaymentAmount   *float64
	LastPaymentCurrency *user.Currency

	TrialEndsAt *time.Time
}

func (s *Service) GetSubscription(ctx context.Context, req GetSubscriptionReq) (*GetSubscriptionResp, error) {
	sub, err := s.userRepo.GetSubscriptionByID(ctx, req.ID)
	if err != nil {
		logr.Get().Errorf("failed to get subscription: %v", err)
		return nil, fmt.Errorf("failed to get subscription: %w", err)
	}

	return &GetSubscriptionResp{
		Plan:                sub.Plan,
		Period:              sub.BillingPeriod,
		StartAt:             sub.StartedAt,
		ExpiresAt:           sub.ExpiresAT,
		AutoRenew:           sub.AutoRenew,
		CancelledAt:         sub.CancelledAt,
		LastPaymentAt:       sub.LastPaymentAt,
		LastPaymentAmount:   sub.LastPaymentAmount,
		LastPaymentCurrency: &sub.LastPaymentCurrency,
		TrialEndsAt:         sub.TrialEndsAt,
	}, nil
}
