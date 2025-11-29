package stats

import "github.com/gin-gonic/gin"

// RegisterStatsRoutes gắn endpoint thống kê (cần auth).
func RegisterStatsRoutes(r *gin.RouterGroup, h *Controller) {
	group := r.Group("/stats")
	group.GET("", h.GetStats)
}
