package patients

import "github.com/gin-gonic/gin"

// RegisterPatientRoutes gắn CRUD patient. Nên đặt group này sau middleware auth.
func RegisterPatientRoutes(r *gin.RouterGroup, h *Controller, auth gin.HandlerFunc) {
	group := r.Group("/patients")
	if auth != nil {
		group.Use(auth)
	}

	group.POST("", h.CreatePatient)
	group.GET("", h.ListPatients)
	group.GET("/:id", h.GetPatient)
	group.PUT("/:id", h.UpdatePatient)
	group.PATCH("/:id", h.UpdatePatient)
	group.DELETE("/:id", h.DeletePatient)
}
