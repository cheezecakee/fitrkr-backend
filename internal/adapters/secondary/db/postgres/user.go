// Package postgres
package postgres

import (
	"context"
	"database/sql"

	"github.com/cheezecakee/logr"

	"github.com/cheezecakee/fitrkr-backend/internal/core/domain/user"
)

type UserRepo struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) (*UserRepo, error) {
	return &UserRepo{
		db: db,
	}, nil
}

const CreateUser = `INSERT INTO users (id, username, full_name, email, roles, password_hash, created_at, updated_at) VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`

func (ur *UserRepo) Add(ctx context.Context, u user.User) error {
	tx, err := ur.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, CreateUser, u.ID, u.Username, u.FullName, u.Email, u.Roles(), u.Password, u.CreatedAt, u.UpdatedAt)
	if err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	logr.Get().Info("New user created!")

	return nil
}
