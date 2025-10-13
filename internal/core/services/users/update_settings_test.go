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

func TestUpdateSettings(t *testing.T) {
	ctx := context.Background()
	testUserID := "test-user-id"

	tests := []struct {
		name          string
		req           users.UpdateSettingsReq
		setupMock     func(*MockUserRepo)
		expectedErr   error
		shouldSucceed bool
	}{
		{
			name: "success - update all settings",
			req: users.UpdateSettingsReq{
				UserID:          testUserID,
				WeightUnit:      stringPtr("kg"),
				HeightUnit:      stringPtr("cm"),
				Theme:           stringPtr("dark"),
				Visibility:      stringPtr("private"),
				EmailNotif:      boolPtr(true),
				PushNotif:       boolPtr(false),
				WorkoutReminder: boolPtr(true),
				StreakReminder:  boolPtr(true),
			},
			setupMock: func(m *MockUserRepo) {
				settings := &user.Settings{
					WeightUnit:      user.Kg,
					HeightUnit:      user.Cm,
					Theme:           user.Light,
					Visibility:      user.Public,
					EmailNotif:      false,
					PushNotif:       true,
					WorkoutReminder: false,
					StreakReminder:  false,
					CreatedAt:       time.Now(),
					UpdatedAt:       time.Now(),
				}
				m.On("GetSettingsByID", ctx, testUserID).Return(settings, nil)
				m.On("UpdateSettings", ctx, mock.MatchedBy(func(s user.Settings) bool {
					return s.WeightUnit == user.Kg && s.Theme == user.Dark && s.EmailNotif == true
				}), testUserID).Return(nil)
			},
			shouldSucceed: true,
		},
		{
			name: "success - update partial settings (only theme)",
			req: users.UpdateSettingsReq{
				UserID: testUserID,
				Theme:  stringPtr("light"),
			},
			setupMock: func(m *MockUserRepo) {
				settings := &user.Settings{
					WeightUnit:      user.Kg,
					HeightUnit:      user.Cm,
					Theme:           user.Dark,
					Visibility:      user.Private,
					EmailNotif:      true,
					PushNotif:       false,
					WorkoutReminder: true,
					StreakReminder:  true,
					CreatedAt:       time.Now(),
					UpdatedAt:       time.Now(),
				}
				m.On("GetSettingsByID", ctx, testUserID).Return(settings, nil)
				m.On("UpdateSettings", ctx, mock.MatchedBy(func(s user.Settings) bool {
					return s.Theme == user.Light
				}), testUserID).Return(nil)
			},
			shouldSucceed: true,
		},
		{
			name: "success - update notifications only",
			req: users.UpdateSettingsReq{
				UserID:          testUserID,
				EmailNotif:      boolPtr(false),
				PushNotif:       boolPtr(false),
				WorkoutReminder: boolPtr(false),
				StreakReminder:  boolPtr(false),
			},
			setupMock: func(m *MockUserRepo) {
				settings := &user.Settings{
					WeightUnit:      user.Kg,
					HeightUnit:      user.Cm,
					Theme:           user.Dark,
					Visibility:      user.Private,
					EmailNotif:      true,
					PushNotif:       true,
					WorkoutReminder: true,
					StreakReminder:  true,
					CreatedAt:       time.Now(),
					UpdatedAt:       time.Now(),
				}
				m.On("GetSettingsByID", ctx, testUserID).Return(settings, nil)
				m.On("UpdateSettings", ctx, mock.MatchedBy(func(s user.Settings) bool {
					return s.EmailNotif == false && s.PushNotif == false
				}), testUserID).Return(nil)
			},
			shouldSucceed: true,
		},
		{
			name: "error - invalid weight unit",
			req: users.UpdateSettingsReq{
				UserID:     testUserID,
				WeightUnit: stringPtr("invalid"),
			},
			setupMock: func(m *MockUserRepo) {
				settings := &user.Settings{
					WeightUnit: user.Kg,
					HeightUnit: user.Cm,
					CreatedAt:  time.Now(),
					UpdatedAt:  time.Now(),
				}
				m.On("GetSettingsByID", ctx, testUserID).Return(settings, nil)
			},
			expectedErr: errors.New("failed to create new weight unit: invalid weight unit"),
		},
		{
			name: "error - invalid height unit",
			req: users.UpdateSettingsReq{
				UserID:     testUserID,
				HeightUnit: stringPtr("invalid"),
			},
			setupMock: func(m *MockUserRepo) {
				settings := &user.Settings{
					WeightUnit: user.Kg,
					HeightUnit: user.Cm,
					CreatedAt:  time.Now(),
					UpdatedAt:  time.Now(),
				}
				m.On("GetSettingsByID", ctx, testUserID).Return(settings, nil)
			},
			expectedErr: errors.New("failed to create new height unit: invalid height unit"),
		},
		{
			name: "error - invalid theme",
			req: users.UpdateSettingsReq{
				UserID: testUserID,
				Theme:  stringPtr("invalid"),
			},
			setupMock: func(m *MockUserRepo) {
				settings := &user.Settings{
					Theme:     user.Dark,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}
				m.On("GetSettingsByID", ctx, testUserID).Return(settings, nil)
			},
			expectedErr: errors.New("failed to create new theme unit: invalid theme"),
		},
		{
			name: "error - invalid visibility",
			req: users.UpdateSettingsReq{
				UserID:     testUserID,
				Visibility: stringPtr("invalid"),
			},
			setupMock: func(m *MockUserRepo) {
				settings := &user.Settings{
					Visibility: user.Private,
					CreatedAt:  time.Now(),
					UpdatedAt:  time.Now(),
				}
				m.On("GetSettingsByID", ctx, testUserID).Return(settings, nil)
			},
			expectedErr: errors.New("failed to create new visibility: invalid visibility"),
		},
		{
			name: "error - get settings fails",
			req: users.UpdateSettingsReq{
				UserID: testUserID,
				Theme:  stringPtr("dark"),
			},
			setupMock: func(m *MockUserRepo) {
				m.On("GetSettingsByID", ctx, testUserID).Return(nil, errors.New("db error"))
			},
			expectedErr: errors.New("failed to get user settings: db error"),
		},
		{
			name: "error - update fails",
			req: users.UpdateSettingsReq{
				UserID: testUserID,
				Theme:  stringPtr("dark"),
			},
			setupMock: func(m *MockUserRepo) {
				settings := &user.Settings{
					Theme:     user.Light,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}
				m.On("GetSettingsByID", ctx, testUserID).Return(settings, nil)
				m.On("UpdateSettings", ctx, mock.Anything, testUserID).Return(errors.New("db error"))
			},
			expectedErr: errors.New("failed to update user settings: db error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockUserRepo)
			tt.setupMock(mockRepo)
			svc := users.NewService(mockRepo)

			err := svc.UpdateSettings(ctx, tt.req)

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
