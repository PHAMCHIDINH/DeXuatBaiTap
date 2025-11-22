package exercises

import (
	"time"

	"chidinh/modules/predictions"
)

type CreateTemplateRequest struct {
	Name           string   `json:"name" binding:"required"`
	Intensity      string   `json:"intensity" binding:"required"`
	Description    string   `json:"description" binding:"required"`
	DurationMin    int32    `json:"duration_min" binding:"required"`
	FreqPerWeek    int32    `json:"freq_per_week" binding:"required"`
	TargetRiskLevel string  `json:"target_risk_level" binding:"required"` // low/medium/high/none
	Tags           []string `json:"tags"`
}

type TemplateResponse struct {
	ID             int64    `json:"id"`
	Name           string   `json:"name"`
	Intensity      string   `json:"intensity"`
	Description    string   `json:"description"`
	DurationMin    int32    `json:"duration_min"`
	FreqPerWeek    int32    `json:"freq_per_week"`
	TargetRiskLevel string  `json:"target_risk_level"`
	Tags           []string `json:"tags"`
}

type ListTemplatesResponse struct {
	Templates []TemplateResponse `json:"templates"`
}

type ListRecommendationsParams struct {
	Limit  int32 `form:"limit,default=10"`
	Offset int32 `form:"offset,default=0"`
}

type RecommendationResponse struct {
	ID           int64                        `json:"id"`
	PatientID    string                       `json:"patient_id"`
	PredictionID string                       `json:"prediction_id"`
	Plan         *predictions.RecommendationPlan `json:"plan,omitempty"`
	CreatedAt    time.Time                    `json:"created_at"`
}

type ListRecommendationsResponse struct {
	Recommendations []RecommendationResponse `json:"recommendations"`
}
