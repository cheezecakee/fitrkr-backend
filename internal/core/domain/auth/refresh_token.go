// Package auth
package auth

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type RefreshToken struct {
	Token     string     `json:"token"`
	UserID    uuid.UUID  `json:"user_id"`
	IsRevoked bool       `json:"is_revoked"`
	ExpiresAt time.Time  `json:"expires_at"`
	RevokedAt *time.Time `json:"revoked_at"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

func NewRefreshToken(userID uuid.UUID) (RefreshToken, error) {
	token, err := MakeRefreshToken()
	if err != nil {
		return RefreshToken{}, fmt.Errorf("failed to generate refresh token: %w", err)
	}
	now := time.Now()
	return RefreshToken{
		Token:     token,
		UserID:    userID,
		IsRevoked: false,
		ExpiresAt: now.Add(7 * 24 * time.Hour),
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

func (rt *RefreshToken) IsExpired() bool {
	return time.Now().After(rt.ExpiresAt)
}

func (rt *RefreshToken) Revoke() {
	rt.IsRevoked = true
	now := time.Now()
	rt.RevokedAt = &now
	rt.UpdatedAt = now
}

func (rt *RefreshToken) Touch() {
	rt.UpdatedAt = time.Now()
}

// Helper function

func MakeRefreshToken() (string, error) {
	token := make([]byte, 32)
	_, err := rand.Read(token)
	if err != nil {
		// Add custom logger and err later
		return "", err
	}

	return hex.EncodeToString(token), nil
}
