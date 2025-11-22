-- name: CreatePatient :one
INSERT INTO patients (user_id, name, gender, dob)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetPatientByID :one
SELECT * FROM patients
WHERE id = $1
LIMIT 1;

-- name: ListPatientsByUser :many
SELECT * FROM patients
WHERE user_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: ListPatientsWithLatestPrediction :many
WITH pa AS (
    SELECT *
    FROM patients
    WHERE user_id = $1
),
latest AS (
    SELECT DISTINCT ON (pr.patient_id)
        pr.patient_id,
        pr.probability,
        pr.risk_label,
        pr.created_at
    FROM predictions pr
    JOIN pa ON pa.id = pr.patient_id
    ORDER BY pr.patient_id, pr.created_at DESC
)
SELECT
    p.id,
    p.user_id,
    p.name,
    p.gender,
    p.dob,
    p.created_at,
    l.probability AS latest_probability,
    l.risk_label AS latest_risk_label,
    l.created_at AS latest_prediction_at
FROM pa p
LEFT JOIN latest l ON l.patient_id = p.id
WHERE ($4 = '' OR COALESCE(l.risk_label, 'none') = $4)
ORDER BY p.created_at DESC
LIMIT $2 OFFSET $3;

-- name: UpdatePatient :one
UPDATE patients
SET
    name = COALESCE($2, name),
    gender = COALESCE($3, gender),
    dob = COALESCE($4, dob)
WHERE id = $1
RETURNING *;

-- name: DeletePatient :exec
DELETE FROM patients
WHERE id = $1;

-- name: CountPatientsByUser :one
SELECT COUNT(*) FROM patients WHERE user_id = $1;
