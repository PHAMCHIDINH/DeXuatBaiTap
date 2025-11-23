import { RecommendationPlan } from '../../exercises/types';
import { Card } from '../../../components/ui/Card';

interface Props {
  recommendation: RecommendationPlan;
}

export function RecommendationCard({ recommendation }: Props) {
  return (
    <Card title="Gợi ý tập luyện">
      <div className="space-y-3 text-sm">
        <p className="text-slate-700">{recommendation.summary}</p>
        <div className="divide-y divide-slate-100 rounded-lg border border-slate-100 bg-slate-50">
          {recommendation.items.map((item, idx) => (
            <div key={`${item.name}-${idx}`} className="p-3">
              <div className="flex items-center justify-between">
                <p className="font-semibold text-slate-900">{item.name}</p>
                <span className="text-xs font-medium uppercase text-slate-500">{item.intensity}</span>
              </div>
              <p className="text-xs text-slate-600">
                {item.duration_min} phút/buổi • {item.freq_per_week} buổi/tuần
              </p>
              {item.notes && <p className="mt-1 text-xs text-slate-600">Ghi chú: {item.notes}</p>}
            </div>
          ))}
        </div>
      </div>
    </Card>
  );
}
