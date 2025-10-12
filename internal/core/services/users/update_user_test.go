package users_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/cheezecakee/fitrkr-athena/internal/core/ports"
	"github.com/cheezecakee/fitrkr-athena/internal/core/services/users"
)

var testUserID = uuid.New()

func validUpdateUserReq() users.UpdateUserReq {
	return users.UpdateUserReq{
		ID:        testUserID.String(),
		Username:  "newuser",
		Email:     "new@example.com",
		FirstName: "Jane",
		LastName:  "Doe",
	}
}

// Helper to setup basic existing user
func basicExistingUser() *ports.User {
	return &ports.User{
		ID:        testUserID,
		Username:  "olduser",
		Email:     "old@example.com",
		FullName:  "Old Name",
		Roles:     []string{"user"},
		CreatedAt: time.Now().Add(-24 * time.Hour),
	}
}

// Helper to setup common username/email checks
func setupUniqueChecks(m *MockUserRepo, username, email string) {
	m.On("GetByUsername", mock.Anything, username).Return(nil, ports.ErrUserNotFound)
	m.On("GetByEmail", mock.Anything, email).Return(nil, ports.ErrUserNotFound)
}

func TestUpdateUser(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name          string
		req           users.UpdateUserReq
		setupMock     func(*MockUserRepo)
		expectedErr   error
		shouldSucceed bool
	}{
		{
			name: "success - update all fields",
			req:  validUpdateUserReq(),
			setupMock: func(m *MockUserRepo) {
				existing := basicExistingUser()
				m.On("GetByID", ctx, mock.Anything).Return(existing, nil)
				setupUniqueChecks(m, "newuser", "new@example.com")
				m.On("Update", ctx, mock.Anything).Return(nil)
			},
			shouldSucceed: true,
		},
		{
			name: "success - update partial fields (only email)",
			req: users.UpdateUserReq{
				ID:    uuid.New().String(),
				Email: "partial@example.com",
			},
			setupMock: func(m *MockUserRepo) {
				existing := basicExistingUser()
				m.On("GetByID", ctx, mock.Anything).Return(existing, nil)
				m.On("GetByEmail", mock.Anything, "partial@example.com").Return(nil, ports.ErrUserNotFound)
				m.On("Update", ctx, mock.Anything).Return(nil)
			},
			shouldSucceed: true,
		},
		{
			name: "error - user not found",
			req:  validUpdateUserReq(),
			setupMock: func(m *MockUserRepo) {
				m.On("GetByID", ctx, mock.Anything).Return(nil, errors.New("not found"))
			},
			expectedErr: errors.New("failed to get user: not found"),
		},
		{
			name: "error - duplicate username",
			req:  validUpdateUserReq(),
			setupMock: func(m *MockUserRepo) {
				existing := basicExistingUser()
				m.On("GetByID", ctx, mock.Anything).Return(existing, nil)
				m.On("GetByUsername", ctx, "newuser").Return(&ports.User{ID: uuid.New()}, nil)
			},
			expectedErr: users.ErrDuplicateUsername,
		},
		{
			name: "error - duplicate email",
			req:  validUpdateUserReq(),
			setupMock: func(m *MockUserRepo) {
				existing := basicExistingUser()
				m.On("GetByID", ctx, mock.Anything).Return(existing, nil)
				m.On("GetByUsername", ctx, "newuser").Return(nil, ports.ErrUserNotFound)
				m.On("GetByEmail", ctx, "new@example.com").Return(&ports.User{ID: uuid.New()}, nil)
			},
			expectedErr: users.ErrDuplicateEmail,
		},
		{
			name: "error - invalid username",
			req: users.UpdateUserReq{
				ID:       uuid.New().String(),
				Username: "ab", // too short
			},
			setupMock: func(m *MockUserRepo) {
				existing := basicExistingUser()
				m.On("GetByID", ctx, mock.Anything).Return(existing, nil)
			},
			expectedErr: errors.New("invalid username: username too short"),
		},
		{
			name: "error - invalid email",
			req: users.UpdateUserReq{
				ID:    uuid.New().String(),
				Email: "notanemail",
			},
			setupMock: func(m *MockUserRepo) {
				existing := basicExistingUser()
				m.On("GetByID", ctx, mock.Anything).Return(existing, nil)
			},
			expectedErr: errors.New("invalid email: invalid email"),
		},
		{
			name: "error - incomplete full name",
			req: users.UpdateUserReq{
				ID:        uuid.New().String(),
				FirstName: "OnlyFirst",
			},
			setupMock: func(m *MockUserRepo) {
				existing := basicExistingUser()
				m.On("GetByID", ctx, mock.Anything).Return(existing, nil)
			},
			expectedErr: errors.New("both first and last name must be provided together"),
		},
		{
			name: "error - repo update fails",
			req:  validUpdateUserReq(),
			setupMock: func(m *MockUserRepo) {
				existing := basicExistingUser()
				m.On("GetByID", ctx, mock.Anything).Return(existing, nil)
				setupUniqueChecks(m, "newuser", "new@example.com")
				m.On("Update", ctx, mock.Anything).Return(errors.New("db error"))
			},
			expectedErr: errors.New("error update user: db error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockUserRepo)
			tt.setupMock(mockRepo)
			svc := users.NewService(mockRepo)

			err := svc.Update(ctx, tt.req)

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
