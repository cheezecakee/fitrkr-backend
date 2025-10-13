package users_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/cheezecakee/fitrkr-athena/internal/core/domain/user"
	"github.com/cheezecakee/fitrkr-athena/internal/core/services/users"
)

func TestGetStats(t *testing.T) {
	ctx := context.Background()
	testID := "test-user-id"

	tests := []struct {
		name          string
		req           users.GetStatsReq
		setupMock     func(*MockUserRepo)
		expectedErr   error
		shouldSucceed bool
	}{
		{
			name: "success - stats retrieved with all data",
			req: users.GetStatsReq{
				ID: testID,
			},
			setupMock: func(m *MockUserRepo) {
				weight := user.WeightValue(75.5)
				height := user.HeightValue(180.0)
				bfp := user.BFP(15.5)
				stats := &user.Stats{
					Weight: &weight,
					Height: &height,
					BFP:    &bfp,
					Streak: user.Streak{
						RestDays:    2,
						Current:     5,
						Longest:     10,
						LastWorkout: time.Now(),
					},
					Totals: user.Totals{
						Workouts: 20,
						Lifted:   1000.0,
						Time:     300,
					},
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}
				settings := &user.Settings{
					WeightUnit:      user.Kg,
					HeightUnit:      user.Cm,
					Theme:           user.Dark,
					Visibility:      user.Private,
					EmailNotif:      true,
					PushNotif:       false,
					WorkoutReminder: true,
					StreakReminder:  true,
				}
				m.On("GetStatsByID", ctx, testID).Return(stats, nil)
				m.On("GetSettingsByID", ctx, testID).Return(settings, nil)
			},
			shouldSucceed: true,
		},
		{
			name: "success - stats retrieved with nil metrics",
			req: users.GetStatsReq{
				ID: testID,
			},
			setupMock: func(m *MockUserRepo) {
				stats := &user.Stats{
					Weight: nil,
					Height: nil,
					BFP:    nil,
					Streak: user.Streak{
						RestDays:    0,
						Current:     0,
						Longest:     0,
						LastWorkout: time.Time{},
					},
					Totals: user.Totals{
						Workouts: 0,
						Lifted:   0.0,
						Time:     0,
					},
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}
				settings := &user.Settings{
					WeightUnit:      user.Lb,
					HeightUnit:      user.Ft,
					Theme:           user.Light,
					Visibility:      user.Public,
					EmailNotif:      false,
					PushNotif:       false,
					WorkoutReminder: false,
					StreakReminder:  false,
				}
				m.On("GetStatsByID", ctx, testID).Return(stats, nil)
				m.On("GetSettingsByID", ctx, testID).Return(settings, nil)
			},
			shouldSucceed: true,
		},
		{
			name: "error - get stats fails",
			req: users.GetStatsReq{
				ID: testID,
			},
			setupMock: func(m *MockUserRepo) {
				m.On("GetStatsByID", ctx, testID).Return(nil, errors.New("stats db error"))
			},
			expectedErr: errors.New("failed to get stats: stats db error"),
		},
		{
			name: "error - get settings fails",
			req: users.GetStatsReq{
				ID: testID,
			},
			setupMock: func(m *MockUserRepo) {
				weight := user.WeightValue(75.5)
				height := user.HeightValue(180.0)
				stats := &user.Stats{
					Weight: &weight,
					Height: &height,
					Totals: user.Totals{
						Workouts: 20,
					},
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}
				m.On("GetStatsByID", ctx, testID).Return(stats, nil)
				m.On("GetSettingsByID", ctx, testID).Return(nil, errors.New("settings db error"))
			},
			expectedErr: errors.New("failed to get settings: settings db error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockUserRepo)
			tt.setupMock(mockRepo)
			svc := users.NewService(mockRepo)

			resp, err := svc.GetStats(ctx, tt.req)

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
