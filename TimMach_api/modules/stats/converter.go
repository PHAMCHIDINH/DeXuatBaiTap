package stats

import (
	"strings"

	db "chidinh/db/sqlc"
)

// RiskSummary represents normalized risk distribution in the domain layer.
type RiskSummary struct {
	Total int64
	Items []RiskCount
}

func buildStatsResponse(total int64, rows []db.CountLatestRiskByUserRow) StatsResponse {
	counts := map[string]int64{
		"low":    0,
		"medium": 0,
		"high":   0,
		"none":   0,
	}
	for _, r := range rows {
		label := strings.ToLower(r.RiskLabel)
		counts[label] += r.Count
	}

	return StatsResponse{
		TotalPatients: total,
		RiskCounts: []RiskCount{
			{RiskLabel: "high", Count: counts["high"]},
			{RiskLabel: "medium", Count: counts["medium"]},
			{RiskLabel: "low", Count: counts["low"]},
			{RiskLabel: "none", Count: counts["none"]},
		},
	}
}
