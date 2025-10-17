// Package postgres package postgres
package postgres

import (
	"context"
	"database/sql"

	"github.com/cheezecakee/logr"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/cheezecakee/fitrkr-athena/internal/core/domain/auth"
	"github.com/cheezecakee/fitrkr-athena/internal/core/ports"
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

const CreateRefreshToken = `INSERT INTO refresh_tokens (token, user_id, is_revoked, expires_at, created_at, updated_at) VALUES ($1,$2,$3,$4,$5,$6)`

func (r *AuthRepo) Add(ctx context.Context, refreshToken auth.RefreshToken) error {
	return WithTransaction(ctx, r.db, func(tx *sql.Tx) error {
		_, err := tx.ExecContext(ctx, CreateRefreshToken, refreshToken.Token, refreshToken.UserID, refreshToken.IsRevoked, refreshToken.ExpiresAt, refreshToken.CreatedAt, refreshToken.UpdatedAt)
		if err != nil {
			return err
		}

		logr.Get().Info("New refresh token created!")

		return nil
	})
}

const GetRefreshTokenByToken = `SELECT user_id, is_revoked, expires_at, revoked_at, created_at, updated_at FROM refresh_tokens WHERE token = $1`

func (r *AuthRepo) GetByToken(ctx context.Context, token string) (*auth.RefreshToken, error) {
	var row auth.RefreshToken

	err := r.db.QueryRowContext(ctx, GetRefreshTokenByToken, token).Scan(
		&row.UserID,
		&row.IsRevoked,
		&row.ExpiresAt,
		&row.RevokedAt,
		&row.CreatedAt,
		&row.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ports.ErrUserNotFound
		}
		return nil, err
	}

	row.Token = token
	return &row, nil
}

const GetRefreshTokenByID = `SELECT token, is_revoked, expires_at, revoked_at, created_at, updated_at FROM refresh_tokens WHERE user_id = $1`

func (r *AuthRepo) GetByID(ctx context.Context, userID string) ([]*auth.RefreshToken, error) {
	rows, err := r.db.QueryContext(ctx, GetRefreshTokenByID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tokens []*auth.RefreshToken
	for rows.Next() {
		var token auth.RefreshToken
		err := rows.Scan(
			&token.Token,
			&token.IsRevoked,
			&token.ExpiresAt,
			&token.RevokedAt,
			&token.CreatedAt,
			&token.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		token.UserID = uuid.MustParse(userID)
		tokens = append(tokens, &token)
	}
	return tokens, rows.Err()
}

const UpdateRefreshToken = `UPDATE refresh_tokens
	SET is_revoked = $2,
		expires_at = $3,
		revoked_at = $4,
		updated_at = $5
	WHERE token = $1
	`

func (r *AuthRepo) Update(ctx context.Context, refreshToken auth.RefreshToken) error {
	return WithTransaction(ctx, r.db, func(tx *sql.Tx) error {
		result, err := tx.ExecContext(ctx, UpdateRefreshToken, refreshToken.Token, refreshToken.IsRevoked, refreshToken.ExpiresAt, refreshToken.RevokedAt, refreshToken.UpdatedAt)
		if err != nil {
			return err
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			return err
		}

		if rowsAffected == 0 {
			return ports.ErrInvalidToken
		}

		logr.Get().Info("Refresh token updated!")
		return nil
	})
}

const DeleteRefreshToken = `DELETE from refresh_tokens WHERE token = $1`

func (r *AuthRepo) Delete(ctx context.Context, token string) error {
	return WithTransaction(ctx, r.db, func(tx *sql.Tx) error {
		result, err := tx.ExecContext(ctx, DeleteRefreshToken, token)
		if err != nil {
			return err
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			return err
		}

		if rowsAffected == 0 {
			return ports.ErrInvalidToken
		}

		logr.Get().Info("Refresh token deleted!")
		return nil
	})
}
