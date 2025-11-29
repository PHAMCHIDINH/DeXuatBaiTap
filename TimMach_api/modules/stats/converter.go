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

func buildStatsResponse(
	totalPatients int64,
	totalPredictions int64,
	distRows []db.GetRiskDistributionRow,
) StatsResponse {
	counts := map[string]int64{
		"low":    0,
		"medium": 0,
		"high":   0,
	}
	for _, r := range distRows {
		label := strings.ToLower(r.RiskLabel)
		if _, ok := counts[label]; ok {
			counts[label] = r.Count
		}
	}

	return StatsResponse{
		TotalPatients:    totalPatients,
		TotalPredictions: totalPredictions,
		RiskCounts: []RiskCount{
			{RiskLabel: "high", Count: counts["high"]},
			{RiskLabel: "medium", Count: counts["medium"]},
			{RiskLabel: "low", Count: counts["low"]},
		},
	}
}
