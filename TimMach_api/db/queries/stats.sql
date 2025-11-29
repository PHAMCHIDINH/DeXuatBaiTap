-- name: GetTotalPatients :one
SELECT COUNT(*) FROM patients
WHERE user_id = $1;

-- name: GetTotalPredictions :one
SELECT COUNT(*) 
FROM predictions p
JOIN patients pa ON p.patient_id = pa.id
WHERE pa.user_id = $1;

-- name: GetRiskDistribution :many
WITH latest_predictions AS (
    SELECT DISTINCT ON (p.patient_id) p.risk_label
    FROM predictions p
    JOIN patients pa ON p.patient_id = pa.id
    WHERE pa.user_id = $1
    ORDER BY p.patient_id, p.created_at DESC
)
SELECT 
    risk_label,
    COUNT(*) AS count
FROM latest_predictions
GROUP BY risk_label;
