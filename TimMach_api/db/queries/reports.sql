-- name: CreateReport :one
INSERT INTO reports (patient_id, filename, file_url, recipients)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetReportByID :one
SELECT * FROM reports
WHERE id = $1
LIMIT 1;

-- name: ListReportsByPatient :many
SELECT * FROM reports
WHERE patient_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: UpdateReportRecipients :one
UPDATE reports
SET recipients = $2
WHERE id = $1
RETURNING *;

-- name: DeleteReport :exec
DELETE FROM reports
WHERE id = $1;

-- name: CountReportsByPatient :one
SELECT COUNT(*) FROM reports
WHERE patient_id = $1;
