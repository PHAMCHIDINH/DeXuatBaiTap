package exercises

import (
	"net/http"
	"strconv"
	"strings"

	db "chidinh/db/sqlc"
	"chidinh/utils"

	"github.com/gin-gonic/gin"
)

// Handler gom dependencies cho module exercise templates/recommendations.
type Handler struct {
	Queries *db.Queries
}

func NewHandler(q *db.Queries) *Handler {
	return &Handler{Queries: q}
}

// POST /exercise-templates
func (h *Handler) CreateTemplate(c *gin.Context) {
	if _, ok := utils.UserIDFromContext(c); !ok {
		utils.RespondError(c, http.StatusUnauthorized, "missing user in context")
		return
	}

	var req CreateTemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondError(c, http.StatusBadRequest, "invalid request body")
		return
	}

	risk := strings.ToLower(strings.TrimSpace(req.TargetRiskLevel))
	if risk != "low" && risk != "medium" && risk != "high" && risk != "none" {
		utils.RespondError(c, http.StatusBadRequest, "target_risk_level must be low/medium/high/none")
		return
	}

	template, err := h.Queries.CreateExerciseTemplate(c, db.CreateExerciseTemplateParams{
		Name:           req.Name,
		Intensity:      req.Intensity,
		Description:    req.Description,
		DurationMin:    req.DurationMin,
		FreqPerWeek:    req.FreqPerWeek,
		TargetRiskLevel: risk,
		Tags:           req.Tags,
	})
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "cannot create template")
		return
	}

	c.JSON(http.StatusCreated, toTemplateResponse(toTemplateDomain(template)))
}

// GET /exercise-templates
func (h *Handler) ListTemplates(c *gin.Context) {
	if _, ok := utils.UserIDFromContext(c); !ok {
		utils.RespondError(c, http.StatusUnauthorized, "missing user in context")
		return
	}

	items, err := h.Queries.ListExerciseTemplates(c)
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "cannot list templates")
		return
	}
	resp := ListTemplatesResponse{Templates: make([]TemplateResponse, 0, len(items))}
	for _, t := range items {
		resp.Templates = append(resp.Templates, toTemplateResponse(toTemplateDomain(t)))
	}
	c.JSON(http.StatusOK, resp)
}

// GET /patients/:id/recommendations
func (h *Handler) ListRecommendationsByPatient(c *gin.Context) {
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
	if err != nil {
		utils.RespondError(c, http.StatusNotFound, "patient not found")
		return
	}
	if patient.UserID != userID {
		utils.RespondError(c, http.StatusForbidden, "patient does not belong to user")
		return
	}

	var req ListRecommendationsParams
	if err := c.ShouldBindQuery(&req); err != nil {
		utils.RespondError(c, http.StatusBadRequest, "invalid query parameters")
		return
	}

	templates, err := h.Queries.ListExerciseTemplates(c)
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "cannot load templates")
		return
	}
	tplByID := indexTemplates(templates)

	items, err := h.Queries.ListExerciseRecommendationsByPatient(c, db.ListExerciseRecommendationsByPatientParams{
		PatientID: int64(patientID),
		Limit:     req.Limit,
		Offset:    req.Offset,
	})
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "cannot list recommendations")
		return
	}

	resp := ListRecommendationsResponse{Recommendations: make([]RecommendationResponse, 0, len(items))}
	for _, r := range items {
		rec := toRecommendationDomain(r, tplByID)
		resp.Recommendations = append(resp.Recommendations, toRecommendationResponse(rec))
	}

	c.JSON(http.StatusOK, resp)
}
