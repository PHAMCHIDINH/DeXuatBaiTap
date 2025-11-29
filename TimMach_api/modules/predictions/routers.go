package predictions

import (

	"github.com/gin-gonic/gin"
)

// RegisterPredictionRoutes gắn endpoint predict + history vào group đã có auth.
func RegisterPredictionRoutes(r *gin.RouterGroup, h *Controller) {
	group := r.Group("/patients")

	group.POST("/:id/predict", h.CreatePrediction)
	group.GET("/:id/predictions", h.ListPredictions)
}
