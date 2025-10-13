// Package postgres
package postgres

import (
	"context"
	"database/sql"

	"github.com/cheezecakee/logr"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/cheezecakee/fitrkr-athena/internal/core/domain/user"
	"github.com/cheezecakee/fitrkr-athena/internal/core/ports"
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

const CreateUser = `INSERT INTO users (id, username, full_name, email, roles, password_hash, created_at, updated_at) VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`

func (ur *UserRepo) Add(ctx context.Context, u user.User) error {
	return WithTransaction(ctx, ur.db, func(tx *sql.Tx) error {
		_, err := tx.ExecContext(ctx, CreateUser, u.ID, u.Username, u.FullName, u.Email, u.Roles, u.Password, u.CreatedAt, u.UpdatedAt)
		if err != nil {
			return err
		}

		logr.Get().Info("New user created!")

		return nil
	})
}

const GetByUserID = `SELECT id, username, email, full_name, roles,created_at, updated_at from users WHERE id = $1`

func (ur *UserRepo) GetByID(ctx context.Context, id string) (*ports.User, error) {
	var row ports.User

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

	return &row, nil
}

const GetByUsername = `SELECT id, username, email, full_name, roles,created_at, updated_at from users WHERE username = $1`

func (ur *UserRepo) GetByUsername(ctx context.Context, username string) (*ports.User, error) {
	var row ports.User

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

	return &row, nil
}

const GetByUserEmail = `SELECT id, username, email, full_name, roles,created_at, updated_at from users WHERE email = $1`

func (ur *UserRepo) GetByEmail(ctx context.Context, email string) (*ports.User, error) {
	var row ports.User
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

	return &row, nil
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

const CreateUserSettings = `INSERT INTO user_settings (user_id, preferred_weight_unit, preferred_height_unit, theme, profile_visibility, email_notifications, push_notifications, workout_reminders, streak_reminders, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`

func (ur *UserRepo) AddSettings(ctx context.Context, us user.Settings, id string) error {
	return WithTransaction(ctx, ur.db, func(tx *sql.Tx) error {
		_, err := tx.ExecContext(ctx, CreateUserSettings, id, us.WeightUnit, us.HeightUnit, us.Theme, us.Visibility, us.EmailNotif, us.PushNotif, us.WorkoutReminder, us.StreakReminder, us.CreatedAt, us.UpdatedAt)
		if err != nil {
			return err
		}

		logr.Get().Info("New user settings created!")

		return nil
	})
}

const CreateUserStats = `INSERT INTO user_stats (user_id, rest_days, current_streak, longest_streak,  total_workouts, total_lifted, total_time_minutes, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`

func (ur *UserRepo) AddStats(ctx context.Context, us user.Stats, id string) error {
	return WithTransaction(ctx, ur.db, func(tx *sql.Tx) error {
		_, err := tx.ExecContext(ctx, CreateUserStats, id, us.Streak.RestDays, us.Streak.Current, us.Streak.Longest, us.Totals.Workouts, us.Totals.Lifted, us.Totals.Time, us.CreatedAt, us.UpdatedAt)
		if err != nil {
			return err
		}

		logr.Get().Info("New user stats created!")

		return nil
	})
}

const CreateUserSubscription = `INSERT INTO user_subscription (user_id, plan, billing_period, started_at, auto_renew, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7)`

func (ur *UserRepo) AddSubscription(ctx context.Context, us user.Subscription, id string) error {
	return WithTransaction(ctx, ur.db, func(tx *sql.Tx) error {
		_, err := tx.ExecContext(ctx, CreateUserSubscription, id, us.Plan, us.BillingPeriod, us.StartedAt, us.AutoRenew, us.CreatedAt, us.UpdatedAt)
		if err != nil {
			return err
		}

		logr.Get().Info("New user subscription created!")

		return nil
	})
}

const GetUserStats = `SELECT weight, height, body_fat_percent, rest_days, current_streak, longest_streak, last_workout_date, total_workouts, total_lifted, total_time_minutes, created_at, updated_at FROM user_stats WHERE user_id = $1`

func (ur *UserRepo) GetStatsByID(ctx context.Context, userID string) (*user.Stats, error) {
	var row user.Stats

	err := ur.db.QueryRowContext(ctx, GetUserStats, userID).Scan(
		&row.Weight,
		&row.Height,
		&row.BFP,
		&row.Streak.RestDays,
		&row.Streak.Current,
		&row.Streak.Longest,
		&row.Streak.LastWorkout,
		&row.Totals.Workouts,
		&row.Totals.Lifted,
		&row.Totals.Time,
		&row.CreatedAt,
		&row.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ports.ErrUserNotFound
		}
		return nil, err
	}

	return &row, nil
}

const GetUserSettings = `SELECT preferred_weight_unit, preferred_height_unit, theme, profile_visibility, email_notifications, push_notifications, workout_reminders, streak_reminders , created_at, updated_at FROM user_settings WHERE user_id = $1`

func (ur *UserRepo) GetSettingsByID(ctx context.Context, userID string) (*user.Settings, error) {
	var row user.Settings

	err := ur.db.QueryRowContext(ctx, GetUserSettings, userID).Scan(
		&row.WeightUnit,
		&row.HeightUnit,
		&row.Theme,
		&row.Visibility,
		&row.EmailNotif,
		&row.PushNotif,
		&row.WorkoutReminder,
		&row.StreakReminder,
		&row.CreatedAt,
		&row.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ports.ErrUserNotFound
		}
		return nil, err
	}

	return &row, nil
}

const GetUserSubscription = `SELECT plan, billing_period, started_at, expires_at, auto_renew, cancelled_at, last_payment_at, last_payment_amount, last_payment_currency, trial_ends_at, created_at, updated_at FROM user_subscription WHERE user_id = $1`

func (ur *UserRepo) GetSubscriptionByID(ctx context.Context, userID string) (*user.Subscription, error) {
	var row user.Subscription

	err := ur.db.QueryRowContext(ctx, GetUserSubscription, userID).Scan(
		&row.Plan,
		&row.BillingPeriod,
		&row.StartedAt,
		&row.ExpiresAt,
		&row.AutoRenew,
		&row.CancelledAt,
		&row.LastPaymentAt,
		&row.LastPaymentAmount,
		&row.LastPaymentCurrency,
		&row.TrialEndsAt,
		&row.CreatedAt,
		&row.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ports.ErrUserNotFound
		}
		return nil, err
	}

	return &row, nil
}

const UpdateUserSubscription = `UPDATE user_subscription
	SET plan = $2, 
		billing_period = $3, 
		started_at = $4,
		expires_at = $5,
		auto_renew = $6,
		cancelled_at = $7,
		last_payment_at = $8,
		last_payment_amount = $9,
		last_paymet_currency= $10,
		trial_ends_at = $11,
		updated_at = $12
	WHERE user_id = $1
`

func (ur *UserRepo) UpdateSubscription(ctx context.Context, sub user.Subscription, userID string) error {
	return WithTransaction(ctx, ur.db, func(tx *sql.Tx) error {
		result, err := tx.ExecContext(ctx, UpdateUserSubscription, userID, sub.Plan, sub.BillingPeriod, sub.StartedAt, sub.ExpiresAt, sub.AutoRenew, sub.CancelledAt, sub.LastPaymentAt, sub.LastPaymentAmount, sub.LastPaymentCurrency, sub.TrialEndsAt, sub.UpdatedAt)
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

		logr.Get().Info("User subscription updated!")
		return nil
	})
}

const UpdateUserSettings = `UPDATE users_settings
	SET preferred_weight_unit = $2,
    	preferred_height_unit = $3,
    	theme = $4,
    	profile_visibility= $5,
    	email_notifications= $6,
    	push_notifications= $7,
    	workout_reminders= $8,
    	streak_reminders= $9,
		updated_at = $10
	WHERE user_id = $1
`

func (ur *UserRepo) UpdateSettings(ctx context.Context, settings user.Settings, userID string) error {
	return WithTransaction(ctx, ur.db, func(tx *sql.Tx) error {
		result, err := tx.ExecContext(ctx, UpdateUserSettings, userID, settings.WeightUnit, settings.HeightUnit, settings.Theme, settings.Visibility, settings.EmailNotif, settings.PushNotif, settings.WorkoutReminder, settings.StreakReminder, settings.UpdatedAt)
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

		logr.Get().Info("User settings updated!")
		return nil
	})
}
