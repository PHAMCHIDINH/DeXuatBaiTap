import { RecommendationResponse } from '../../../types/api';
import { Card } from '../../../components/ui/Card';
import { formatDate } from '../../../utils/format';

interface Props {
  recommendations: RecommendationResponse[];
}

export function RecommendationsList({ recommendations }: Props) {
  if (!recommendations || recommendations.length === 0) {
    return <p className="text-sm text-slate-600">Chưa có kế hoạch tập luyện.</p>;
  }

  return (
    <div className="space-y-3">
      {recommendations.map((rec) => {
        const plan = rec.plan;
        const items = plan?.items ?? [];

        return (
          <Card key={rec.id} title={`Kế hoạch #${rec.id}`}>
            <div className="mb-2 flex items-center justify-between text-sm text-slate-600">
              <span>Prediction ID: {rec.prediction_id}</span>
              <span>{formatDate(rec.created_at)}</span>
            </div>
            {plan ? (
              <div className="space-y-2">
                <p className="text-slate-800">{plan.summary || 'Kế hoạch tập luyện'}</p>
                {items.length > 0 ? (
                  <ul className="list-disc pl-5 text-sm text-slate-700">
                    {items.map((item, idx) => (
                      <li key={idx}>
                        <span className="font-medium">{item.name}</span> — {item.intensity},{' '}
                        {item.duration_min} phút, {item.freq_per_week} buổi/tuần. {item.notes ?? ''}
                      </li>
                    ))}
                  </ul>
                ) : (
                  <p className="text-sm text-slate-600">Không có chi tiết bài tập.</p>
                )}
              </div>
            ) : (
              <p className="text-sm text-slate-600">Không có dữ liệu kế hoạch.</p>
            )}
          </Card>
        );
      })}
    </div>
  );
}
