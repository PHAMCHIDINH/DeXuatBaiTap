package exercises

import "github.com/gin-gonic/gin"

// RegisterExerciseRoutes đăng ký endpoint template và recommendation (cần auth).
func RegisterExerciseRoutes(r *gin.RouterGroup, h *Handler, auth gin.HandlerFunc) {
	templates := r.Group("/exercise-templates")
	if auth != nil {
		templates.Use(auth)
	}
	templates.GET("", h.ListTemplates)
	templates.POST("", h.CreateTemplate)

	recs := r.Group("/patients")
	if auth != nil {
		recs.Use(auth)
	}
	recs.GET("/:id/recommendations", h.ListRecommendationsByPatient)
}
