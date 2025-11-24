package stats

import (
	"net/http"

	db "chidinh/db/sqlc"
	"chidinh/utils"

	"github.com/gin-gonic/gin"
)

// Controller gom dependencies cho module thống kê.
type Controller struct {
	Queries *db.Queries
}

func NewController(q *db.Queries) *Controller {
	return &Controller{Queries: q}
}

// GET /stats
func (h *Controller) GetStats(c *gin.Context) {
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

	c.JSON(http.StatusOK, buildStatsResponse(total, rows))
}
