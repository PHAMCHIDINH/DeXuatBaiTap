-- +goose Up
ALTER TABLE users ADD COLUMN IF NOT EXISTS keycloak_id TEXT UNIQUE;
ALTER TABLE users ALTER COLUMN password_hash DROP NOT NULL;
CREATE INDEX IF NOT EXISTS idx_users_keycloak_id ON users(keycloak_id);

-- +goose Down
DROP INDEX IF EXISTS idx_users_keycloak_id;
ALTER TABLE users ALTER COLUMN password_hash SET NOT NULL;
ALTER TABLE users DROP COLUMN IF EXISTS keycloak_id;
