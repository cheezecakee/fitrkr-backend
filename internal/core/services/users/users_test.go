package users_test

import (
	"context"
	"os"
	"testing"

	"github.com/cheezecakee/logr"
	"github.com/stretchr/testify/mock"

	"github.com/cheezecakee/fitrkr-athena/internal/core/domain/user"
	"github.com/cheezecakee/fitrkr-athena/internal/core/ports"
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

func TestMain(m *testing.M) {
	logr.Init(&logr.PlainTextFormatter{}, logr.LevelInfo, nil)

	exitCode := m.Run()

	os.Exit(exitCode)
}
