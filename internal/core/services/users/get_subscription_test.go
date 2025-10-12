package users_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/cheezecakee/fitrkr-athena/internal/core/domain/user"
	"github.com/cheezecakee/fitrkr-athena/internal/core/ports"
	"github.com/cheezecakee/fitrkr-athena/internal/core/services/users"
)

func TestGetSubscription(t *testing.T) {
	ctx := context.Background()
	testID := "test-user-id"

	tests := []struct {
		name          string
		req           users.GetSubscriptionReq
		setupMock     func(*MockUserRepo)
		expectedErr   error
		shouldSucceed bool
	}{
		{
			name: "success - subscription retrieved",
			req: users.GetSubscriptionReq{
				ID: testID,
			},
			setupMock: func(m *MockUserRepo) {
				period := user.Monthly
				sub := &ports.Subscription{
					Plan:          user.Premium,
					BillingPeriod: &period,
					StartedAt:     time.Now(),
					ExpiresAT:     nil,
					AutoRenew:     true,
					CancelledAt:   nil,
				}
				m.On("GetSubscriptionByID", ctx, testID).Return(sub, nil)
			},
			shouldSucceed: true,
		},
		{
			name: "error - repo fails",
			req: users.GetSubscriptionReq{
				ID: testID,
			},
			setupMock: func(m *MockUserRepo) {
				m.On("GetSubscriptionByID", ctx, testID).Return(nil, errors.New("db error"))
			},
			expectedErr: errors.New("failed to get subscription: db error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockUserRepo)
			tt.setupMock(mockRepo)
			svc := users.NewService(mockRepo)

			resp, err := svc.GetSubscription(ctx, tt.req)

			if tt.shouldSucceed {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
			} else {
				assert.Error(t, err)
				assert.Nil(t, resp)
				assert.EqualError(t, err, tt.expectedErr.Error())
			}

			mockRepo.AssertExpectations(t)
		})
	}
}
