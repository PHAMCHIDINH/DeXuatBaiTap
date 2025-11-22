-- +goose Up
CREATE SEQUENCE IF NOT EXISTS user_id_seq;

CREATE TABLE users (
    id TEXT PRIMARY KEY,
    email TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE patients (
    id BIGSERIAL PRIMARY KEY,
    user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    gender SMALLINT NOT NULL,
    dob DATE NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE predictions (
    id BIGSERIAL PRIMARY KEY,
    patient_id BIGINT NOT NULL REFERENCES patients(id) ON DELETE CASCADE,
    probability DOUBLE PRECISION NOT NULL,
    risk_label TEXT NOT NULL,
    raw_features JSONB NOT NULL DEFAULT '{}',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE exercise_templates (
    id BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    intensity TEXT NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    duration_min INT NOT NULL DEFAULT 0,
    freq_per_week INT NOT NULL DEFAULT 0,
    target_risk_level TEXT NOT NULL DEFAULT '',
    tags TEXT[] NOT NULL DEFAULT '{}'
);

CREATE TABLE exercise_recommendations (
    id BIGSERIAL PRIMARY KEY,
    patient_id BIGINT NOT NULL REFERENCES patients(id) ON DELETE CASCADE,
    prediction_id BIGINT NOT NULL REFERENCES predictions(id) ON DELETE CASCADE,
    plan JSONB NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_patients_user_id ON patients(user_id);
CREATE INDEX idx_predictions_patient_id ON predictions(patient_id);
CREATE INDEX idx_predictions_patient_created_at ON predictions(patient_id, created_at DESC);
CREATE INDEX idx_exercise_templates_intensity ON exercise_templates(intensity);
CREATE INDEX idx_exercise_templates_target_risk ON exercise_templates(target_risk_level);
CREATE INDEX idx_ex_rec_patient ON exercise_recommendations(patient_id);
CREATE INDEX idx_ex_rec_prediction ON exercise_recommendations(prediction_id);

-- +goose Down
DROP TABLE IF EXISTS exercise_recommendations;
DROP TABLE IF EXISTS exercise_templates;
DROP TABLE IF EXISTS predictions;
DROP TABLE IF EXISTS patients;
DROP TABLE IF EXISTS users;
DROP SEQUENCE IF EXISTS user_id_seq;
