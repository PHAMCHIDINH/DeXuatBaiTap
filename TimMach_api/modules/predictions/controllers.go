package predictions

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	db "chidinh/db/sqlc"
	"chidinh/utils"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

// MLClient định nghĩa cách gọi sang service FastAPI.
type MLClient interface {
	Predict(ctx context.Context, payload MLRequest) (MLResponse, error)
}

// Handler gom dependency cho module predictions.
type Handler struct {
	Queries *db.Queries
	ML      MLClient
}

func NewHandler(queries *db.Queries, ml MLClient) *Handler {
	return &Handler{
		Queries: queries,
		ML:      ml,
	}
}

// POST /patients/:id/predict
func (h *Handler) CreatePrediction(c *gin.Context) {
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

	if h.ML == nil {
		utils.RespondError(c, http.StatusInternalServerError, "ml client is not configured")
		return
	}

	mlResp, err := h.ML.Predict(c, mlPayload)
	if err != nil {
		utils.RespondError(c, http.StatusBadGateway, fmt.Sprintf("ml service error: %v", err))
		return
	}

	rawFeatures, _ := json.Marshal(mlPayload)

	pred, err := h.Queries.CreatePrediction(c, db.CreatePredictionParams{
		PatientID:   int64(patientID),
		Probability: mlResp.Probability,
		RiskLabel:   mlResp.RiskLevel,
		RawFeatures: rawFeatures,
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
func (h *Handler) ListPredictions(c *gin.Context) {
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
