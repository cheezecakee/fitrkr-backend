package users

import (
	"context"
	"fmt"

	"github.com/cheezecakee/logr"

	"github.com/cheezecakee/fitrkr-athena/internal/core/domain/user"
)

type UpgradePlanReq struct {
	UserID        string `json:"user_id"`
	Plan          string `json:"plan"`
	BillingPeriod string `json:"billing_period"`
}

func (s *Service) UpgradePlan(ctx context.Context, req UpgradePlanReq) error {
	existing, err := s.userRepo.GetSubscriptionByID(ctx, req.UserID)
	if err != nil {
		logr.Get().Errorf("failed to get subscription: %v", err)
		return fmt.Errorf("failed to get subscription: %w", err)
	}

	plan, err := user.NewPlan(req.Plan)
	if err != nil {
		logr.Get().Errorf("failed to create a new plan: %v", err)
		return fmt.Errorf("failed to create a new plan: %w", err)
	}
	period, err := user.NewBillingPeriod(req.BillingPeriod)
	if err != nil {
		logr.Get().Errorf("failed to create a new billing period: %v", err)
		return fmt.Errorf("failed to create a new billing period: %w", err)
	}

	err = existing.Upgrade(plan, period)
	if err != nil {
		logr.Get().Errorf("failed to update plan: %v", err)
		return fmt.Errorf("failed to update plan: %w", err)
	}

	err = s.userRepo.UpdateSubscription(ctx, *existing, req.UserID)
	if err != nil {
		logr.Get().Errorf("failed to update subscription: %v", err)
		return fmt.Errorf("failed to update subscription: %w", err)
	}

	return nil
}

type RecordPaymentReq struct {
	UserID   string  `json:"user_id"`
	Amount   float64 `json:"amount"`
	Currency string  `json:"currency"`
}

func (s *Service) RecordPayment(ctx context.Context, req RecordPaymentReq) error {
	existing, err := s.userRepo.GetSubscriptionByID(ctx, req.UserID)
	if err != nil {
		logr.Get().Errorf("failed to get subscription: %v", err)
		return fmt.Errorf("failed to get subscription: %w", err)
	}

	currency, err := user.NewCurrency(req.Currency)
	if err != nil {
		logr.Get().Errorf("failed to create a new currency: %v", err)
		return fmt.Errorf("failed to create a new currency: %w", err)
	}

	existing.ProcessPayment(req.Amount, currency)

	err = s.userRepo.UpdateSubscription(ctx, *existing, req.UserID)
	if err != nil {
		logr.Get().Errorf("failed to update subscription: %v", err)
		return fmt.Errorf("failed to update subscription: %w", err)
	}
	return nil
}

type CancelSubscriptionReq struct {
	UserID string `json:"user_id"`
}

func (s *Service) CancelSubscription(ctx context.Context, req CancelSubscriptionReq) error {
	existing, err := s.userRepo.GetSubscriptionByID(ctx, req.UserID)
	if err != nil {
		logr.Get().Errorf("failed to get subscription: %v", err)
		return fmt.Errorf("failed to get subscription: %w", err)
	}

	err = existing.Cancel()
	if err != nil {
		logr.Get().Errorf("failed to cancel subscription: %v", err)
		return fmt.Errorf("failed to cancel subscription: %w", err)
	}

	err = s.userRepo.UpdateSubscription(ctx, *existing, req.UserID)
	if err != nil {
		logr.Get().Errorf("failed to update subscription: %v", err)
		return fmt.Errorf("failed to update subscription: %w", err)
	}
	return nil
}

type StartTrialReq struct {
	UserID string `json:"user_id"`
}

func (s *Service) StartTrial(ctx context.Context, req StartTrialReq) error {
	existing, err := s.userRepo.GetSubscriptionByID(ctx, req.UserID)
	if err != nil {
		logr.Get().Errorf("failed to get subscription: %v", err)
		return fmt.Errorf("failed to get subscription: %w", err)
	}

	existing.StartTrial()

	err = s.userRepo.UpdateSubscription(ctx, *existing, req.UserID)
	if err != nil {
		logr.Get().Errorf("failed to update subscription: %v", err)
		return fmt.Errorf("failed to update subscription: %w", err)
	}
	return nil
}
