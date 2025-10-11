package users_test

import (
	"context"
	"errors"
	"testing"

	"github.com/cheezecakee/logr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/cheezecakee/fitrkr-athena/internal/core/domain/user"
	"github.com/cheezecakee/fitrkr-athena/internal/core/ports"
	"github.com/cheezecakee/fitrkr-athena/internal/core/services/users"
)

type MockUserRepo struct {
	mock.Mock
}

func (m *MockUserRepo) Add(ctx context.Context, u user.User) error {
	args := m.Called(ctx, u)
	return args.Error(0)
}

func (m *MockUserRepo) GetByUsername(ctx context.Context, username string) (*ports.User, error) {
	args := m.Called(ctx, username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ports.User), args.Error(1)
}

func (m *MockUserRepo) GetByEmail(ctx context.Context, email string) (*ports.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ports.User), args.Error(1)
}

func (m *MockUserRepo) GetByID(ctx context.Context, id string) (*ports.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ports.User), args.Error(1)
}

func (m *MockUserRepo) Update(ctx context.Context, u user.User) error {
	args := m.Called(ctx, u)
	return args.Error(0)
}

func (m *MockUserRepo) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockUserRepo) AddStats(ctx context.Context, stats user.Stats, userID string) error {
	args := m.Called(ctx, stats, userID)
	return args.Error(0)
}

func (m *MockUserRepo) GetStatsByID(ctx context.Context, userID string) (*ports.Stats, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ports.Stats), args.Error(1)
}

func (m *MockUserRepo) UpdateStats(ctx context.Context, stats user.Stats) error {
	args := m.Called(ctx, stats)
	return args.Error(0)
}

func (m *MockUserRepo) AddSubscription(ctx context.Context, sub user.Subscription, userID string) error {
	args := m.Called(ctx, sub, userID)
	return args.Error(0)
}

func (m *MockUserRepo) GetSubscriptionByID(ctx context.Context, userID string) (*ports.Subscription, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ports.Subscription), args.Error(1)
}

func (m *MockUserRepo) UpdateSubscription(ctx context.Context, sub user.Subscription) error {
	args := m.Called(ctx, sub)
	return args.Error(0)
}

func (m *MockUserRepo) AddSettings(ctx context.Context, settings user.Settings, userID string) error {
	args := m.Called(ctx, settings, userID)
	return args.Error(0)
}

func (m *MockUserRepo) GetSettingsByID(ctx context.Context, userID string) (*ports.Settings, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ports.Settings), args.Error(1)
}

func (m *MockUserRepo) UpdateSettings(ctx context.Context, settings user.Settings) error {
	args := m.Called(ctx, settings)
	return args.Error(0)
}

// Test helper to build valid test requests
func validCreateAccountReq() users.CreateAccountReq {
	return users.CreateAccountReq{
		Username:  "testuser123",
		Email:     "test@example.com",
		FirstName: "John",
		LastName:  "Doe",
		Roles:     []string{"user"},
		Password:  "SecurePass123!",
	}
}

func TestCreateAccount(t *testing.T) {
	logr.Init(&logr.PlainTextFormatter{}, logr.LevelInfo, nil)
	ctx := context.Background()

	tests := []struct {
		name          string
		req           users.CreateAccountReq
		setupMock     func(*MockUserRepo)
		expectedErr   error
		shouldSucceed bool
	}{
		{
			name: "success - creates user with all defaults",
			req:  validCreateAccountReq(),
			setupMock: func(m *MockUserRepo) {
				// User doesn't exist yet
				m.On("GetByUsername", ctx, "testuser123").Return(nil, ports.ErrUserNotFound)
				m.On("GetByEmail", ctx, "test@example.com").Return(nil, ports.ErrUserNotFound)
				// All repo calls succeed
				m.On("Add", ctx, mock.MatchedBy(func(u user.User) bool {
					return u.Username == "testuser123" && u.Email == "test@example.com"
				})).Return(nil)
				m.On("AddStats", ctx, mock.Anything, mock.Anything).Return(nil)
				m.On("AddSubscription", ctx, mock.Anything, mock.Anything).Return(nil)
				m.On("AddSettings", ctx, mock.Anything, mock.Anything).Return(nil)
			},
			shouldSucceed: true,
		},
		{
			name: "error - empty username",
			req: users.CreateAccountReq{
				Username:  "",
				Email:     "test@example.com",
				FirstName: "John",
				LastName:  "Doe",
				Roles:     []string{"user"},
				Password:  "SecurePass123!",
			},
			setupMock:   func(m *MockUserRepo) {}, // No repo calls expected
			expectedErr: user.ErrEmptyUsername,
		},
		{
			name: "error - username too short",
			req: users.CreateAccountReq{
				Username:  "ab",
				Email:     "test@example.com",
				FirstName: "John",
				LastName:  "Doe",
				Roles:     []string{"user"},
				Password:  "SecurePass123!",
			},
			setupMock:   func(m *MockUserRepo) {},
			expectedErr: user.ErrUsernameTooShort,
		},
		{
			name: "error - invalid email",
			req: users.CreateAccountReq{
				Username:  "testuser123",
				Email:     "notanemail",
				FirstName: "John",
				LastName:  "Doe",
				Roles:     []string{"user"},
				Password:  "SecurePass123!",
			},
			setupMock:   func(m *MockUserRepo) {},
			expectedErr: user.ErrInvalidEmail,
		},
		{
			name: "error - empty name",
			req: users.CreateAccountReq{
				Username:  "testuser123",
				Email:     "test@example.com",
				FirstName: "",
				LastName:  "Doe",
				Roles:     []string{"user"},
				Password:  "SecurePass123!",
			},
			setupMock:   func(m *MockUserRepo) {},
			expectedErr: user.ErrEmptyName,
		},
		{
			name: "error - name too short",
			req: users.CreateAccountReq{
				Username:  "testuser123",
				Email:     "test@example.com",
				FirstName: "J",
				LastName:  "Doe",
				Roles:     []string{"user"},
				Password:  "SecurePass123!",
			},
			setupMock:   func(m *MockUserRepo) {},
			expectedErr: user.ErrNameTooShort,
		},
		{
			name: "error - password too short",
			req: users.CreateAccountReq{
				Username:  "testuser123",
				Email:     "test@example.com",
				FirstName: "John",
				LastName:  "Doe",
				Roles:     []string{"user"},
				Password:  "Short1!",
			},
			setupMock:   func(m *MockUserRepo) {},
			expectedErr: user.ErrPasswordTooShort,
		},
		{
			name: "error - password missing uppercase",
			req: users.CreateAccountReq{
				Username:  "testuser123",
				Email:     "test@example.com",
				FirstName: "John",
				LastName:  "Doe",
				Roles:     []string{"user"},
				Password:  "securepass123!",
			},
			setupMock:   func(m *MockUserRepo) {},
			expectedErr: user.ErrPasswordNoUpper,
		},
		{
			name: "error - invalid role",
			req: users.CreateAccountReq{
				Username:  "testuser123",
				Email:     "test@example.com",
				FirstName: "John",
				LastName:  "Doe",
				Roles:     []string{"superadmin"},
				Password:  "SecurePass123!",
			},
			setupMock:   func(m *MockUserRepo) {},
			expectedErr: user.ErrInvalidRole,
		},
		{
			name: "error - duplicate username",
			req:  validCreateAccountReq(),
			setupMock: func(m *MockUserRepo) {
				// Username already exists
				m.On("GetByUsername", ctx, "testuser123").Return(&ports.User{}, nil)
			},
			expectedErr: users.ErrDuplicateUsername,
		},
		{
			name: "error - duplicate email",
			req:  validCreateAccountReq(),
			setupMock: func(m *MockUserRepo) {
				// Username is unique
				m.On("GetByUsername", ctx, "testuser123").Return(nil, ports.ErrUserNotFound)
				// But email exists
				m.On("GetByEmail", ctx, "test@example.com").Return(&ports.User{}, nil)
			},
			expectedErr: users.ErrDuplicateEmail,
		},
		{
			name: "error - repo.GetByUsername fails",
			req:  validCreateAccountReq(),
			setupMock: func(m *MockUserRepo) {
				// Repo connection error
				m.On("GetByUsername", ctx, "testuser123").Return(nil, errors.New("db connection failed"))
			},
			expectedErr: errors.New("db connection failed"),
		},
		{
			name: "error - repo.Add fails",
			req:  validCreateAccountReq(),
			setupMock: func(m *MockUserRepo) {
				m.On("GetByUsername", ctx, "testuser123").Return(nil, ports.ErrUserNotFound)
				m.On("GetByEmail", ctx, "test@example.com").Return(nil, ports.ErrUserNotFound)
				// User creation fails
				m.On("Add", ctx, mock.Anything).Return(errors.New("constraint violation"))
			},
			expectedErr: errors.New("constraint violation"),
		},
		{
			name: "error - repo.AddStats fails",
			req:  validCreateAccountReq(),
			setupMock: func(m *MockUserRepo) {
				m.On("GetByUsername", ctx, "testuser123").Return(nil, ports.ErrUserNotFound)
				m.On("GetByEmail", ctx, "test@example.com").Return(nil, ports.ErrUserNotFound)
				m.On("Add", ctx, mock.Anything).Return(nil)
				// Stats creation fails
				m.On("AddStats", ctx, mock.Anything, mock.Anything).Return(errors.New("stats insert failed"))
			},
			expectedErr: errors.New("stats insert failed"),
		},
		{
			name: "error - repo.AddSubscription fails",
			req:  validCreateAccountReq(),
			setupMock: func(m *MockUserRepo) {
				m.On("GetByUsername", ctx, "testuser123").Return(nil, ports.ErrUserNotFound)
				m.On("GetByEmail", ctx, "test@example.com").Return(nil, ports.ErrUserNotFound)
				m.On("Add", ctx, mock.Anything).Return(nil)
				m.On("AddStats", ctx, mock.Anything, mock.Anything).Return(nil)
				// Subscription creation fails
				m.On("AddSubscription", ctx, mock.Anything, mock.Anything).Return(errors.New("subscription insert failed"))
			},
			expectedErr: errors.New("subscription insert failed"),
		},
		{
			name: "error - repo.AddSettings fails",
			req:  validCreateAccountReq(),
			setupMock: func(m *MockUserRepo) {
				m.On("GetByUsername", ctx, "testuser123").Return(nil, ports.ErrUserNotFound)
				m.On("GetByEmail", ctx, "test@example.com").Return(nil, ports.ErrUserNotFound)
				m.On("Add", ctx, mock.Anything).Return(nil)
				m.On("AddStats", ctx, mock.Anything, mock.Anything).Return(nil)
				m.On("AddSubscription", ctx, mock.Anything, mock.Anything).Return(nil)
				// Settings creation fails
				m.On("AddSettings", ctx, mock.Anything, mock.Anything).Return(errors.New("settings insert failed"))
			},
			expectedErr: errors.New("settings insert failed"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock repo
			mockRepo := new(MockUserRepo)
			tt.setupMock(mockRepo)

			// Create service with mock
			svc := users.NewService(mockRepo)

			// Execute
			resp, err := svc.CreateAccount(ctx, tt.req)

			// Assert
			if tt.shouldSucceed {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				assert.NotEmpty(t, resp.UserID)
			} else {
				assert.Error(t, err)
				assert.Nil(t, resp)
				assert.ErrorContains(t, err, tt.expectedErr.Error())
			}

			// Verify all expectations were met
			mockRepo.AssertExpectations(t)
		})
	}
}
