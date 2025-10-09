// Package postgres
package postgres

import (
	"context"
	"database/sql"
	"time"

	"github.com/cheezecakee/logr"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/cheezecakee/fitrkr-athena/internal/core/domain/user"
	"github.com/cheezecakee/fitrkr-athena/internal/ports"
)

type UserRepo struct {
	db      *sql.DB
	typeMap *pgtype.Map
}

func NewUserRepo(db *sql.DB) (*UserRepo, error) {
	return &UserRepo{
		db:      db,
		typeMap: pgtype.NewMap(),
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
	snapshot := user.UserSnapshot{
		ID:           r.ID,
		Username:     user.Username(r.Username),
		FullName:     r.FullName,
		Email:        user.Email(r.Email),
		Roles:        user.StringsToRoles(r.Roles),
		Stats:        user.Stats{},        // zero-value
		Subscription: user.Subscription{}, // zero-value
		Settings:     user.Settings{},     // zero-value
		CreatedAt:    r.CreatedAt,
		UpdatedAt:    r.UpdatedAt,
	}
	return snapshot.Reconstitute()
}

const CreateUser = `INSERT INTO users (id, username, full_name, email, roles, password_hash, created_at, updated_at) VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`

func (ur *UserRepo) Add(ctx context.Context, u user.User) error {
	return WithTransaction(ctx, ur.db, func(tx *sql.Tx) error {
		_, err := tx.ExecContext(ctx, CreateUser, u.ID, u.Username, u.FullName, u.Email, u.Roles(), u.Password, u.CreatedAt, u.UpdatedAt)
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

	var rolesArray pgtype.Array[string]

	err := ur.db.QueryRowContext(ctx, GetByUserID, id).Scan(
		&row.ID,
		&row.Username,
		&row.Email,
		&row.FullName,
		ur.typeMap.SQLScanner(&rolesArray),
		&row.CreatedAt,
		&row.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ports.ErrUserNotFound
		}
		return nil, err
	}

	row.Roles = rolesArray.Elements

	return row.ToDomain(), nil
}

const GetByUsername = `SELECT id, username, email, full_name, roles,created_at, updated_at from users WHERE username = $1`

func (ur *UserRepo) GetByUsername(ctx context.Context, username string) (*user.User, error) {
	var row UserRow

	var rolesArray pgtype.Array[string]

	err := ur.db.QueryRowContext(ctx, GetByUsername, username).Scan(
		&row.ID,
		&row.Username,
		&row.Email,
		&row.FullName,
		ur.typeMap.SQLScanner(&rolesArray),
		&row.CreatedAt,
		&row.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ports.ErrUserNotFound
		}
		return nil, err
	}

	row.Roles = rolesArray.Elements

	return row.ToDomain(), nil
}

const GetByUserEmail = `SELECT id, username, email, full_name, roles,created_at, updated_at from users WHERE email = $1`

func (ur *UserRepo) GetByEmail(ctx context.Context, email string) (*user.User, error) {
	var row UserRow
	var rolesArray pgtype.Array[string]

	err := ur.db.QueryRowContext(ctx, GetByUserEmail, email).Scan(
		&row.ID,
		&row.Username,
		&row.Email,
		&row.FullName,
		ur.typeMap.SQLScanner(&rolesArray),
		&row.CreatedAt,
		&row.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ports.ErrUserNotFound
		}
		return nil, err
	}

	row.Roles = rolesArray.Elements

	return row.ToDomain(), nil
}

// Password and roles will be updated separately at a later date

const UpdateUser = `UPDATE users 
	SET username = $2, 
		full_name = $3, 
		email = $4,
		updated_at = $5
	WHERE id = $1
`

func (ur *UserRepo) Update(ctx context.Context, u user.User) error {
	return WithTransaction(ctx, ur.db, func(tx *sql.Tx) error {
		result, err := tx.ExecContext(ctx, UpdateUser, u.ID, u.Username, u.FullName, u.Email, u.UpdatedAt)
		if err != nil {
			return err
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			return err
		}

		if rowsAffected == 0 {
			return ports.ErrUserNotFound
		}

		logr.Get().Info("User updated!")
		return nil
	})
}

const DeleteUser = `Delete from users WHERE id = $1`

func (ur *UserRepo) Delete(ctx context.Context, id string) error {
	return WithTransaction(ctx, ur.db, func(tx *sql.Tx) error {
		result, err := tx.ExecContext(ctx, DeleteUser, id)
		if err != nil {
			return err
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			return err
		}

		if rowsAffected == 0 {
			return ports.ErrUserNotFound
		}

		logr.Get().Info("User deleted!")
		return nil
	})
}
