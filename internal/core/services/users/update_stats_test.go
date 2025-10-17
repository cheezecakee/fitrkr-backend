package users_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/cheezecakee/fitrkr-athena/internal/core/domain/user"
	"github.com/cheezecakee/fitrkr-athena/internal/core/ports"
	"github.com/cheezecakee/fitrkr-athena/internal/core/services/users"
)

func TestUpdateBodyMetrics(t *testing.T) {
	ctx := context.Background()
	testUserID := "test-user-id"

	tests := []struct {
		name          string
		req           users.UpdateBodyMetricsReq
		setupMock     func(*MockUserRepo)
		expectedErr   error
		shouldSucceed bool
	}{
		{
			name: "success - update weight only",
			req: users.UpdateBodyMetricsReq{
				UserID:      testUserID,
				WeightValue: ptrFloat64(75.5),
			},
			setupMock: func(m *MockUserRepo) {
				m.On("GetSettingsByID", ctx, testUserID).Return(
					&user.Settings{WeightUnit: user.Kg, HeightUnit: user.Cm},
					nil,
				)
				m.On("UpdateBodyMetrics", ctx, mock.MatchedBy(func(u ports.UpdateBodyMetrics) bool {
					return u.WeightValue != nil
				}), testUserID).Return(nil)
			},
			shouldSucceed: true,
		},
		{
			name: "success - update all fields",
			req: users.UpdateBodyMetricsReq{
				UserID:      testUserID,
				WeightValue: ptrFloat64(80),
				HeightValue: ptrFloat64(180),
				BFP:         ptrFloat64(20),
			},
			setupMock: func(m *MockUserRepo) {
				m.On("GetSettingsByID", ctx, testUserID).Return(
					&user.Settings{WeightUnit: user.Kg, HeightUnit: user.Cm},
					nil,
				)
				m.On("UpdateBodyMetrics", ctx, mock.MatchedBy(func(u ports.UpdateBodyMetrics) bool {
					return u.WeightValue != nil && u.HeightValue != nil && u.BFP != nil
				}), testUserID).Return(nil)
			},
			shouldSucceed: true,
		},
		{
			name: "error - invalid weight (negative)",
			req: users.UpdateBodyMetricsReq{
				UserID:      testUserID,
				WeightValue: ptrFloat64(-10),
			},
			setupMock: func(m *MockUserRepo) {
				m.On("GetSettingsByID", ctx, testUserID).Return(
					&user.Settings{WeightUnit: user.Kg, HeightUnit: user.Cm},
					nil,
				)
			},
			expectedErr: errors.New("failed to create user weight: weight cannot be negative"),
		},
		{
			name: "error - invalid height (zero)",
			req: users.UpdateBodyMetricsReq{
				UserID:      testUserID,
				HeightValue: ptrFloat64(0),
			},
			setupMock: func(m *MockUserRepo) {
				m.On("GetSettingsByID", ctx, testUserID).Return(
					&user.Settings{WeightUnit: user.Kg, HeightUnit: user.Cm},
					nil,
				)
			},
			expectedErr: errors.New("failed to create user height: height cannot be zero"),
		},
		{
			name: "error - invalid BFP (over 100)",
			req: users.UpdateBodyMetricsReq{
				UserID: testUserID,
				BFP:    ptrFloat64(150),
			},
			setupMock: func(m *MockUserRepo) {
				m.On("GetSettingsByID", ctx, testUserID).Return(
					&user.Settings{WeightUnit: user.Kg, HeightUnit: user.Cm},
					nil,
				)
			},
			expectedErr: errors.New("failed to create user BFP: invalid bodyfat percentage"),
		},
		{
			name: "error - GetSettingsByID fails",
			req: users.UpdateBodyMetricsReq{
				UserID:      testUserID,
				WeightValue: ptrFloat64(75),
			},
			setupMock: func(m *MockUserRepo) {
				m.On("GetSettingsByID", ctx, testUserID).Return(nil, errors.New("db error"))
			},
			expectedErr: errors.New("failed to get user settings: db error"),
		},
		{
			name: "error - UpdateBodyMetrics fails",
			req: users.UpdateBodyMetricsReq{
				UserID:      testUserID,
				WeightValue: ptrFloat64(70),
			},
			setupMock: func(m *MockUserRepo) {
				m.On("GetSettingsByID", ctx, testUserID).Return(
					&user.Settings{WeightUnit: user.Kg, HeightUnit: user.Cm},
					nil,
				)
				m.On("UpdateBodyMetrics", ctx, mock.Anything, testUserID).Return(errors.New("update failed"))
			},
			expectedErr: errors.New("failed to update body metrics update failed"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockUserRepo)
			tt.setupMock(mockRepo)
			svc := users.NewService(mockRepo)

			err := svc.UpdateBodyMetrics(ctx, tt.req)

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
