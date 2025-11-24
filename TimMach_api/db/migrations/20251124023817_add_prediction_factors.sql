-- +goose Up
ALTER TABLE predictions
ADD COLUMN IF NOT EXISTS factors JSONB NOT NULL DEFAULT '[]';

-- Backfill from raw_features if it already stores {factors:[...]}
UPDATE predictions
SET factors = COALESCE(raw_features->'factors', '[]'::jsonb)
WHERE factors = '[]'::jsonb;

-- +goose Down
ALTER TABLE predictions DROP COLUMN IF EXISTS factors;
