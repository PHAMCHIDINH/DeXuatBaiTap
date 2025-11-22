package predictions

import (
	"context"
	"encoding/json"
	"strings"

	db "chidinh/db/sqlc"
)

// buildRecommendation chọn template phù hợp theo risk_level, fallback nếu thiếu dữ liệu.
// Trả về plan cho response và blob (chỉ chứa summary + template_ids) để lưu DB.
func buildRecommendation(ctx context.Context, q *db.Queries, risk string) (RecommendationPlan, []byte, error) {
	respPlan := RecommendationPlan{
		Summary: defaultSummary(risk),
		Items:   []RecommendationItem{},
	}

	templates, err := q.ListExerciseTemplates(ctx)
	if err != nil {
		return respPlan, nil, err
	}

	riskLower := strings.ToLower(risk)
	filtered := make([]db.ExerciseTemplate, 0, len(templates))
	for _, t := range templates {
		target := strings.ToLower(t.TargetRiskLevel)
		if target == "" || target == riskLower {
			filtered = append(filtered, t)
		}
	}
	// nếu không có template nào match, dùng tất cả
	if len(filtered) == 0 {
		filtered = templates
	}

	maxItems := 3
	templateIDs := make([]int64, 0, maxItems)
	for i, t := range filtered {
		if i >= maxItems {
			break
		}
		templateIDs = append(templateIDs, t.ID)
		respPlan.Items = append(respPlan.Items, RecommendationItem{
			Name:        t.Name,
			Intensity:   t.Intensity,
			DurationMin: int(t.DurationMin),
			FreqPerWeek: int(t.FreqPerWeek),
			Notes:       t.Description,
		})
	}

	if len(respPlan.Items) == 0 {
		respPlan.Items = fallbackItems(riskLower)
	}

	storePlan := RecommendationPlan{
		Summary:     respPlan.Summary,
		TemplateIDs: templateIDs,
	}

	blob, err := json.Marshal(storePlan)
	return respPlan, blob, err
}

func defaultSummary(risk string) string {
	switch strings.ToLower(risk) {
	case "high":
		return "Cường độ nhẹ 5-6 buổi/tuần, theo dõi nhịp tim và tham khảo bác sĩ."
	case "medium":
		return "Tập cường độ nhẹ-vừa 4-5 buổi/tuần, kết hợp cardio và kéo giãn."
	default:
		return "Duy trì cường độ nhẹ 4-5 buổi/tuần, chú ý khởi động và giãn cơ."
	}
}

func fallbackItems(risk string) []RecommendationItem {
	base := []RecommendationItem{
		{
			Name:        "Đi bộ nhanh",
			Intensity:   "low",
			DurationMin: 30,
			FreqPerWeek: 5,
			Notes:       "Theo dõi nhịp tim, khởi động 5 phút.",
		},
		{
			Name:        "Yoga/giãn cơ",
			Intensity:   "low",
			DurationMin: 20,
			FreqPerWeek: 3,
			Notes:       "Hít thở đều, tránh sau ăn no.",
		},
	}
	if risk == "high" {
		return base
	}
	return append(base, RecommendationItem{
		Name:        "Đạp xe nhẹ",
		Intensity:   "medium",
		DurationMin: 25,
		FreqPerWeek: 2,
		Notes:       "Giữ nhịp ổn định, không quá sức.",
	})
}
