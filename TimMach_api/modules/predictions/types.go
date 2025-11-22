package predictions

import (
	"encoding/json"
	"time"
)

// MLRequest định nghĩa payload gửi sang service FastAPI.
type MLRequest struct {
	AgeYears    float64 `json:"age_years"`
	Gender      int     `json:"gender"`
	Height      float64 `json:"height"`
	Weight      float64 `json:"weight"`
	APHi        int     `json:"ap_hi"`
	APLo        int     `json:"ap_lo"`
	Cholesterol int     `json:"cholesterol"`
	Gluc        int     `json:"gluc"`
	Smoke       int     `json:"smoke"`
	Alco        int     `json:"alco"`
	Active      int     `json:"active"`
}

// MLResponse là output từ FastAPI.
type MLResponse struct {
	Probability float64 `json:"probability"`
	Label       int     `json:"label"`
	RiskLevel   string  `json:"risk_level"`
}

// CreatePredictionRequest chứa dữ liệu đầu vào client gửi lên (giống MLRequest).
type CreatePredictionRequest struct {
	AgeYears    float64 `json:"age_years" binding:"required"`
	Gender      int     `json:"gender" binding:"required"`
	Height      float64 `json:"height" binding:"required"`
	Weight      float64 `json:"weight" binding:"required"`
	APHi        int     `json:"ap_hi" binding:"required"`
	APLo        int     `json:"ap_lo" binding:"required"`
	Cholesterol int     `json:"cholesterol" binding:"required"`
	Gluc        int     `json:"gluc" binding:"required"`
	Smoke       *int    `json:"smoke" binding:"required,oneof=0 1"`
	Alco        *int    `json:"alco" binding:"required,oneof=0 1"`
	Active      *int    `json:"active" binding:"required,oneof=0 1"`
}

type PredictionResponse struct {
	ID          string          `json:"id"`
	PatientID   string          `json:"patient_id"`
	Probability float64         `json:"probability"`
	RiskLabel   string          `json:"risk_label"`
	RawFeatures json.RawMessage `json:"raw_features,omitempty"`
	CreatedAt   time.Time       `json:"created_at"`
}

type ListPredictionsResponse struct {
	Predictions []PredictionResponse `json:"predictions"`
}

// Recommendation response struct
type RecommendationItem struct {
	Name        string `json:"name"`
	Intensity   string `json:"intensity"`
	DurationMin int    `json:"duration_min"`
	FreqPerWeek int    `json:"freq_per_week"`
	Notes       string `json:"notes,omitempty"`
}

type RecommendationPlan struct {
	Summary     string               `json:"summary"`
	Items       []RecommendationItem `json:"items,omitempty"`
	TemplateIDs []int64              `json:"template_ids,omitempty"`
}

type CreatePredictionResponse struct {
	Prediction     PredictionResponse `json:"prediction"`
	Recommendation RecommendationPlan `json:"recommendation"`
}
type ListPredictionsParams struct {
	Limit  int32 `form:"limit,default=10"`
	Offset int32 `form:"offset,default=0"`
}
