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

func TestGetStats(t *testing.T) {
	ctx := context.Background()
	testID := "test-user-id"
	weight := 75.5
	height := 180.0

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
				stats := &ports.Stats{
					WeightValue:   &weight,
					HeightValue:   &height,
					BFP:           ptrFloat64(15.5),
					RestDays:      2,
					Current:       5,
					Longest:       10,
					TotalWorkouts: 20,
					TotalLifted:   1000.0,
					TotalTime:     300,
				}
				settings := &ports.Settings{
					WeightUnit: user.Kg,
					HeightUnit: user.Cm,
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
				stats := &ports.Stats{
					WeightValue:   nil,
					HeightValue:   nil,
					BFP:           nil,
					RestDays:      0,
					Current:       0,
					Longest:       0,
					TotalWorkouts: 0,
					TotalLifted:   0.0,
					TotalTime:     0,
				}
				settings := &ports.Settings{
					WeightUnit: user.Lb,
					HeightUnit: user.Ft,
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
				stats := &ports.Stats{
					WeightValue:   &weight,
					HeightValue:   &height,
					TotalWorkouts: 20,
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
