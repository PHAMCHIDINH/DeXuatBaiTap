package predictions

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	db "chidinh/db/sqlc"
	"chidinh/utils"
	"chidinh/utils/httpclient"

	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"github.com/jackc/pgx/v5"
)

// Controller gom dependency cho module predictions.
type Controller struct {
	Queries *db.Queries
	MLHTTP  *resty.Client
}

// NewController khởi tạo controller với DB queries và HTTP client gọi ML.
func NewController(queries *db.Queries, mlHTTP *resty.Client) *Controller {
	return &Controller{Queries: queries, MLHTTP: mlHTTP}
}

// POST /patients/:id/predict
// Nhận input từ client -> gửi ML -> lưu prediction + factors -> trả kết quả và recommendation.
func (h *Controller) CreatePrediction(c *gin.Context) {
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

	var req CreatePredictionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.L().Warnf("create prediction: bind error: %v", err)
		utils.RespondError(c, http.StatusBadRequest, "invalid request body")
		return
	}

	mlPayload := toMLRequest(req)

	if h.MLHTTP == nil {
		utils.RespondError(c, http.StatusInternalServerError, "ml client is not configured")
		return
	}

	var mlResp MLResponse
	if _, err := httpclient.PostJSON(c, h.MLHTTP, "/predict", mlPayload, &mlResp); err != nil {
		utils.RespondError(c, http.StatusBadGateway, fmt.Sprintf("ml service error: %v", err))
		return
	}

	rawFeatures := encodeStoredFeatures(mlPayload)
	factorsJSON, _ := json.Marshal(mlResp.Factors)

	pred, err := h.Queries.CreatePrediction(c, db.CreatePredictionParams{
		PatientID:   int64(patientID),
		Probability: mlResp.Probability,
		RiskLabel:   mlResp.RiskLevel,
		RawFeatures: rawFeatures,
		Factors:     factorsJSON,
	})
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "cannot save prediction")
		return
	}

	plan, planBlob, err := buildRecommendation(c, h.Queries, mlResp.RiskLevel)
	if err != nil {
		utils.L().Warnf("create recommendation error: %v", err)
		planBlob, _ = json.Marshal(plan)
	}

	if planBlob != nil {
		_, recErr := h.Queries.CreateExerciseRecommendation(c, db.CreateExerciseRecommendationParams{
			PatientID:    int64(patientID),
			PredictionID: pred.ID,
			Plan:         planBlob,
		})
		if recErr != nil {
			utils.L().Warnf("save recommendation failed: %v", recErr)
		}
	}

	c.JSON(http.StatusCreated, CreatePredictionResponse{
		Prediction:     toPredictionResponse(toPredictionDomain(pred)),
		Recommendation: plan,
	})
}

// GET /patients/:id/predictions
// Lấy danh sách dự đoán của bệnh nhân (theo thời gian giảm dần).
func (h *Controller) ListPredictions(c *gin.Context) {
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

	var req ListPredictionsParams
	if err := c.ShouldBindQuery(&req); err != nil {
		utils.RespondError(c, http.StatusBadRequest, "invalid query parameters")
		return
	}

	items, err := h.Queries.ListPredictionsByPatient(c, db.ListPredictionsByPatientParams{
		PatientID: int64(patientID),
		Limit:     req.Limit,
		Offset:    req.Offset,
	})
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "cannot list predictions")
		return
	}

	resp := ListPredictionsResponse{Predictions: make([]PredictionResponse, 0, len(items))}
	for _, p := range items {
		resp.Predictions = append(resp.Predictions, toPredictionResponse(toPredictionDomain(p)))
	}

	c.JSON(http.StatusOK, resp)
}

// buildRecommendation selects up to 3 templates matching risk level (no hardcoded defaults).
// Returns both response plan and a compact blob (summary + template IDs) to persist.
func buildRecommendation(ctx context.Context, q *db.Queries, risk string) (RecommendationPlan, []byte, error) {
	respPlan := RecommendationPlan{
		Summary: "",
		Items:   []RecommendationItem{},
	}

	templates, err := q.ListExerciseTemplates(ctx)
	if err != nil {
		return respPlan, nil, err
	}

	riskLower := strings.ToLower(risk)
	filtered := make([]db.ExerciseTemplate, 0, len(templates))
	for _, t := range templates {
		target := strings.ToLower(t.TargetRiskLevel)
		if target == "" || target == riskLower {
			filtered = append(filtered, t)
		}
	}
	if len(filtered) == 0 {
		filtered = templates
	}

	maxItems := 3
	templateIDs := make([]int64, 0, maxItems)
	for i, t := range filtered {
		if i >= maxItems {
			break
		}
		templateIDs = append(templateIDs, t.ID)
		respPlan.Items = append(respPlan.Items, RecommendationItem{
			Name:        t.Name,
			Intensity:   t.Intensity,
			DurationMin: int(t.DurationMin),
			FreqPerWeek: int(t.FreqPerWeek),
			Notes:       t.Description,
		})
	}

	storePlan := RecommendationPlan{
		Summary:     respPlan.Summary,
		TemplateIDs: templateIDs,
	}

	blob, err := json.Marshal(storePlan)
	return respPlan, blob, err
}
