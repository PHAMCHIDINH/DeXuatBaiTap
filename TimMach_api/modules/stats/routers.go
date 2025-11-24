package stats

import "github.com/gin-gonic/gin"

// RegisterStatsRoutes gắn endpoint thống kê (cần auth).
func RegisterStatsRoutes(r *gin.RouterGroup, h *Controller, auth gin.HandlerFunc) {
	group := r.Group("/stats")
	if auth != nil {
		group.Use(auth)
	}
	group.GET("", h.GetStats)
}
