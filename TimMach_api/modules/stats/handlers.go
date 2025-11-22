package stats

import (
	"net/http"
	"strings"

	db "chidinh/db/sqlc"
	"chidinh/utils"

	"github.com/gin-gonic/gin"
)

// Handler gom dependencies cho module thống kê.
type Handler struct {
	Queries *db.Queries
}

func NewHandler(q *db.Queries) *Handler {
	return &Handler{Queries: q}
}

// GET /stats
func (h *Handler) GetStats(c *gin.Context) {
	userID, ok := utils.UserIDFromContext(c)
	if !ok {
		utils.RespondError(c, http.StatusUnauthorized, "missing user in context")
		return
	}

	total, err := h.Queries.CountPatientsByUser(c, userID)
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "cannot count patients")
		return
	}

	rows, err := h.Queries.CountLatestRiskByUser(c, userID)
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "cannot count risk groups")
		return
	}

	counts := map[string]int64{
		"low":    0,
		"medium": 0,
		"high":   0,
		"none":   0,
	}
	for _, r := range rows {
		label := strings.ToLower(r.RiskLabel)
		counts[label] += r.Count
	}

	resp := StatsResponse{
		TotalPatients: total,
		RiskCounts: []RiskCount{
			{RiskLabel: "high", Count: counts["high"]},
			{RiskLabel: "medium", Count: counts["medium"]},
			{RiskLabel: "low", Count: counts["low"]},
			{RiskLabel: "none", Count: counts["none"]},
		},
	}

	c.JSON(http.StatusOK, resp)
}
