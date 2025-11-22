package exercises

import (
	"encoding/json"
	"strconv"
	"time"

	db "chidinh/db/sqlc"
	"chidinh/modules/predictions"

	"github.com/jackc/pgx/v5/pgtype"
)

// Template represents a domain exercise template.
type Template struct {
	ID              int64
	Name            string
	Intensity       string
	Description     string
	DurationMin     int32
	FreqPerWeek     int32
	TargetRiskLevel string
	Tags            []string
}

// Recommendation represents a domain recommendation with hydrated plan.
type Recommendation struct {
	ID           int64
	PatientID    int64
	PredictionID int64
	Plan         *predictions.RecommendationPlan
	CreatedAt    time.Time
}

func toTemplateDomain(t db.ExerciseTemplate) Template {
	return Template{
		ID:              t.ID,
		Name:            t.Name,
		Intensity:       t.Intensity,
		Description:     t.Description,
		DurationMin:     t.DurationMin,
		FreqPerWeek:     t.FreqPerWeek,
		TargetRiskLevel: t.TargetRiskLevel,
		Tags:            t.Tags,
	}
}

func toTemplateResponse(t Template) TemplateResponse {
	return TemplateResponse{
		ID:              t.ID,
		Name:            t.Name,
		Intensity:       t.Intensity,
		Description:     t.Description,
		DurationMin:     t.DurationMin,
		FreqPerWeek:     t.FreqPerWeek,
		TargetRiskLevel: t.TargetRiskLevel,
		Tags:            t.Tags,
	}
}

func toRecommendationDomain(r db.ExerciseRecommendation, tplByID map[int64]db.ExerciseTemplate) Recommendation {
	var stored predictions.RecommendationPlan
	if len(r.Plan) > 0 {
		_ = json.Unmarshal(r.Plan, &stored)
	}

	return Recommendation{
		ID:           r.ID,
		PatientID:    r.PatientID,
		PredictionID: r.PredictionID,
		Plan:         hydratePlan(stored, tplByID),
		CreatedAt:    safeTime(r.CreatedAt),
	}
}

func toRecommendationResponse(r Recommendation) RecommendationResponse {
	return RecommendationResponse{
		ID:           r.ID,
		PatientID:    strconv.FormatInt(r.PatientID, 10),
		PredictionID: strconv.FormatInt(r.PredictionID, 10),
		Plan:         r.Plan,
		CreatedAt:    r.CreatedAt,
	}
}

func indexTemplates(templates []db.ExerciseTemplate) map[int64]db.ExerciseTemplate {
	tplByID := make(map[int64]db.ExerciseTemplate, len(templates))
	for _, t := range templates {
		tplByID[t.ID] = t
	}
	return tplByID
}

func hydratePlan(stored predictions.RecommendationPlan, tplByID map[int64]db.ExerciseTemplate) *predictions.RecommendationPlan {
	if stored.Summary == "" && len(stored.TemplateIDs) == 0 {
		return nil
	}

	resp := predictions.RecommendationPlan{
		Summary:     stored.Summary,
		TemplateIDs: stored.TemplateIDs,
		Items:       []predictions.RecommendationItem{},
	}

	for _, id := range stored.TemplateIDs {
		if t, ok := tplByID[id]; ok {
			resp.Items = append(resp.Items, predictions.RecommendationItem{
				Name:        t.Name,
				Intensity:   t.Intensity,
				DurationMin: int(t.DurationMin),
				FreqPerWeek: int(t.FreqPerWeek),
				Notes:       t.Description,
			})
		}
	}

	return &resp
}

func safeTime(t pgtype.Timestamptz) time.Time {
	if !t.Valid {
		return time.Time{}
	}
	return t.Time
}
