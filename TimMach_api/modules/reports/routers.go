package reports

import "github.com/gin-gonic/gin"

// RegisterReportRoutes wires report endpoints (including legacy patient routes).
func RegisterReportRoutes(r *gin.RouterGroup, h *Controller, auth gin.HandlerFunc) {
	patientGroup := r.Group("/patients")
	if auth != nil {
		patientGroup.Use(auth)
	}
	patientGroup.GET("/:id/report.pdf", h.GetPatientReport)
	patientGroup.POST("/:id/report/email", h.SendPatientReportEmail)
	patientGroup.POST("/:id/reports", h.CreateReport)
	patientGroup.GET("/:id/reports", h.ListReports)

	reportGroup := r.Group("/reports")
	if auth != nil {
		reportGroup.Use(auth)
	}
	reportGroup.GET("/:id/download", h.DownloadReport)
	reportGroup.POST("/:id/email", h.SendReportEmail)
	reportGroup.DELETE("/:id", h.DeleteReport)
}
