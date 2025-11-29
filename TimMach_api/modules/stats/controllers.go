package stats

import (
	"net/http"

	db "chidinh/db/sqlc"
	"chidinh/utils"

	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
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

	var totalPatients int64
	var totalPredictions int64
	var riskDist []db.GetRiskDistributionRow

	g, ctx := errgroup.WithContext(c)

	g.Go(func() error {
		var err error
		totalPatients, err = h.Queries.GetTotalPatients(ctx, userID)
		return err
	})

	g.Go(func() error {
		var err error
		totalPredictions, err = h.Queries.GetTotalPredictions(ctx, userID)
		return err
	})

	g.Go(func() error {
		var err error
		riskDist, err = h.Queries.GetRiskDistribution(ctx, userID)
		return err
	})

	if err := g.Wait(); err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "cannot query stats")
		return
	}

	c.JSON(http.StatusOK, buildStatsResponse(totalPatients, totalPredictions, riskDist))
}
