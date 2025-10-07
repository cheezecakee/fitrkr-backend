-- +goose Up
CREATE TABLE users (
    id UUID PRIMARY KEY,
    username VARCHAR(255) UNIQUE NOT NULL,
    full_name VARCHAR(255) NOT NULL,
    password_hash TEXT NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    -- bio VARCHAR(160), 
    -- profile_picture TEXT,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    roles TEXT[] NOT NULL
);

CREATE TABLE user_stats (
    user_id UUID PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    
    -- Body metrics
    weight NUMERIC(5,2),
    height NUMERIC(5,2),
    body_fat_percent NUMERIC(4,1),
    
    -- Workout streaks
    rest_days INT NOT NULL,
    current_streak INT NOT NULL,
    longest_streak INT NOT NULL,
    last_workout_date DATE,
    
    -- Lifetime totals
    total_workouts INT NOT NULL,
    total_volume_lifted NUMERIC(10,2) NOT NULL, 
    total_time_minutes INT NOT NULL,
    
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

CREATE TABLE user_settings (
    user_id UUID PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,

    -- Display preferences 
    preferred_weight_unit VARCHAR(2) NOT NULL,
    preferred_height_unit VARCHAR(5) NOT NULL,
    theme VARCHAR(10) NOT NULL,

    -- Privacy 
    profile_visibility VARCHAR(10) NOT NULL,

    -- Notification preferences 
    email_notifications BOOLEAN NOT NULL,
    push_notifications BOOLEAN NOT NULL,
    workout_reminders BOOLEAN NOT NUll,
    streak_reminders BOOLEAN NOT NULL,

    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NUll
);

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION update_timestamp() RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;
-- +goose StatementEnd

CREATE TRIGGER update_users_timestamp
    BEFORE UPDATE ON users
    FOR EACH ROW EXECUTE FUNCTION update_timestamp();

CREATE TRIGGER update_user_stats_timestamp
    BEFORE UPDATE ON user_stats
    FOR EACH ROW EXECUTE FUNCTION update_timestamp();

CREATE TRIGGER update_user_settings_timestamp
    BEFORE UPDATE ON user_settings
    FOR EACH ROW EXECUTE FUNCTION update_timestamp();

-- +goose Down
DROP TABLE IF EXISTS user_settings;
DROP TABLE IF EXISTS user_stats;
DROP TABLE IF EXISTS subscriptions;
DROP TABLE IF EXISTS users;
