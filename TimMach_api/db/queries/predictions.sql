-- name: CreatePrediction :one
INSERT INTO predictions (
    patient_id,
    probability,
    risk_label,
    raw_features
)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetPredictionByID :one
SELECT * FROM predictions
WHERE id = $1
LIMIT 1;

-- name: GetLatestPredictionByPatient :one
SELECT * FROM predictions
WHERE patient_id = $1
ORDER BY created_at DESC
LIMIT 1;

-- name: ListPredictionsByPatient :many
SELECT * FROM predictions
WHERE patient_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: CountLatestRiskByUser :many
WITH pa AS (
    SELECT id
    FROM patients
    WHERE user_id = $1
),
latest AS (
    SELECT DISTINCT ON (pr.patient_id)
        pr.patient_id,
        pr.risk_label
    FROM predictions pr
    JOIN pa ON pa.id = pr.patient_id
    ORDER BY pr.patient_id, pr.created_at DESC
)
SELECT COALESCE(latest.risk_label, 'none') AS risk_label, COUNT(*) AS count
FROM pa p
LEFT JOIN latest ON latest.patient_id = p.id
GROUP BY COALESCE(latest.risk_label, 'none');
