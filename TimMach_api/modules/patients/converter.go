package patients

import (
	"strconv"
	"time"

	db "chidinh/db/sqlc"

	"github.com/jackc/pgx/v5/pgtype"
)

// Patient represents the domain object independent from transport/DB types.
type Patient struct {
	ID               int64
	UserID           string
	Name             string
	Gender           int16
	Dob              time.Time
	CreatedAt        time.Time
	LatestPrediction *PredictionSummary
}

// PredictionSummary keeps only the last prediction needed for patient views.
type PredictionSummary struct {
	Probability float64
	RiskLabel   string
	CreatedAt   time.Time
}

func toPatientDomain(p db.Patient, latest *db.Prediction) Patient {
	return Patient{
		ID:               p.ID,
		UserID:           p.UserID,
		Name:             p.Name,
		Gender:           p.Gender,
		Dob:              safeDate(p.Dob),
		CreatedAt:        safeTime(p.CreatedAt),
		LatestPrediction: toPredictionSummaryDomain(latest),
	}
}

func toPatientDomainFromJoined(row db.ListPatientsWithLatestPredictionRow) Patient {
	return Patient{
		ID:               row.ID,
		UserID:           row.UserID,
		Name:             row.Name,
		Gender:           row.Gender,
		Dob:              safeDate(row.Dob),
		CreatedAt:        safeTime(row.CreatedAt),
		LatestPrediction: toPredictionSummaryFromJoined(row),
	}
}

func toPatientResponse(p Patient) PatientResponse {
	return PatientResponse{
		ID:               strconv.FormatInt(p.ID, 10),
		UserID:           p.UserID,
		Name:             p.Name,
		Gender:           p.Gender,
		Dob:              formatDate(p.Dob),
		CreatedAt:        p.CreatedAt,
		LatestPrediction: toPredictionSummaryDTO(p.LatestPrediction),
	}
}

func toPredictionSummaryDomain(pred *db.Prediction) *PredictionSummary {
	if pred == nil {
		return nil
	}
	return &PredictionSummary{
		Probability: pred.Probability,
		RiskLabel:   pred.RiskLabel,
		CreatedAt:   safeTime(pred.CreatedAt),
	}
}

func toPredictionSummaryFromJoined(row db.ListPatientsWithLatestPredictionRow) *PredictionSummary {
	if row.LatestProbability == nil && row.LatestRiskLabel == nil {
		return nil
	}

	summary := PredictionSummary{
		Probability: 0,
		RiskLabel:   "",
		CreatedAt:   safeTime(row.LatestPredictionAt),
	}
	if row.LatestProbability != nil {
		summary.Probability = *row.LatestProbability
	}
	if row.LatestRiskLabel != nil {
		summary.RiskLabel = *row.LatestRiskLabel
	}

	return &summary
}

func toPredictionSummaryDTO(pred *PredictionSummary) *PatientPredictionSummary {
	if pred == nil {
		return nil
	}
	return &PatientPredictionSummary{
		Probability: pred.Probability,
		RiskLabel:   pred.RiskLabel,
		CreatedAt:   pred.CreatedAt,
	}
}

func parseDate(input string) (time.Time, error) {
	return time.Parse("2006-01-02", input)
}

func toDBDate(t time.Time) pgtype.Date {
	return pgtype.Date{Time: t, Valid: true}
}

func formatDate(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format("2006-01-02")
}

func safeTime(t pgtype.Timestamptz) time.Time {
	if !t.Valid {
		return time.Time{}
	}
	return t.Time
}

func safeDate(d pgtype.Date) time.Time {
	if !d.Valid {
		return time.Time{}
	}
	return d.Time
}
