-- +goose Up
ALTER TABLE users DROP COLUMN password_hash;

-- +goose Down
ALTER TABLE users ADD COLUMN password_hash TEXT;
