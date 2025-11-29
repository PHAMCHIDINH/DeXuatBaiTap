package reports

import "github.com/gin-gonic/gin"

// RegisterReportRoutes wires report endpoints (including legacy patient routes).
func RegisterReportRoutes(r *gin.RouterGroup, h *Controller) {
	patientGroup := r.Group("/patients")
	patientGroup.GET("/:id/report.pdf", h.GetPatientReport)
	patientGroup.POST("/:id/report/email", h.SendPatientReportEmail)
	patientGroup.POST("/:id/reports", h.CreateReport)
	patientGroup.GET("/:id/reports", h.ListReports)

	reportGroup := r.Group("/reports")
	reportGroup.GET("/:id/download", h.DownloadReport)
	reportGroup.POST("/:id/email", h.SendReportEmail)
	reportGroup.DELETE("/:id", h.DeleteReport)
}
