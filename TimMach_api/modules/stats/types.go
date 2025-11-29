package stats

type RiskCount struct {
	RiskLabel string `json:"risk_label"`
	Count     int64  `json:"count"`
}

type StatsResponse struct {
	TotalPatients    int64       `json:"total_patients"`
	TotalPredictions int64       `json:"total_predictions"`
	RiskCounts       []RiskCount `json:"risk_counts"`
}
