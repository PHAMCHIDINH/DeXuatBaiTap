package reports

import "time"

type SendReportEmailRequest struct {
	Email   string `json:"email" binding:"omitempty,email"`
	Subject string `json:"subject"`
	Message string `json:"message"`
}

type ListReportsRequest struct {
	Limit  int32 `form:"limit" binding:"omitempty,min=1,max=100"`
	Offset int32 `form:"offset" binding:"omitempty,min=0"`
}

type ReportRecipient struct {
	Email  string    `json:"email"`
	SentAt time.Time `json:"sent_at"`
	Status string    `json:"status"`
}

type ReportResponse struct {
	ID         int64             `json:"id"`
	PatientID  int64             `json:"patient_id"`
	Filename   string            `json:"filename"`
	FileURL    string            `json:"file_url"`
	Recipients []ReportRecipient `json:"recipients"`
	CreatedAt  time.Time         `json:"created_at"`
}

type ListReportsResponse struct {
	Reports []ReportResponse `json:"reports"`
	Total   int64            `json:"total"`
}
