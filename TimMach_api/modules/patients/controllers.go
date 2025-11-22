package patients

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	db "chidinh/db/sqlc"
	"chidinh/utils"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
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
		Dob:    toDBDate(dobVal),
	})
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "cannot create patient")
		return
	}

	c.JSON(http.StatusCreated, toPatientResponse(toPatientDomain(patient, nil)))
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
		resp.Patients = append(resp.Patients, toPatientResponse(toPatientDomainFromJoined(p)))
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

	c.JSON(http.StatusOK, toPatientResponse(toPatientDomain(patient, latest)))
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
		newDob = toDBDate(dob)
	}

	updated, err := h.Queries.UpdatePatient(c, db.UpdatePatientParams{
		ID:     int64(patientID),
		Name:   newName,
		Gender: newGender,
		Dob:    newDob,
	})
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "cannot update patient")
		return
	}

	c.JSON(http.StatusOK, toPatientResponse(toPatientDomain(updated, nil)))
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
