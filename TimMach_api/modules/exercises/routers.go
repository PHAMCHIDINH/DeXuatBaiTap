package exercises

import "github.com/gin-gonic/gin"

// RegisterExerciseRoutes đăng ký endpoint template và recommendation (cần auth).
func RegisterExerciseRoutes(r *gin.RouterGroup, h *Controller) {
	templates := r.Group("/exercise-templates")
	templates.GET("", h.ListTemplates)
	templates.POST("", h.CreateTemplate)

	recs := r.Group("/patients")
	recs.GET("/:id/recommendations", h.ListRecommendationsByPatient)
}
