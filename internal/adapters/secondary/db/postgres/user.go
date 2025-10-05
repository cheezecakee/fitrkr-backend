// Package postgres
package postgres

import (
	"context"
	"database/sql"
	"time"

	"github.com/cheezecakee/logr"
	"github.com/google/uuid"

	"github.com/cheezecakee/fitrkr-backend/internal/core/domain/user"
	"github.com/cheezecakee/fitrkr-backend/internal/ports"
)

type UserRepo struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) (*UserRepo, error) {
	return &UserRepo{
		db: db,
	}, nil
}

type UserRow struct {
	ID        uuid.UUID
	Username  string
	Email     string
	FullName  string
	Roles     []string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (r UserRow) ToDomain() *user.User {
	// roles := strings.Split(r.Roles, " ")
	return user.Reconstitute(r.ID, user.Username(r.Username), r.FullName, user.Email(r.Email), user.StringsToRoles(r.Roles), r.CreatedAt, r.UpdatedAt)
}

const CreateUser = `INSERT INTO users (id, username, full_name, email, roles, password_hash, created_at, updated_at) VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`

func (ur *UserRepo) Add(ctx context.Context, u user.User) error {
	return WithTransaction(ctx, ur.db, func(tx *sql.Tx) error {
		_, err := tx.ExecContext(ctx, CreateUser, u.ID, u.Username, u.FullName, u.Email, u.Roles(), u.Password, u.CreatedAt, u.UpdatedAt)
		logr.Get().Debugf("Roles: %v Type: %T", u.Roles(), u.Roles())
		if err != nil {
			return err
		}

		logr.Get().Info("New user created!")

		return nil
	})
}

const GetByUserID = `SELECT id, username, email, full_name, roles,created_at, updated_at from users WHERE id = $1`

func (ur *UserRepo) GetByID(ctx context.Context, id string) (*user.User, error) {
	var row UserRow

	err := ur.db.QueryRowContext(ctx, GetByUserID, id).Scan(
		&row.ID,
		&row.Username,
		&row.Email,
		&row.FullName,
		&row.Roles,
		&row.CreatedAt,
		&row.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ports.ErrUserNotFound
		}
		return nil, err
	}

	return row.ToDomain(), nil
}

const GetByUsername = `SELECT id, username, email, full_name, roles,created_at, updated_at from users WHERE username = $1`

func (ur *UserRepo) GetByUsername(ctx context.Context, username string) (*user.User, error) {
	var row UserRow

	err := ur.db.QueryRowContext(ctx, GetByUsername, username).Scan(
		&row.ID,
		&row.Username,
		&row.Email,
		&row.FullName,
		&row.Roles,
		&row.CreatedAt,
		&row.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ports.ErrUserNotFound
		}
		return nil, err
	}

	return row.ToDomain(), nil
}

const GetByUserEmail = `SELECT id, username, email, full_name, roles,created_at, updated_at from users WHERE email = $1`

func (ur *UserRepo) GetByEmail(ctx context.Context, email string) (*user.User, error) {
	var row UserRow

	err := ur.db.QueryRowContext(ctx, GetByUserEmail, email).Scan(
		&row.ID,
		&row.Username,
		&row.Email,
		&row.FullName,
		&row.Roles,
		&row.CreatedAt,
		&row.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ports.ErrUserNotFound
		}
		return nil, err
	}

	return row.ToDomain(), nil
}
