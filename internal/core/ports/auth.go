package ports

import (
	"context"
	"errors"

	"github.com/cheezecakee/fitrkr-athena/internal/core/domain/auth"
)

var ErrInvalidToken = errors.New("invalid token")

type AuthRepo interface {
	Add(ctx context.Context, refreshToken auth.RefreshToken) error
	GetByToken(ctx context.Context, token string) (auth.RefreshToken, error)
	GetByID(ctx context.Context, userID string) ([]auth.RefreshToken, error)
	Update(ctx context.Context, refreshToken auth.RefreshToken) error
	Delete(ctx context.Context, token string) error
}
