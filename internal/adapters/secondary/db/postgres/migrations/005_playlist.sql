
-- +goose Up
CREATE TABLE playlist (
    id SERIAL PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title VARCHAR(100) NOT NULL,
    description TEXT,
    is_active BOOLEAN DEFAULT FALSE,
    last_worked BOOLEAN DEFAULT FALSE,
    visibility VARCHAR(20) NOT NULL DEFAULT 'private', -- 'private', 'public', 'unlisted'
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    CONSTRAINT unique_user_playlist_title UNIQUE (user_id, title)
);

CREATE TABLE block (
    id SERIAL PRIMARY KEY,
    playlist_id INT NOT NULL REFERENCES playlist(id) ON DELETE CASCADE,
    exercise_id INT NOT NULL REFERENCES exercise(id) ON DELETE CASCADE,
    block_type VARCHAR(20) DEFAULT 'standard', -- 'standard', 'superset'
    block_order INT NOT NULL,
    rest_after_block_seconds INT DEFAULT 60 -- Rest after completing entire block
);

CREATE TABLE config (
    id SERIAL PRIMARY KEY,
    block_id INT NOT NULL REFERENCES block(id) ON DELETE CASCADE,

    -- Strength fields
    sets INT,
    reps_min INT,
    reps_max INT,
    weight NUMERIC(6,2),
    
    -- Cardio fields  
    duration_seconds INT,
    distance NUMERIC(6,2),
    target_pace NUMERIC(5,2),
    target_heart_rate INT,
    incline NUMERIC(4,1),
    
    rest_seconds INT DEFAULT 60,
    tempo INT[] CHECK (array_length(tempo, 1) = 4), -- [eccentric, pause, concentric, pause]
    
    -- Common fields
    notes TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Indexes
CREATE INDEX idx_block_playlist_id ON block(playlist_id);
CREATE INDEX idx_block_exercise_id ON block(exercise_id);
CREATE INDEX idx_config_block_id ON config(block_id);

-- Triggers
CREATE TRIGGER update_playlist_timestamp
    BEFORE UPDATE ON playlist
    FOR EACH ROW EXECUTE FUNCTION update_timestamp();

CREATE TRIGGER update_block_timestamp
    BEFORE UPDATE ON block
    FOR EACH ROW EXECUTE FUNCTION update_timestamp();

CREATE TRIGGER update_config_timestamp
    BEFORE UPDATE ON config
    FOR EACH ROW EXECUTE FUNCTION update_timestamp();
-- +goose Down
DROP TABLE IF EXISTS config;
DROP TABLE IF EXISTS block;
DROP TABLE IF EXISTS playlist;

