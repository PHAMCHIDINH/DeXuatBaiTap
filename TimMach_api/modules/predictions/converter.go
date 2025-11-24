package predictions

import (
	"encoding/json"
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
	Factors     []RiskFactor
	RawFeatures []byte
	CreatedAt   time.Time
}

func toMLRequest(req CreatePredictionRequest) MLRequest {
	// Chuyển payload client sang payload gửi cho ML (thay nil bằng 0).
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
	factors := decodeDBFactors(p.Factors)
	return Prediction{
		ID:          p.ID,
		PatientID:   p.PatientID,
		Probability: p.Probability,
		RiskLabel:   p.RiskLabel,
		Factors:     factors,
		RawFeatures: p.RawFeatures,
		CreatedAt:   safeTime(p.CreatedAt),
	}
}

func toPredictionResponse(p Prediction) PredictionResponse {
	// raw_features được giữ nguyên để FE có thể xem input gốc.
	rawJSON := json.RawMessage(p.RawFeatures)
	return PredictionResponse{
		ID:          strconv.FormatInt(p.ID, 10),
		PatientID:   strconv.FormatInt(p.PatientID, 10),
		Probability: p.Probability,
		RiskLabel:   p.RiskLabel,
		RawFeatures: rawJSON,
		Factors:     p.Factors,
		CreatedAt:   p.CreatedAt,
	}
}

func encodeStoredFeatures(input MLRequest) []byte {
	// Lưu input gốc (không nhồi factors) vào raw_features.
	raw, err := json.Marshal(input)
	if err != nil {
		return nil
	}
	return raw
}

func decodeDBFactors(raw []byte) []RiskFactor {
	if len(raw) == 0 {
		return nil
	}
	var factors []RiskFactor
	if err := json.Unmarshal(raw, &factors); err != nil {
		return nil
	}
	return factors
}

func safeTime(t pgtype.Timestamptz) time.Time {
	if !t.Valid {
		return time.Time{}
	}
	return t.Time
}
