package users_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/cheezecakee/fitrkr-athena/internal/core/ports"
	"github.com/cheezecakee/fitrkr-athena/internal/core/services/users"
)

func TestDelete(t *testing.T) {
	ctx := context.Background()
	testID := "test-user-id"

	tests := []struct {
		name          string
		req           users.DeleteAccountReq
		setupMock     func(*MockUserRepo)
		expectedErr   error
		shouldSucceed bool
	}{
		{
			name: "success - user deleted",
			req: users.DeleteAccountReq{
				ID: testID,
			},
			setupMock: func(m *MockUserRepo) {
				m.On("Delete", ctx, testID).Return(nil)
			},
			shouldSucceed: true,
		},
		{
			name: "error - user not found",
			req: users.DeleteAccountReq{
				ID: testID,
			},
			setupMock: func(m *MockUserRepo) {
				m.On("Delete", ctx, testID).Return(ports.ErrUserNotFound)
			},
			expectedErr: users.ErrUserNotFound,
		},
		{
			name: "error - repo delete fails",
			req: users.DeleteAccountReq{
				ID: testID,
			},
			setupMock: func(m *MockUserRepo) {
				m.On("Delete", ctx, testID).Return(errors.New("db error"))
			},
			expectedErr: errors.New("failed to delete user: db error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockUserRepo)
			tt.setupMock(mockRepo)
			svc := users.NewService(mockRepo)

			err := svc.Delete(ctx, tt.req)

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
