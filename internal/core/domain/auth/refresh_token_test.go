package auth_test

import (
	"encoding/hex"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/cheezecakee/fitrkr-athena/internal/core/domain/auth"
)

func TestNewRefreshToken(t *testing.T) {
	userID := uuid.New()
	beforeCreation := time.Now()

	rf, err := auth.NewRefreshToken(userID)

	afterCreation := time.Now()

	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	tests := []struct {
		name    string
		check   func() bool
		message string
	}{
		{
			"userID matches input",
			func() bool { return rf.UserID == userID },
			"expected userID to match",
		},
		{
			"token is not empty",
			func() bool { return rf.Token != "" },
			"expected token to be generated",
		},
		{
			"token is valid hex",
			func() bool {
				_, err := hex.DecodeString(rf.Token)
				return err == nil
			},
			"expected valid hex string",
		},
		{
			"token is 64 chars (32 bytes as hex)",
			func() bool { return len(rf.Token) == 64 },
			"expected token length 64",
		},
		{
			"isRevoked is false",
			func() bool { return rf.IsRevoked == false },
			"expected isRevoked to be false",
		},
		{
			"revokedAt is nil",
			func() bool { return rf.RevokedAt == nil },
			"expected revokedAt to be nil",
		},
		{
			"expiresAt is 7 days in future",
			func() bool {
				minExpiry := beforeCreation.Add(7 * 24 * time.Hour)
				maxExpiry := afterCreation.Add(7 * 24 * time.Hour)
				return rf.ExpiresAt.After(minExpiry) && rf.ExpiresAt.Before(maxExpiry.Add(1*time.Second))
			},
			"expected expiresAt to be ~7 days from now",
		},
		{
			"createdAt is recent",
			func() bool {
				return rf.CreatedAt.After(beforeCreation.Add(-1*time.Second)) &&
					rf.CreatedAt.Before(afterCreation.Add(1*time.Second))
			},
			"expected createdAt to be recent",
		},
		{
			"updatedAt matches createdAt",
			func() bool { return rf.UpdatedAt.Equal(rf.CreatedAt) },
			"expected updatedAt to equal createdAt",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !tt.check() {
				t.Error(tt.message)
			}
		})
	}
}

func TestNewRefreshToken_UniqueTokens(t *testing.T) {
	userID := uuid.New()

	rf1, _ := auth.NewRefreshToken(userID)
	rf2, _ := auth.NewRefreshToken(userID)

	if rf1.Token == rf2.Token {
		t.Error("expected different tokens for each call")
	}
}

func TestMakeRefreshToken(t *testing.T) {
	token, err := auth.MakeRefreshToken()
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if len(token) != 64 {
		t.Errorf("expected token length 64, got %d", len(token))
	}

	if _, err := hex.DecodeString(token); err != nil {
		t.Errorf("expected valid hex string, got %v", err)
	}
}

func TestIsExpired(t *testing.T) {
	userID := uuid.New()
	rf, _ := auth.NewRefreshToken(userID)

	t.Run("not expired immediately after creation", func(t *testing.T) {
		if rf.IsExpired() {
			t.Error("expected token to not be expired right after creation")
		}
	})

	t.Run("not expired before expiration time", func(t *testing.T) {
		// This will always pass unless your clock is weird
		// but it's good for documentation
		if rf.IsExpired() {
			t.Error("expected token to not be expired")
		}
	})
}

func TestRevoke(t *testing.T) {
	userID := uuid.New()
	rf, _ := auth.NewRefreshToken(userID)
	beforeRevoke := time.Now()

	rf.Revoke()

	afterRevoke := time.Now()

	if !rf.IsRevoked {
		t.Error("expected isRevoked to be true")
	}

	if rf.RevokedAt == nil {
		t.Fatal("expected revokedAt to be set")
	}

	if rf.RevokedAt.Before(beforeRevoke) || rf.RevokedAt.After(afterRevoke.Add(1*time.Second)) {
		t.Errorf("expected revokedAt to be between %v and %v, got %v", beforeRevoke, afterRevoke, *rf.RevokedAt)
	}

	if rf.UpdatedAt.Before(beforeRevoke) || rf.UpdatedAt.After(afterRevoke.Add(1*time.Second)) {
		t.Errorf("expected updatedAt to be updated during revoke")
	}
}

func TestTouch(t *testing.T) {
	userID := uuid.New()
	rf, _ := auth.NewRefreshToken(userID)

	// Small delay to ensure time advances
	time.Sleep(10 * time.Millisecond)

	beforeTouch := time.Now()
	rf.Touch()
	afterTouch := time.Now()

	if rf.UpdatedAt.Before(beforeTouch) || rf.UpdatedAt.After(afterTouch.Add(1*time.Second)) {
		t.Errorf("expected updatedAt to be updated during touch")
	}
}
