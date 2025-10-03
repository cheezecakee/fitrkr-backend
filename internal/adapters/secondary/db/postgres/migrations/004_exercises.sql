-- +goose Up

-- Exercise Category
CREATE TABLE exercise_category (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) UNIQUE NOT NULL
);

-- Equipment
CREATE TABLE equipment (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) UNIQUE NOT NULL
);

-- Training Type (independent tags)
CREATE TABLE training_type (
    id SERIAL PRIMARY KEY,
    name VARCHAR(20) UNIQUE NOT NULL
);

-- Exercise
CREATE TABLE exercise (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) UNIQUE NOT NULL,
    description TEXT NOT NULL,
    category_id INT NOT NULL REFERENCES exercise_category(id),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Trigger to auto-update updated_at
CREATE TRIGGER update_exercise_timestamp
    BEFORE UPDATE ON exercise
    FOR EACH ROW EXECUTE FUNCTION update_timestamp();

-- Muscle Group
CREATE TABLE muscle_group (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) UNIQUE NOT NULL
);

-- Join Table: Exercise ↔ Muscle Groups
CREATE TABLE exercise_muscles (
    exercise_id INT NOT NULL REFERENCES exercise(id) ON DELETE CASCADE,
    muscle_group_id INT NOT NULL REFERENCES muscle_group(id) ON DELETE CASCADE,
    PRIMARY KEY (exercise_id, muscle_group_id)
);

-- Join Table: Exercise ↔ Training Types
CREATE TABLE exercise_training_types (
    exercise_id INT NOT NULL REFERENCES exercise(id) ON DELETE CASCADE,
    training_type_id INT NOT NULL REFERENCES training_type(id) ON DELETE CASCADE,
    PRIMARY KEY (exercise_id, training_type_id)
);

-- Join Table: Exercise ↔ Equipment
CREATE TABLE exercise_equipment (
    exercise_id INT NOT NULL REFERENCES exercise(id) ON DELETE CASCADE,
    equipment_id INT NOT NULL REFERENCES equipment(id) ON DELETE CASCADE,
    PRIMARY KEY (exercise_id, equipment_id)
);

-- Indexes for fast lookups
CREATE INDEX idx_exercise_muscles_exercise_id ON exercise_muscles(exercise_id);
CREATE INDEX idx_exercise_muscles_muscle_group_id ON exercise_muscles(muscle_group_id);
CREATE INDEX idx_exercise_training_types_exercise_id ON exercise_training_types(exercise_id);
CREATE INDEX idx_exercise_training_types_training_type_id ON exercise_training_types(training_type_id);
CREATE INDEX idx_exercise_equipment_exercise_id ON exercise_equipment(exercise_id);
CREATE INDEX idx_exercise_equipment_equipment_id ON exercise_equipment(equipment_id);

-- +goose Down
DROP TABLE IF EXISTS exercise_equipment;
DROP TABLE IF EXISTS exercise_training_types;
DROP TABLE IF EXISTS exercise_muscles;
DROP TABLE IF EXISTS exercise;
DROP TABLE IF EXISTS training_type;
DROP TABLE IF EXISTS muscle_group;
DROP TABLE IF EXISTS equipment;
DROP TABLE IF EXISTS exercise_category;

-- TODO: Consider adding equipment_muscles table
-- To associate equipment with muscle groups
-- Useful for recommendation, filtering, or workout generation
