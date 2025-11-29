-- name: CreateUser :one
INSERT INTO users (id, email, keycloak_id)
VALUES ($1, $2, $3)
RETURNING id, email, created_at, keycloak_id;

-- name: CreateKeycloakUser :one
INSERT INTO users (id, email, keycloak_id)
VALUES ($1, $2, $3)
RETURNING id, email, created_at, keycloak_id;

-- name: AttachKeycloakID :one
UPDATE users
SET keycloak_id = $2
WHERE id = $1
RETURNING id, email, created_at, keycloak_id;

-- name: GetUserByID :one
SELECT id, email, created_at, keycloak_id FROM users
WHERE id = $1
LIMIT 1;

-- name: GetUserByEmail :one
SELECT id, email, created_at, keycloak_id FROM users
WHERE email = $1
LIMIT 1;

-- name: GetUserByKeycloakID :one
SELECT id, email, created_at, keycloak_id FROM users
WHERE keycloak_id = $1
LIMIT 1;

-- name: ListUsers :many
SELECT id, email, created_at, keycloak_id FROM users
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: NextUserSeq :one
SELECT nextval('user_id_seq')::bigint;
