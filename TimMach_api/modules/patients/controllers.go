package patients

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	db "chidinh/db/sqlc"
	"chidinh/utils"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

// Handler gom dependency cho module patients.
type Handler struct {
	Queries *db.Queries
}

func NewHandler(queries *db.Queries) *Handler {
	return &Handler{Queries: queries}
}

// POST /patients
func (h *Handler) CreatePatient(c *gin.Context) {
	userID, ok := utils.UserIDFromContext(c)
	if !ok {
		utils.RespondError(c, http.StatusUnauthorized, "missing user in context")
		return
	}

	var req CreatePatientRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondError(c, http.StatusBadRequest, "invalid request body")
		return
	}

	dobVal, err := parseDate(req.Dob)
	if err != nil {
		utils.RespondError(c, http.StatusBadRequest, "dob must be YYYY-MM-DD")
		return
	}

	patient, err := h.Queries.CreatePatient(c, db.CreatePatientParams{
		UserID: userID,
		Name:   req.Name,
		Gender: req.Gender,
		Dob:    pgtype.Date{Time: dobVal, Valid: true},
	})
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "cannot create patient")
		return
	}

	c.JSON(http.StatusCreated, mapPatientResponse(patient, nil))
}

// GET /patients
func (h *Handler) ListPatients(c *gin.Context) {
	userID, ok := utils.UserIDFromContext(c)
	if !ok {
		utils.RespondError(c, http.StatusUnauthorized, "missing user in context")
		return
	}

	var req ListPatientsParams
	if err := c.ShouldBindQuery(&req); err != nil {
		utils.RespondError(c, http.StatusBadRequest, "invalid query parameters")
		return
	}

	riskFilter := strings.ToLower(strings.TrimSpace(req.Risk))
	switch riskFilter {
	case "", "low", "medium", "high", "none":
	default:
		utils.RespondError(c, http.StatusBadRequest, "risk must be low/medium/high/none")
		return
	}

	items, err := h.Queries.ListPatientsWithLatestPrediction(c, db.ListPatientsWithLatestPredictionParams{
		UserID:  userID,
		Limit:   req.Limit,
		Offset:  req.Offset,
		Column4: riskFilter,
	})
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "cannot list patients")
		return
	}

	resp := ListPatientsResponse{Patients: make([]PatientResponse, 0, len(items))}
	for _, p := range items {
		resp.Patients = append(resp.Patients, mapPatientFromJoinedRow(p))
	}

	c.JSON(http.StatusOK, resp)
}

// GET /patients/:id
func (h *Handler) GetPatient(c *gin.Context) {
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

	patient, err := h.Queries.GetPatientByID(c, int64(patientID))
	if errors.Is(err, pgx.ErrNoRows) {
		utils.RespondError(c, http.StatusNotFound, "patient not found")
		return
	}
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "cannot fetch patient")
		return
	}

	if patient.UserID != userID {
		utils.RespondError(c, http.StatusForbidden, "patient does not belong to user")
		return
	}

	var latest *db.Prediction
	if pred, err := h.Queries.GetLatestPredictionByPatient(c, int64(patientID)); err == nil {
		latest = &pred
	}

	c.JSON(http.StatusOK, mapPatientResponse(patient, latest))
}

// PUT/PATCH /patients/:id
func (h *Handler) UpdatePatient(c *gin.Context) {
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

	existing, err := h.Queries.GetPatientByID(c, int64(patientID))
	if errors.Is(err, pgx.ErrNoRows) {
		utils.RespondError(c, http.StatusNotFound, "patient not found")
		return
	}
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "cannot fetch patient")
		return
	}
	if existing.UserID != userID {
		utils.RespondError(c, http.StatusForbidden, "patient does not belong to user")
		return
	}

	var req UpdatePatientRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondError(c, http.StatusBadRequest, "invalid request body")
		return
	}

	newName := existing.Name
	if req.Name != nil {
		newName = *req.Name
	}

	newGender := existing.Gender
	if req.Gender != nil {
		newGender = *req.Gender
	}

	newDob := existing.Dob
	if req.Dob != nil {
		dob, err := parseDate(*req.Dob)
		if err != nil {
			utils.RespondError(c, http.StatusBadRequest, "dob must be YYYY-MM-DD")
			return
		}
		newDob = pgtype.Date{Time: dob, Valid: true}
	}

	updated, err := h.Queries.UpdatePatient(c, db.UpdatePatientParams{
		ID:     int64(patientID),
		Name:   newName,
		Gender: newGender,
		Dob:    pgtype.Date{Time: newDob.Time, Valid: newDob.Valid},
	})
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "cannot update patient")
		return
	}

	c.JSON(http.StatusOK, mapPatientResponse(updated, nil))
}

// DELETE /patients/:id
func (h *Handler) DeletePatient(c *gin.Context) {
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

	patient, err := h.Queries.GetPatientByID(c, int64(patientID))
	if errors.Is(err, pgx.ErrNoRows) {
		utils.RespondError(c, http.StatusNotFound, "patient not found")
		return
	}
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "cannot fetch patient")
		return
	}
	if patient.UserID != userID {
		utils.RespondError(c, http.StatusForbidden, "patient does not belong to user")
		return
	}

	if err := h.Queries.DeletePatient(c, int64(patientID)); err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "cannot delete patient")
		return
	}

	c.Status(http.StatusNoContent)
}

func mapPatientResponse(p db.Patient, latest *db.Prediction) PatientResponse {
	resp := PatientResponse{
		ID:        strconv.FormatInt(p.ID, 10),
		UserID:    p.UserID,
		Name:      p.Name,
		Gender:    p.Gender,
		Dob:       formatDate(p.Dob),
		CreatedAt: safeTime(p.CreatedAt),
	}
	resp.LatestPrediction = mapPredictionSummary(latest)
	return resp
}

func mapPatientFromJoinedRow(row db.ListPatientsWithLatestPredictionRow) PatientResponse {
	summary := mapPredictionSummaryFromJoined(row)
	return PatientResponse{
		ID:               strconv.FormatInt(row.ID, 10),
		UserID:           row.UserID,
		Name:             row.Name,
		Gender:           row.Gender,
		Dob:              formatDate(row.Dob),
		CreatedAt:        safeTime(row.CreatedAt),
		LatestPrediction: summary,
	}
}

func parseDate(input string) (time.Time, error) {
	return time.Parse("2006-01-02", input)
}

func formatDate(d pgtype.Date) string {
	if !d.Valid {
		return ""
	}
	return d.Time.Format("2006-01-02")
}

func safeTime(t pgtype.Timestamptz) time.Time {
	if !t.Valid {
		return time.Time{}
	}
	return t.Time
}

func mapPredictionSummary(pred *db.Prediction) *PatientPredictionSummary {
	if pred == nil {
		return nil
	}
	return &PatientPredictionSummary{
		Probability: pred.Probability,
		RiskLabel:   pred.RiskLabel,
		CreatedAt:   safeTime(pred.CreatedAt),
	}
}

func mapPredictionSummaryFromJoined(row db.ListPatientsWithLatestPredictionRow) *PatientPredictionSummary {
	if row.LatestProbability == nil && row.LatestRiskLabel == nil {
		return nil
	}

	summary := PatientPredictionSummary{
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
