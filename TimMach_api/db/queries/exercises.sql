-- name: CreateExerciseTemplate :one
INSERT INTO exercise_templates (
    name, intensity, description, duration_min, freq_per_week, target_risk_level, tags
) VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: ListExerciseTemplates :many
SELECT * FROM exercise_templates
ORDER BY name ASC;

-- name: CreateExerciseRecommendation :one
INSERT INTO exercise_recommendations (
    patient_id, prediction_id, plan
) VALUES ($1, $2, $3)
RETURNING *;

-- name: ListExerciseRecommendationsByPatient :many
SELECT * FROM exercise_recommendations
WHERE patient_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: GetExerciseRecommendationByPrediction :one
SELECT * FROM exercise_recommendations
WHERE prediction_id = $1
LIMIT 1;
