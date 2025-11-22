package patients

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"math"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	db "chidinh/db/sqlc"
	"chidinh/modules/predictions"
	"chidinh/utils"

	"github.com/SebastiaanKlippert/go-wkhtmltopdf"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

type PatientInfoView struct {
	Name     string
	DOB      string
	Gender   string
	HeightCm float64
	WeightKg float64
	BMI      float64
	AgeYears int
}

type PredictionView struct {
	Time            string
	ProbabilityPct  int
	RiskLevel       string
	Label           string
	MainRiskFactors []string
}

type ExerciseSessionView struct {
	Name        string
	DurationMin int
	Notes       string
	Frequency   string
}

type ExercisePlanView struct {
	Summary  string
	Sessions []ExerciseSessionView
}

type HistoryItemView struct {
	Time           string
	ProbabilityPct int
	RiskLabel      string
}

type PatientReportViewModel struct {
	ClinicName       string
	GeneratedAt      string
	Patient          PatientInfoView
	LatestPrediction *PredictionView
	ExercisePlan     ExercisePlanView
	History          []HistoryItemView
}

// GET /patients/:id/report.pdf
func (h *Handler) GetPatientReport(c *gin.Context) {
	userID, ok := utils.UserIDFromContext(c)
	if !ok {
		utils.RespondError(c, http.StatusUnauthorized, "missing user in context")
		return
	}

	patientID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.RespondError(c, http.StatusBadRequest, "invalid patient id")
		return
	}

	vm, err := buildPatientReportViewModel(c, h.Queries, userID, int64(patientID))
	if err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, pgx.ErrNoRows) {
			status = http.StatusNotFound
		}
		utils.RespondError(c, status, err.Error())
		return
	}

	tplPath := filepath.Join("templates", "patient_report.html")
	tpl, err := template.ParseFiles(tplPath)
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, fmt.Sprintf("cannot load template: %v", err))
		return
	}

	var htmlBuf bytes.Buffer
	if err := tpl.Execute(&htmlBuf, vm); err != nil {
		utils.RespondError(c, http.StatusInternalServerError, fmt.Sprintf("cannot render template: %v", err))
		return
	}

	pdfg, err := wkhtmltopdf.NewPDFGenerator()
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, fmt.Sprintf("cannot init pdf generator: %v", err))
		return
	}
	pdfg.AddPage(wkhtmltopdf.NewPageReader(bytes.NewReader(htmlBuf.Bytes())))
	pdfg.Orientation.Set(wkhtmltopdf.OrientationPortrait)
	pdfg.PageSize.Set(wkhtmltopdf.PageSizeA4)
	pdfg.Dpi.Set(150)
	pdfg.MarginLeft.Set(10)
	pdfg.MarginRight.Set(10)
	pdfg.MarginTop.Set(10)
	pdfg.MarginBottom.Set(10)

	if err := pdfg.Create(); err != nil {
		utils.RespondError(c, http.StatusInternalServerError, fmt.Sprintf("cannot generate pdf: %v", err))
		return
	}

	filename := fmt.Sprintf("patient_%d_report.pdf", patientID)
	c.Header("Content-Type", "application/pdf")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
	c.Data(http.StatusOK, "application/pdf", pdfg.Bytes())
}

func buildPatientReportViewModel(ctx context.Context, q *db.Queries, userID string, patientID int64) (PatientReportViewModel, error) {
	now := time.Now()
	vm := PatientReportViewModel{
		ClinicName:  "HeartCare Clinic",
		GeneratedAt: now.Format("2006-01-02 15:04"),
	}

	patient, err := q.GetPatientByID(ctx, patientID)
	if err != nil {
		return vm, err
	}
	if patient.UserID != userID {
		return vm, fmt.Errorf("patient does not belong to user")
	}

	vm.Patient = mapPatientInfo(patient)

	latestPred, err := q.GetLatestPredictionByPatient(ctx, patientID)
	if err == nil {
		vm.LatestPrediction = mapPredictionView(latestPred)
		vm.Patient = mergeFeaturesIntoPatient(vm.Patient, latestPred.RawFeatures)
	}

	historyRows, err := q.ListPredictionsByPatient(ctx, db.ListPredictionsByPatientParams{
		PatientID: patientID,
		Limit:     50,
		Offset:    0,
	})
	if err == nil {
		vm.History = mapHistory(historyRows)
	}

	recPlan, err := loadRecommendationPlan(ctx, q, patientID, latestPred)
	if err == nil {
		vm.ExercisePlan = recPlan
	}

	return vm, nil
}

func loadRecommendationPlan(ctx context.Context, q *db.Queries, patientID int64, latestPred db.Prediction) (ExercisePlanView, error) {
	var plan ExercisePlanView

	var rec db.ExerciseRecommendation
	var err error
	if latestPred.ID != 0 {
		rec, err = q.GetExerciseRecommendationByPrediction(ctx, latestPred.ID)
	} else {
		// fallback: lấy bản ghi mới nhất của bệnh nhân
		items, listErr := q.ListExerciseRecommendationsByPatient(ctx, db.ListExerciseRecommendationsByPatientParams{
			PatientID: patientID,
			Limit:     1,
			Offset:    0,
		})
		if listErr == nil && len(items) > 0 {
			rec = items[0]
			err = nil
		} else {
			err = listErr
		}
	}
	if err != nil {
		return plan, nil
	}

	var stored predictions.RecommendationPlan
	if len(rec.Plan) > 0 {
		_ = json.Unmarshal(rec.Plan, &stored)
	}

	templates, tplErr := q.ListExerciseTemplates(ctx)
	if tplErr != nil {
		return plan, tplErr
	}
	tplByID := make(map[int64]db.ExerciseTemplate, len(templates))
	for _, t := range templates {
		tplByID[t.ID] = t
	}

	plan.Summary = stored.Summary
	for _, id := range stored.TemplateIDs {
		if t, ok := tplByID[id]; ok {
			plan.Sessions = append(plan.Sessions, ExerciseSessionView{
				Name:        t.Name,
				DurationMin: int(t.DurationMin),
				Notes:       t.Description,
				Frequency:   fmt.Sprintf("%d buoi/tuần", t.FreqPerWeek),
			})
		}
	}

	return plan, nil
}

func mapPatientInfo(p db.Patient) PatientInfoView {
	age := 0
	if p.Dob.Valid {
		diff := time.Since(p.Dob.Time)
		age = int(diff.Hours() / 24 / 365)
	}
	return PatientInfoView{
		Name:     p.Name,
		DOB:      p.Dob.Time.Format("2006-01-02"),
		Gender:   genderLabel(p.Gender),
		HeightCm: 0,
		WeightKg: 0,
		BMI:      0,
		AgeYears: age,
	}
}

func mapPredictionView(pred db.Prediction) *PredictionView {
	rawFactors := extractRiskFactors(pred.RawFeatures)

	return &PredictionView{
		Time:            pred.CreatedAt.Time.Format("2006-01-02 15:04"),
		ProbabilityPct:  int(pred.Probability * 100),
		RiskLevel:       pred.RiskLabel,
		Label:           riskLabel(pred.RiskLabel),
		MainRiskFactors: rawFactors,
	}
}

func mapHistory(items []db.Prediction) []HistoryItemView {
	out := make([]HistoryItemView, 0, len(items))
	for _, p := range items {
		out = append(out, HistoryItemView{
			Time:           p.CreatedAt.Time.Format("2006-01-02 15:04"),
			ProbabilityPct: int(p.Probability * 100),
			RiskLabel:      riskLabel(p.RiskLabel),
		})
	}
	return out
}

type featuresPayload struct {
	AgeYears    float64 `json:"age_years"`
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

func extractRiskFactors(raw []byte) []string {
	var f featuresPayload
	_ = json.Unmarshal(raw, &f)

	factors := []string{}
	if f.Cholesterol >= 3 {
		factors = append(factors, "Cholesterol cao")
	}
	if f.Gluc >= 3 {
		factors = append(factors, "Đường huyết cao")
	}
	if f.APHi > 140 || f.APLo > 90 {
		factors = append(factors, "Huyết áp cao")
	}
	if f.Smoke == 1 {
		factors = append(factors, "Hút thuốc")
	}
	if f.Alco == 1 {
		factors = append(factors, "Uống rượu")
	}
	if f.Active == 0 {
		factors = append(factors, "Ít vận động")
	}
	return factors
}

func mergeFeaturesIntoPatient(info PatientInfoView, raw []byte) PatientInfoView {
	var f featuresPayload
	if err := json.Unmarshal(raw, &f); err != nil {
		return info
	}
	if f.Height > 0 {
		info.HeightCm = f.Height
	}
	if f.Weight > 0 {
		info.WeightKg = f.Weight
	}
	if f.Height > 0 && f.Weight > 0 {
		hm := f.Height / 100
		info.BMI = math.Round((f.Weight/(hm*hm))*10) / 10
	}
	if f.AgeYears > 0 {
		info.AgeYears = int(math.Round(f.AgeYears))
	}
	return info
}

func genderLabel(g int16) string {
	switch g {
	case 1:
		return "Male"
	case 2:
		return "Female"
	default:
		return "Other"
	}
}

func riskLabel(r string) string {
	switch strings.ToLower(r) {
	case "high":
		return "High"
	case "medium":
		return "Medium"
	case "low":
		return "Low"
	default:
		return r
	}
}
