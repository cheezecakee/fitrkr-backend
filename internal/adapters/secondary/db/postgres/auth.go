// Package postgres package postgres
package postgres

import (
	"context"
	"database/sql"

	"github.com/jackc/pgx/v5/pgtype"

	"github.com/cheezecakee/fitrkr-athena/internal/core/domain/auth"
)

type AuthRepo struct {
	db      *sql.DB
	typeMap *pgtype.Map
}

func NewAuthRepo(db *sql.DB) (*AuthRepo, error) {
	return &AuthRepo{
		db:      db,
		typeMap: pgtype.NewMap(),
	}, nil
}

func (ar *AuthRepo) Add(ctx context.Context, refreshToken auth.RefreshToken) error {
	return nil
}

func (ar *AuthRepo) GetByToken(ctx context.Context, token string) (auth.RefreshToken, error) {
	return auth.RefreshToken{}, nil
}

func (ar *AuthRepo) GetByID(ctx context.Context, userID string) ([]auth.RefreshToken, error) {
	return nil, nil
}

func (ar *AuthRepo) Update(ctx context.Context, refreshToken auth.RefreshToken) error {
	return nil
}

func (ar *AuthRepo) Delete(ctx context.Context, token string) error {
	return nil
}
