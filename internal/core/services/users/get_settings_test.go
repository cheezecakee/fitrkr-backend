package users_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/cheezecakee/fitrkr-athena/internal/core/domain/user"
	"github.com/cheezecakee/fitrkr-athena/internal/core/ports"
	"github.com/cheezecakee/fitrkr-athena/internal/core/services/users"
)

func TestGetSettings(t *testing.T) {
	ctx := context.Background()
	testID := "test-user-id"

	tests := []struct {
		name          string
		req           users.GetSettingsReq
		setupMock     func(*MockUserRepo)
		expectedErr   error
		shouldSucceed bool
	}{
		{
			name: "success - settings retrieved",
			req: users.GetSettingsReq{
				ID: testID,
			},
			setupMock: func(m *MockUserRepo) {
				settings := &ports.Settings{
					WeightUnit:       user.Kg,
					HeightUnit:       user.Cm,
					Theme:            user.Dark,
					Visibility:       user.Private,
					EmailNotifs:      true,
					PushNotifs:       false,
					WorkoutReminders: true,
					StreakReminders:  true,
				}
				m.On("GetSettingsByID", ctx, testID).Return(settings, nil)
			},
			shouldSucceed: true,
		},
		{
			name: "success - settings with all disabled",
			req: users.GetSettingsReq{
				ID: testID,
			},
			setupMock: func(m *MockUserRepo) {
				settings := &ports.Settings{
					WeightUnit:       user.Lb,
					HeightUnit:       user.Ft,
					Theme:            user.Light,
					Visibility:       user.Public,
					EmailNotifs:      false,
					PushNotifs:       false,
					WorkoutReminders: false,
					StreakReminders:  false,
				}
				m.On("GetSettingsByID", ctx, testID).Return(settings, nil)
			},
			shouldSucceed: true,
		},
		{
			name: "error - repo fails",
			req: users.GetSettingsReq{
				ID: testID,
			},
			setupMock: func(m *MockUserRepo) {
				m.On("GetSettingsByID", ctx, testID).Return(nil, errors.New("db error"))
			},
			expectedErr: errors.New("failed to get settings: db error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockUserRepo)
			tt.setupMock(mockRepo)
			svc := users.NewService(mockRepo)

			resp, err := svc.GetSettings(ctx, tt.req)

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

// Helper function for pointer to float64
func ptrFloat64(f float64) *float64 {
	return &f
}
