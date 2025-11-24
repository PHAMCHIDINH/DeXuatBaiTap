package reports

import (
	"encoding/json"
	"fmt"

	db "chidinh/db/sqlc"
)

func mapReportResponse(r db.Report) (ReportResponse, error) {
	recipients, err := decodeRecipients(r.Recipients)
	if err != nil {
		return ReportResponse{}, err
	}
	return ReportResponse{
		ID:         r.ID,
		PatientID:  r.PatientID,
		Filename:   r.Filename,
		FileURL:    fmt.Sprintf("/reports/%d/download", r.ID),
		Recipients: recipients,
		CreatedAt:  r.CreatedAt.Time,
	}, nil
}

func decodeRecipients(raw []byte) ([]ReportRecipient, error) {
	if len(raw) == 0 {
		return []ReportRecipient{}, nil
	}
	var recipients []ReportRecipient
	if err := json.Unmarshal(raw, &recipients); err != nil {
		return nil, err
	}
	return recipients, nil
}
