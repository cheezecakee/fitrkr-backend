package users_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/cheezecakee/fitrkr-athena/internal/core/ports"
	"github.com/cheezecakee/fitrkr-athena/internal/core/services/users"
)

func validGetUserByIDReq() users.GetUserByIDReq {
	return users.GetUserByIDReq{
		ID: uuid.New().String(),
	}
}

func TestGetUserByID(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name          string
		req           users.GetUserByIDReq
		setupMock     func(*MockUserRepo)
		expectedErr   error
		shouldSucceed bool
	}{
		{
			name: "success - user exists",
			req:  validGetUserByIDReq(),
			setupMock: func(m *MockUserRepo) {
				m.On("GetByID", ctx, mock.Anything).Return(&ports.User{
					ID:       uuid.New(),
					Username: "testuser",
					Email:    "test@example.com",
					FullName: "John Doe",
					Roles:    []string{"user"},
				}, nil)
			},
			shouldSucceed: true,
		},
		{
			name: "error - user not found",
			req:  validGetUserByIDReq(),
			setupMock: func(m *MockUserRepo) {
				m.On("GetByID", ctx, mock.Anything).Return(nil, errors.New("not found"))
			},
			expectedErr: users.ErrUserNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockUserRepo)
			tt.setupMock(mockRepo)
			svc := users.NewService(mockRepo)

			resp, err := svc.GetByID(ctx, tt.req)

			if tt.shouldSucceed {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				assert.Equal(t, "testuser", string(resp.Username))
			} else {
				assert.Error(t, err)
				assert.Nil(t, resp)
				assert.ErrorIs(t, err, tt.expectedErr)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func validGetUserByUsernameReq() users.GetUserByUsernameReq {
	return users.GetUserByUsernameReq{
		Username: "testuser123",
	}
}

func TestGetUserByUsername(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name          string
		req           users.GetUserByUsernameReq
		setupMock     func(*MockUserRepo)
		expectedErr   error
		shouldSucceed bool
	}{
		{
			name: "success - user exists",
			req:  validGetUserByUsernameReq(),
			setupMock: func(m *MockUserRepo) {
				m.On("GetByUsername", ctx, mock.Anything).Return(&ports.User{
					ID:       uuid.New(),
					Username: "testuser",
					Email:    "test@example.com",
					FullName: "John Doe",
					Roles:    []string{"user"},
				}, nil)
			},
			shouldSucceed: true,
		},
		{
			name: "error - user not found",
			req:  validGetUserByUsernameReq(),
			setupMock: func(m *MockUserRepo) {
				m.On("GetByUsername", ctx, mock.Anything).Return(nil, errors.New("not found"))
			},
			expectedErr: users.ErrUserNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockUserRepo)
			tt.setupMock(mockRepo)
			svc := users.NewService(mockRepo)

			resp, err := svc.GetByUsername(ctx, tt.req)

			if tt.shouldSucceed {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				assert.Equal(t, "testuser", string(resp.Username))
			} else {
				assert.Error(t, err)
				assert.Nil(t, resp)
				assert.ErrorIs(t, err, tt.expectedErr)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func validGetUserByEmailReq() users.GetUserByEmailReq {
	return users.GetUserByEmailReq{
		Email: "test@example.com",
	}
}

func TestGetUserByEmail(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name          string
		req           users.GetUserByEmailReq
		setupMock     func(*MockUserRepo)
		expectedErr   error
		shouldSucceed bool
	}{
		{
			name: "success - user exists",
			req:  validGetUserByEmailReq(),
			setupMock: func(m *MockUserRepo) {
				m.On("GetByEmail", ctx, mock.Anything).Return(&ports.User{
					ID:       uuid.New(),
					Username: "testuser",
					Email:    "test@example.com",
					FullName: "John Doe",
					Roles:    []string{"user"},
				}, nil)
			},
			shouldSucceed: true,
		},
		{
			name: "error - user not found",
			req:  validGetUserByEmailReq(),
			setupMock: func(m *MockUserRepo) {
				m.On("GetByEmail", ctx, mock.Anything).Return(nil, errors.New("not found"))
			},
			expectedErr: users.ErrUserNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockUserRepo)
			tt.setupMock(mockRepo)
			svc := users.NewService(mockRepo)

			resp, err := svc.GetByEmail(ctx, tt.req)

			if tt.shouldSucceed {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				assert.Equal(t, "testuser", string(resp.Username))
			} else {
				assert.Error(t, err)
				assert.Nil(t, resp)
				assert.ErrorIs(t, err, tt.expectedErr)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}
