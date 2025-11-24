package reports

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"math"
	"net/http"
	"net/mail"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	db "chidinh/db/sqlc"
	"chidinh/modules/predictions"
	"chidinh/utils"
	"chidinh/utils/mailer"

	"github.com/SebastiaanKlippert/go-wkhtmltopdf"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

const reportStorageDir = "tmp/reports"

type Controller struct {
	Queries            *db.Queries
	Mailer             *mailer.Mailer
	DefaultReportEmail string
}

func NewController(queries *db.Queries, mailerSvc *mailer.Mailer, defaultReportEmail string) *Controller {
	return &Controller{
		Queries:            queries,
		Mailer:             mailerSvc,
		DefaultReportEmail: defaultReportEmail,
	}
}

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

var errPatientForbidden = errors.New("patient does not belong to user")

// GetPatientReport tải xuống báo cáo PDF của bệnh nhân (legacy endpoint).
// GET /patients/:id/report.pdf
// Generate PDF on-the-fly và trả về dưới dạng file attachment.
func (h *Controller) GetPatientReport(c *gin.Context) {
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

	filename, pdfBytes, err := h.buildPatientReportPDF(c, userID, int64(patientID))
	if err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, pgx.ErrNoRows) {
			status = http.StatusNotFound
		} else if errors.Is(err, errPatientForbidden) {
			status = http.StatusForbidden
		}
		utils.RespondError(c, status, err.Error())
		return
	}

	c.Header("Content-Type", "application/pdf")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
	c.Data(http.StatusOK, "application/pdf", pdfBytes)
}

// SendPatientReportEmail gửi báo cáo PDF qua email (legacy endpoint).
// POST /patients/:id/report/email
// Generate PDF on-the-fly và gửi trực tiếp qua email, không lưu vào database.
func (h *Controller) SendPatientReportEmail(c *gin.Context) {
	userID, ok := utils.UserIDFromContext(c)
	if !ok {
		utils.RespondError(c, http.StatusUnauthorized, "missing user in context")
		return
	}

	if h.Mailer == nil || !h.Mailer.Enabled() {
		utils.RespondError(c, http.StatusInternalServerError, "email service is not configured")
		return
	}

	patientID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.RespondError(c, http.StatusBadRequest, "invalid patient id")
		return
	}

	var req SendReportEmailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondError(c, http.StatusBadRequest, "invalid request body")
		return
	}

	email := strings.TrimSpace(req.Email)
	if email == "" {
		email = h.DefaultReportEmail
	}
	if email == "" {
		utils.RespondError(c, http.StatusBadRequest, "email is required")
		return
	}
	if _, err := mail.ParseAddress(email); err != nil {
		utils.RespondError(c, http.StatusBadRequest, "email is invalid")
		return
	}

	filename, pdfBytes, err := h.buildPatientReportPDF(c, userID, int64(patientID))
	if err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, pgx.ErrNoRows) {
			status = http.StatusNotFound
		} else if errors.Is(err, errPatientForbidden) {
			status = http.StatusForbidden
		}
		utils.RespondError(c, status, err.Error())
		return
	}

	subject := req.Subject
	if subject == "" {
		subject = "Báo cáo kết quả tim mạch"
	}
	body := req.Message
	if body == "" {
		body = "Đính kèm báo cáo kết quả tim mạch của bạn."
	}

	err = h.Mailer.Send(email, subject, body, []mailer.Attachment{{
		Filename: filename,
		MimeType: "application/pdf",
		Content:  pdfBytes,
	}})
	if err != nil {
		if errors.Is(err, mailer.ErrNotConfigured) {
			utils.RespondError(c, http.StatusInternalServerError, "email service is not configured")
			return
		}
		utils.RespondError(c, http.StatusInternalServerError, fmt.Sprintf("cannot send email: %v", err))
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "sent"})
}

// CreateReport tạo báo cáo mới và lưu vào database.
// POST /patients/:id/reports
// Workflow:
// 1. Validate user ownership
// 2. Generate PDF
// 3. Lưu file vào disk (tmp/reports/{patient_id}/{filename}.pdf)
// 4. Insert record vào bảng reports
// 5. Trả về ReportResponse với id, filename, file_url
func (h *Controller) CreateReport(c *gin.Context) {
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

	_, err = h.getOwnedPatient(c, userID, int64(patientID))
	if err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, pgx.ErrNoRows) {
			status = http.StatusNotFound
		} else if errors.Is(err, errPatientForbidden) {
			status = http.StatusForbidden
		}
		utils.RespondError(c, status, err.Error())
		return
	}

	_, pdfBytes, err := h.buildPatientReportPDF(c, userID, int64(patientID))
	if err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, pgx.ErrNoRows) {
			status = http.StatusNotFound
		} else if errors.Is(err, errPatientForbidden) {
			status = http.StatusForbidden
		}
		utils.RespondError(c, status, err.Error())
		return
	}

	filename := fmt.Sprintf("report_%d_%s.pdf", patientID, time.Now().Format("20060102_150405"))
	filePath, err := saveReportFile(int64(patientID), filename, pdfBytes)
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, fmt.Sprintf("cannot save report file: %v", err))
		return
	}

	recipientsJSON, err := json.Marshal([]ReportRecipient{})
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "cannot initialize recipients")
		return
	}

	report, err := h.Queries.CreateReport(c, db.CreateReportParams{
		PatientID:  int64(patientID),
		Filename:   filename,
		FileUrl:    filePath,
		Recipients: recipientsJSON,
	})
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "cannot create report record")
		return
	}

	resp, err := mapReportResponse(report)
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "cannot build response")
		return
	}

	c.JSON(http.StatusCreated, resp)
}

// ListReports lấy danh sách báo cáo của bệnh nhân.
// GET /patients/:id/reports?limit=10&offset=0
// Trả về danh sách reports với thông tin recipients đã gửi email.
func (h *Controller) ListReports(c *gin.Context) {
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

	_, err = h.getOwnedPatient(c, userID, int64(patientID))
	if err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, pgx.ErrNoRows) {
			status = http.StatusNotFound
		} else if errors.Is(err, errPatientForbidden) {
			status = http.StatusForbidden
		}
		utils.RespondError(c, status, err.Error())
		return
	}

	var req ListReportsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		utils.RespondError(c, http.StatusBadRequest, "invalid query parameters")
		return
	}
	limit := req.Limit
	if limit == 0 {
		limit = 10
	}

	items, err := h.Queries.ListReportsByPatient(c, db.ListReportsByPatientParams{
		PatientID: int64(patientID),
		Limit:     limit,
		Offset:    req.Offset,
	})
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "cannot list reports")
		return
	}

	total, err := h.Queries.CountReportsByPatient(c, int64(patientID))
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "cannot count reports")
		return
	}

	resp := ListReportsResponse{Reports: make([]ReportResponse, 0, len(items)), Total: total}
	for _, r := range items {
		item, mapErr := mapReportResponse(r)
		if mapErr != nil {
			utils.RespondError(c, http.StatusInternalServerError, "cannot build response")
			return
		}
		resp.Reports = append(resp.Reports, item)
	}

	c.JSON(http.StatusOK, resp)
}

// DownloadReport tải xuống file PDF báo cáo đã được lưu.
// GET /reports/:id/download
// Đọc file từ disk dựa trên file_url trong database và trả về dưới dạng attachment.
func (h *Controller) DownloadReport(c *gin.Context) {
	userID, ok := utils.UserIDFromContext(c)
	if !ok {
		utils.RespondError(c, http.StatusUnauthorized, "missing user in context")
		return
	}

	reportID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.RespondError(c, http.StatusBadRequest, "invalid report id")
		return
	}

	report, _, err := h.getOwnedReport(c, userID, int64(reportID))
	if err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, pgx.ErrNoRows) {
			status = http.StatusNotFound
		} else if errors.Is(err, errPatientForbidden) {
			status = http.StatusForbidden
		}
		utils.RespondError(c, status, err.Error())
		return
	}

	fileBytes, err := os.ReadFile(report.FileUrl)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			utils.RespondError(c, http.StatusNotFound, "report file not found")
			return
		}
		utils.RespondError(c, http.StatusInternalServerError, fmt.Sprintf("cannot read file: %v", err))
		return
	}

	c.Header("Content-Type", "application/pdf")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", report.Filename))
	c.Data(http.StatusOK, "application/pdf", fileBytes)
}

// SendReportEmail gửi báo cáo đã lưu qua email.
// POST /reports/:id/email
// Workflow:
// 1. Lấy report từ database
// 2. Validate ownership
// 3. Đọc PDF file từ disk
// 4. Gửi email qua mailer service
// 5. Cập nhật recipients JSONB trong database (track email đã gửi)
func (h *Controller) SendReportEmail(c *gin.Context) {
	userID, ok := utils.UserIDFromContext(c)
	if !ok {
		utils.RespondError(c, http.StatusUnauthorized, "missing user in context")
		return
	}

	if h.Mailer == nil || !h.Mailer.Enabled() {
		utils.RespondError(c, http.StatusInternalServerError, "email service is not configured")
		return
	}

	reportID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.RespondError(c, http.StatusBadRequest, "invalid report id")
		return
	}

	var req SendReportEmailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondError(c, http.StatusBadRequest, "invalid request body")
		return
	}

	email := strings.TrimSpace(req.Email)
	if email == "" {
		utils.RespondError(c, http.StatusBadRequest, "email is required")
		return
	}
	if _, err := mail.ParseAddress(email); err != nil {
		utils.RespondError(c, http.StatusBadRequest, "email is invalid")
		return
	}

	report, _, err := h.getOwnedReport(c, userID, int64(reportID))
	if err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, pgx.ErrNoRows) {
			status = http.StatusNotFound
		} else if errors.Is(err, errPatientForbidden) {
			status = http.StatusForbidden
		}
		utils.RespondError(c, status, err.Error())
		return
	}

	fileBytes, err := os.ReadFile(report.FileUrl)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			utils.RespondError(c, http.StatusNotFound, "report file not found")
			return
		}
		utils.RespondError(c, http.StatusInternalServerError, fmt.Sprintf("cannot read file: %v", err))
		return
	}

	subject := req.Subject
	if subject == "" {
		subject = "Báo cáo kết quả tim mạch"
	}
	body := req.Message
	if body == "" {
		body = "Đính kèm báo cáo kết quả tim mạch của bạn."
	}

	err = h.Mailer.Send(email, subject, body, []mailer.Attachment{{
		Filename: report.Filename,
		MimeType: "application/pdf",
		Content:  fileBytes,
	}})
	if err != nil {
		if errors.Is(err, mailer.ErrNotConfigured) {
			utils.RespondError(c, http.StatusInternalServerError, "email service is not configured")
			return
		}
		utils.RespondError(c, http.StatusInternalServerError, fmt.Sprintf("cannot send email: %v", err))
		return
	}

	recipients, err := decodeRecipients(report.Recipients)
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "cannot parse recipients")
		return
	}
	recipients = append(recipients, ReportRecipient{
		Email:  email,
		SentAt: time.Now(),
		Status: "sent",
	})
	recipientsJSON, err := json.Marshal(recipients)
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "cannot update recipients")
		return
	}

	_, err = h.Queries.UpdateReportRecipients(c, db.UpdateReportRecipientsParams{
		ID:         report.ID,
		Recipients: recipientsJSON,
	})
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "cannot persist recipients")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":   "sent",
		"report_id": report.ID,
		"sent_to":   email,
	})
}

// DeleteReport xóa báo cáo (cả file và database record).
// DELETE /reports/:id
// Xóa file PDF từ disk và xóa record khỏi bảng reports.
func (h *Controller) DeleteReport(c *gin.Context) {
	userID, ok := utils.UserIDFromContext(c)
	if !ok {
		utils.RespondError(c, http.StatusUnauthorized, "missing user in context")
		return
	}

	reportID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.RespondError(c, http.StatusBadRequest, "invalid report id")
		return
	}

	report, _, err := h.getOwnedReport(c, userID, int64(reportID))
	if err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, pgx.ErrNoRows) {
			status = http.StatusNotFound
		} else if errors.Is(err, errPatientForbidden) {
			status = http.StatusForbidden
		}
		utils.RespondError(c, status, err.Error())
		return
	}

	if err := os.Remove(report.FileUrl); err != nil && !errors.Is(err, os.ErrNotExist) {
		utils.RespondError(c, http.StatusInternalServerError, fmt.Sprintf("cannot delete file: %v", err))
		return
	}

	if err := h.Queries.DeleteReport(c, report.ID); err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "cannot delete report")
		return
	}

	c.Status(http.StatusNoContent)
}

// buildPatientReportPDF tạo file PDF báo cáo từ HTML template.
// Trả về: (filename, pdfBytes, error)
// Sử dụng wkhtmltopdf để convert HTML → PDF.
func (h *Controller) buildPatientReportPDF(ctx context.Context, userID string, patientID int64) (string, []byte, error) {
	vm, err := buildPatientReportViewModel(ctx, h.Queries, userID, patientID)
	if err != nil {
		return "", nil, err
	}

	tplPath := filepath.Join("modules", "reports", "templates", "patient_report.html")
	tpl, err := template.ParseFiles(tplPath)
	if err != nil {
		return "", nil, fmt.Errorf("cannot load template: %w", err)
	}

	var htmlBuf bytes.Buffer
	if err := tpl.Execute(&htmlBuf, vm); err != nil {
		return "", nil, fmt.Errorf("cannot render template: %w", err)
	}

	pdfg, err := wkhtmltopdf.NewPDFGenerator()
	if err != nil {
		return "", nil, fmt.Errorf("cannot init pdf generator: %w", err)
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
		return "", nil, fmt.Errorf("cannot generate pdf: %w", err)
	}

	filename := fmt.Sprintf("patient_%d_report.pdf", patientID)
	return filename, pdfg.Bytes(), nil
}

// buildPatientReportViewModel tập hợp dữ liệu cho báo cáo PDF.
// Lấy thông tin từ:
// - patients (name, dob, gender)
// - predictions (latest prediction + history)
// - raw_features (height, weight, BMI từ prediction gần nhất)
// - exercise_recommendations (exercise plan)
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
		return vm, errPatientForbidden
	}

	vm.Patient = mapPatientInfo(patient)

	// Lấy prediction gần nhất (nếu có)
	latestPred, err := q.GetLatestPredictionByPatient(ctx, patientID)
	if err == nil {
		features, factors := decodeFeatures(latestPred.RawFeatures, latestPred.Factors)
		// Map prediction cho view (risk, probability, factors từ ML)
		vm.LatestPrediction = mapPredictionView(latestPred, factors)
		// Merge health metrics từ raw_features vào patient info (height, weight, BMI)
		vm.Patient = mergeFeaturesIntoPatient(vm.Patient, features)
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

// loadRecommendationPlan tải exercise plan từ exercise_recommendations.
// Lấy recommendation plan dựa trên prediction_id hoặc patient_id.
// Resolve template details từ exercise_templates table.
func loadRecommendationPlan(ctx context.Context, q *db.Queries, patientID int64, latestPred db.Prediction) (ExercisePlanView, error) {
	var plan ExercisePlanView

	var rec db.ExerciseRecommendation
	var err error
	if latestPred.ID != 0 {
		rec, err = q.GetExerciseRecommendationByPrediction(ctx, latestPred.ID)
	} else {
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

// mapPatientInfo convert db.Patient → PatientInfoView.
// Tính tuổi từ DOB (date of birth).
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

// mapPredictionView convert db.Prediction → PredictionView.
// Extract risk factors từ raw_features JSONB.
func mapPredictionView(pred db.Prediction, factors []predictions.RiskFactor) *PredictionView {
	mainFactors := make([]string, 0, len(factors))
	for _, f := range factors {
		if f.Message != "" {
			mainFactors = append(mainFactors, f.Message)
		}
	}
	return &PredictionView{
		Time:            pred.CreatedAt.Time.Format("2006-01-02 15:04"),
		ProbabilityPct:  int(pred.Probability * 100),
		RiskLevel:       pred.RiskLabel,
		Label:           riskLabel(pred.RiskLabel),
		MainRiskFactors: mainFactors,
	}
}

// mapHistory convert danh sách predictions → history items cho báo cáo.
// Chỉ hiển thị thông tin cơ bản: time, probability, risk label.
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

// decodeFeatures đọc raw_features + factors để lấy input và factors do ML trả về.
// Hỗ trợ cả định dạng mới {input:{...}, factors:[...]} và định dạng cũ (fields phẳng).
func decodeFeatures(raw []byte, factorsJSON []byte) (featuresPayload, []predictions.RiskFactor) {
	var out featuresPayload
	factors := decodeFactors(factorsJSON)

	var container struct {
		Input   featuresPayload          `json:"input"`
		Factors []predictions.RiskFactor `json:"factors"`
	}
	if err := json.Unmarshal(raw, &container); err == nil {
		if container.Input != (featuresPayload{}) {
			out = container.Input
		}
		factors = container.Factors
	}
	if out == (featuresPayload{}) {
		_ = json.Unmarshal(raw, &out)
	}
	return out, factors
}

func decodeFactors(raw []byte) []predictions.RiskFactor {
	if len(raw) == 0 {
		return nil
	}
	var factors []predictions.RiskFactor
	if err := json.Unmarshal(raw, &factors); err != nil {
		return nil
	}
	return factors
}

// mergeFeaturesIntoPatient merge health metrics từ raw_features vào patient info.
// Lấy height, weight từ prediction để tính BMI.
// Update age từ age_years nếu chính xác hơn tính từ DOB.
func mergeFeaturesIntoPatient(info PatientInfoView, f featuresPayload) PatientInfoView {
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

// saveReportFile lưu file PDF vào disk.
// Cấu trúc: tmp/reports/{patient_id}/{filename}.pdf
// Trả về file path để lưu vào database.
func saveReportFile(patientID int64, filename string, content []byte) (string, error) {
	dir := filepath.Join(reportStorageDir, fmt.Sprintf("%d", patientID))
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", err
	}
	path := filepath.Join(dir, filename)
	if err := os.WriteFile(path, content, 0644); err != nil {
		return "", err
	}
	return path, nil
}

// getOwnedPatient validate rằng patient thuộc về user.
// Trả về error nếu patient không tồn tại hoặc không thuộc về user.
func (h *Controller) getOwnedPatient(ctx context.Context, userID string, patientID int64) (db.Patient, error) {
	patient, err := h.Queries.GetPatientByID(ctx, patientID)
	if err != nil {
		return patient, err
	}
	if patient.UserID != userID {
		return patient, errPatientForbidden
	}
	return patient, nil
}

// getOwnedReport validate rằng report thuộc về user (thông qua patient ownership).
// Trả về report và patient nếu hợp lệ.
func (h *Controller) getOwnedReport(ctx context.Context, userID string, reportID int64) (db.Report, db.Patient, error) {
	report, err := h.Queries.GetReportByID(ctx, reportID)
	if err != nil {
		return report, db.Patient{}, err
	}
	patient, err := h.getOwnedPatient(ctx, userID, report.PatientID)
	if err != nil {
		return report, patient, err
	}
	return report, patient, nil
}
