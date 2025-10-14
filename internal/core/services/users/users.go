// Package users
package users

import (
	"context"
	"errors"

	"github.com/cheezecakee/fitrkr-athena/internal/core/ports"
)

var (
	ErrUserNotFound      = errors.New("user does not exist")
	ErrDuplicateEmail    = errors.New("email already exists")
	ErrDuplicateUsername = errors.New("username already exists")
)

type UserService interface {
	CreateAccount(ctx context.Context, req CreateAccountReq) (*CreateAccountResp, error)
	GetByID(ctx context.Context, req GetUserByIDReq) (*GetUserResp, error)
	GetByUsername(ctx context.Context, req GetUserByUsernameReq) (*GetUserResp, error)
	GetByEmail(ctx context.Context, req GetUserByEmailReq) (*GetUserResp, error)
	Update(ctx context.Context, req UpdateUserReq) error
	Delete(ctx context.Context, req DeleteAccountReq) error

	GetStats(ctx context.Context, req GetStatsReq) (*GetStatsResp, error)
	GetSubscription(ctx context.Context, req GetSubscriptionReq) (*GetSubscriptionResp, error)
	GetSettings(ctx context.Context, req GetSettingsReq) (*GetSettingsResp, error)

	UpgradePlan(ctx context.Context, req UpgradePlanReq) error
	RecordPayment(ctx context.Context, req RecordPaymentReq) error
	CancelSubscription(ctx context.Context, req CancelSubscriptionReq) error
	StartTrial(ctx context.Context, req StartTrialReq) error

	UpdateSettings(ctx context.Context, req UpdateSettingsReq) error

	UpdateBodyMetrics(ctx context.Context, req UpdateBodyMetricsReq) error
}

type Service struct {
	userRepo ports.UserRepo
}

func NewService(userRepo ports.UserRepo) *Service {
	return &Service{
		userRepo: userRepo,
	}
}
