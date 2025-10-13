package users_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/cheezecakee/fitrkr-athena/internal/core/domain/user"
	"github.com/cheezecakee/fitrkr-athena/internal/core/services/users"
)

func TestUpdatePlan(t *testing.T) {
	ctx := context.Background()
	testUserID := "test-user-id"

	tests := []struct {
		name          string
		req           users.UpdatePlanReq
		setupMock     func(*MockUserRepo)
		expectedErr   error
		shouldSucceed bool
	}{
		{
			name: "success - upgrade from Basic to Premium monthly",
			req: users.UpdatePlanReq{
				UserID:        testUserID,
				Plan:          "premium",
				BillingPeriod: "monthly",
			},
			setupMock: func(m *MockUserRepo) {
				sub := &user.Subscription{
					Plan:      user.Basic,
					StartedAt: time.Now(),
					AutoRenew: false,
				}
				m.On("GetSubscriptionByID", ctx, testUserID).Return(sub, nil)
				m.On("UpdateSubscription", ctx, mock.MatchedBy(func(s user.Subscription) bool {
					return s.Plan == user.Premium && s.AutoRenew == true
				}), testUserID).Return(nil)
			},
			shouldSucceed: true,
		},
		{
			name: "success - upgrade from Basic to Premium yearly",
			req: users.UpdatePlanReq{
				UserID:        testUserID,
				Plan:          "premium",
				BillingPeriod: "yearly",
			},
			setupMock: func(m *MockUserRepo) {
				sub := &user.Subscription{
					Plan:      user.Basic,
					StartedAt: time.Now(),
					AutoRenew: false,
				}
				m.On("GetSubscriptionByID", ctx, testUserID).Return(sub, nil)
				m.On("UpdateSubscription", ctx, mock.MatchedBy(func(s user.Subscription) bool {
					return s.Plan == user.Premium && *s.BillingPeriod == user.Yearly
				}), testUserID).Return(nil)
			},
			shouldSucceed: true,
		},
		{
			name: "error - invalid plan",
			req: users.UpdatePlanReq{
				UserID:        testUserID,
				Plan:          "invalid",
				BillingPeriod: "monthly",
			},
			setupMock: func(m *MockUserRepo) {
				sub := &user.Subscription{
					Plan:      user.Basic,
					StartedAt: time.Now(),
				}
				m.On("GetSubscriptionByID", ctx, testUserID).Return(sub, nil)
			},
			expectedErr: errors.New("failed to create a new plan: invalid plan"),
		},
		{
			name: "error - invalid billing period",
			req: users.UpdatePlanReq{
				UserID:        testUserID,
				Plan:          "premium",
				BillingPeriod: "invalid",
			},
			setupMock: func(m *MockUserRepo) {
				sub := &user.Subscription{
					Plan:      user.Basic,
					StartedAt: time.Now(),
				}
				m.On("GetSubscriptionByID", ctx, testUserID).Return(sub, nil)
			},
			expectedErr: errors.New("failed to create a new billing period: invalid billing period"),
		},
		{
			name: "error - already on premium",
			req: users.UpdatePlanReq{
				UserID:        testUserID,
				Plan:          "premium",
				BillingPeriod: "monthly",
			},
			setupMock: func(m *MockUserRepo) {
				sub := &user.Subscription{
					Plan:      user.Premium,
					StartedAt: time.Now(),
				}
				m.On("GetSubscriptionByID", ctx, testUserID).Return(sub, nil)
			},
			expectedErr: errors.New("failed to update plan: upgrade not available: must be on Basic plan"),
		},
		{
			name: "error - get subscription fails",
			req: users.UpdatePlanReq{
				UserID:        testUserID,
				Plan:          "premium",
				BillingPeriod: "monthly",
			},
			setupMock: func(m *MockUserRepo) {
				m.On("GetSubscriptionByID", ctx, testUserID).Return(nil, errors.New("db error"))
			},
			expectedErr: errors.New("failed to get subscription: db error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockUserRepo)
			tt.setupMock(mockRepo)
			svc := users.NewService(mockRepo)

			err := svc.UpdatePlan(ctx, tt.req)

			if tt.shouldSucceed {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.EqualError(t, err, tt.expectedErr.Error())
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestRecordPayment(t *testing.T) {
	ctx := context.Background()
	testUserID := "test-user-id"

	tests := []struct {
		name          string
		req           users.RecordPaymentReq
		setupMock     func(*MockUserRepo)
		expectedErr   error
		shouldSucceed bool
	}{
		{
			name: "success - record payment",
			req: users.RecordPaymentReq{
				UserID:   testUserID,
				Amount:   99.99,
				Currency: "USD",
			},
			setupMock: func(m *MockUserRepo) {
				period := user.Monthly
				sub := &user.Subscription{
					Plan:          user.Premium,
					BillingPeriod: &period,
					StartedAt:     time.Now(),
					ExpiresAt:     ptrTime(time.Now().Add(30 * 24 * time.Hour)),
				}
				m.On("GetSubscriptionByID", ctx, testUserID).Return(sub, nil)
				m.On("UpdateSubscription", ctx, mock.MatchedBy(func(s user.Subscription) bool {
					return s.LastPaymentAmount != nil && *s.LastPaymentAmount == 99.99
				}), testUserID).Return(nil)
			},
			shouldSucceed: true,
		},
		{
			name: "error - invalid currency",
			req: users.RecordPaymentReq{
				UserID:   testUserID,
				Amount:   99.99,
				Currency: "INVALID",
			},
			setupMock: func(m *MockUserRepo) {
				sub := &user.Subscription{
					Plan:      user.Premium,
					StartedAt: time.Now(),
				}
				m.On("GetSubscriptionByID", ctx, testUserID).Return(sub, nil)
			},
			expectedErr: errors.New("failed to create a new currency: invalid currency"),
		},
		{
			name: "error - get subscription fails",
			req: users.RecordPaymentReq{
				UserID:   testUserID,
				Amount:   99.99,
				Currency: "USD",
			},
			setupMock: func(m *MockUserRepo) {
				m.On("GetSubscriptionByID", ctx, testUserID).Return(nil, errors.New("db error"))
			},
			expectedErr: errors.New("failed to get subscription: db error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockUserRepo)
			tt.setupMock(mockRepo)
			svc := users.NewService(mockRepo)

			err := svc.RecordPayment(ctx, tt.req)

			if tt.shouldSucceed {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.EqualError(t, err, tt.expectedErr.Error())
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestCancelSubscription(t *testing.T) {
	ctx := context.Background()
	testUserID := "test-user-id"

	tests := []struct {
		name          string
		req           users.CancelSubscriptionReq
		setupMock     func(*MockUserRepo)
		expectedErr   error
		shouldSucceed bool
	}{
		{
			name: "success - cancel premium subscription",
			req: users.CancelSubscriptionReq{
				UserID: testUserID,
			},
			setupMock: func(m *MockUserRepo) {
				period := user.Monthly
				sub := &user.Subscription{
					Plan:          user.Premium,
					BillingPeriod: &period,
					StartedAt:     time.Now(),
					ExpiresAt:     ptrTime(time.Now().Add(30 * 24 * time.Hour)),
					AutoRenew:     true,
				}
				m.On("GetSubscriptionByID", ctx, testUserID).Return(sub, nil)
				m.On("UpdateSubscription", ctx, mock.MatchedBy(func(s user.Subscription) bool {
					return s.AutoRenew == false && s.CancelledAt != nil
				}), testUserID).Return(nil)
			},
			shouldSucceed: true,
		},
		{
			name: "error - cannot cancel basic plan",
			req: users.CancelSubscriptionReq{
				UserID: testUserID,
			},
			setupMock: func(m *MockUserRepo) {
				sub := &user.Subscription{
					Plan:      user.Basic,
					StartedAt: time.Now(),
				}
				m.On("GetSubscriptionByID", ctx, testUserID).Return(sub, nil)
			},
			expectedErr: errors.New("failed to cancel subscription: already on Basic plan"),
		},
		{
			name: "error - get subscription fails",
			req: users.CancelSubscriptionReq{
				UserID: testUserID,
			},
			setupMock: func(m *MockUserRepo) {
				m.On("GetSubscriptionByID", ctx, testUserID).Return(nil, errors.New("db error"))
			},
			expectedErr: errors.New("failed to get subscription: db error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockUserRepo)
			tt.setupMock(mockRepo)
			svc := users.NewService(mockRepo)

			err := svc.CancelSubscription(ctx, tt.req)

			if tt.shouldSucceed {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.EqualError(t, err, tt.expectedErr.Error())
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestStartTrial(t *testing.T) {
	ctx := context.Background()
	testUserID := "test-user-id"

	tests := []struct {
		name          string
		req           users.StartTrialReq
		setupMock     func(*MockUserRepo)
		expectedErr   error
		shouldSucceed bool
	}{
		{
			name: "success - start trial",
			req: users.StartTrialReq{
				UserID: testUserID,
			},
			setupMock: func(m *MockUserRepo) {
				sub := &user.Subscription{
					Plan:      user.Basic,
					StartedAt: time.Now(),
				}
				m.On("GetSubscriptionByID", ctx, testUserID).Return(sub, nil)
				m.On("UpdateSubscription", ctx, mock.MatchedBy(func(s user.Subscription) bool {
					return s.TrialEndsAt != nil
				}), testUserID).Return(nil)
			},
			shouldSucceed: true,
		},
		{
			name: "error - get subscription fails",
			req: users.StartTrialReq{
				UserID: testUserID,
			},
			setupMock: func(m *MockUserRepo) {
				m.On("GetSubscriptionByID", ctx, testUserID).Return(nil, errors.New("db error"))
			},
			expectedErr: errors.New("failed to get subscription: db error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockUserRepo)
			tt.setupMock(mockRepo)
			svc := users.NewService(mockRepo)

			err := svc.StartTrial(ctx, tt.req)

			if tt.shouldSucceed {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.EqualError(t, err, tt.expectedErr.Error())
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

// Helper function for pointer to time.Time
func ptrTime(t time.Time) *time.Time {
	return &t
}
