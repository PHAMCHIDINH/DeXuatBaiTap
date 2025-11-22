package patients

import "time"

type CreatePatientRequest struct {
	Name   string `json:"name" binding:"required"`
	Gender int16  `json:"gender" binding:"required"`
	Dob    string `json:"dob" binding:"required"` // yyyy-mm-dd
}

type UpdatePatientRequest struct {
	Name   *string `json:"name,omitempty"`
	Gender *int16  `json:"gender,omitempty"`
	Dob    *string `json:"dob,omitempty"` // yyyy-mm-dd
}

type PatientPredictionSummary struct {
	Probability float64   `json:"probability"`
	RiskLabel   string    `json:"risk_label"`
	CreatedAt   time.Time `json:"created_at"`
}

type PatientResponse struct {
	ID                string                   `json:"id"`
	UserID            string                   `json:"user_id"`
	Name              string                   `json:"name"`
	Gender            int16                    `json:"gender"`
	Dob               string                   `json:"dob"`
	CreatedAt         time.Time                `json:"created_at"`
	LatestPrediction  *PatientPredictionSummary `json:"latest_prediction,omitempty"`
}

type ListPatientsResponse struct {
	Patients []PatientResponse `json:"patients"`
}
type ListPatientsParams struct {
	Limit  int32  `form:"limit,default=10"`
	Offset int32  `form:"offset,default=0"`
	Risk   string `form:"risk"`
}
