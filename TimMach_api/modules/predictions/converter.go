package predictions

import (
	"strconv"
	"time"

	db "chidinh/db/sqlc"

	"github.com/jackc/pgx/v5/pgtype"
)

// Prediction represents domain model separated from DB driver types.
type Prediction struct {
	ID          int64
	PatientID   int64
	Probability float64
	RiskLabel   string
	RawFeatures []byte
	CreatedAt   time.Time
}

func toMLRequest(req CreatePredictionRequest) MLRequest {
	toInt := func(v *int) int {
		if v == nil {
			return 0
		}
		return *v
	}
	return MLRequest{
		AgeYears:    req.AgeYears,
		Gender:      req.Gender,
		Height:      req.Height,
		Weight:      req.Weight,
		APHi:        req.APHi,
		APLo:        req.APLo,
		Cholesterol: req.Cholesterol,
		Gluc:        req.Gluc,
		Smoke:       toInt(req.Smoke),
		Alco:        toInt(req.Alco),
		Active:      toInt(req.Active),
	}
}

func toPredictionDomain(p db.Prediction) Prediction {
	return Prediction{
		ID:          p.ID,
		PatientID:   p.PatientID,
		Probability: p.Probability,
		RiskLabel:   p.RiskLabel,
		RawFeatures: p.RawFeatures,
		CreatedAt:   safeTime(p.CreatedAt),
	}
}

func toPredictionResponse(p Prediction) PredictionResponse {
	return PredictionResponse{
		ID:          strconv.FormatInt(p.ID, 10),
		PatientID:   strconv.FormatInt(p.PatientID, 10),
		Probability: p.Probability,
		RiskLabel:   p.RiskLabel,
		RawFeatures: p.RawFeatures,
		CreatedAt:   p.CreatedAt,
	}
}

func safeTime(t pgtype.Timestamptz) time.Time {
	if !t.Valid {
		return time.Time{}
	}
	return t.Time
}
