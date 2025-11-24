import { PredictionResponse } from '../types';
import { Card } from '../../../components/ui/Card';
import { Badge } from '../../../components/ui/Badge';
import { formatDate, formatPercent } from '../../../utils/format';

interface Props {
  prediction: PredictionResponse;
}

function riskColor(label: string): 'green' | 'amber' | 'red' | 'gray' {
  const value = label.toLowerCase();
  if (value.includes('low')) return 'green';
  if (value.includes('medium')) return 'amber';
  if (value.includes('high')) return 'red';
  return 'gray';
}

export function PredictionResultCard({ prediction }: Props) {
  return (
    <Card title="Kết quả dự đoán">
      <div className="space-y-3 text-sm">
        <div className="flex items-center gap-3">
          <p className="text-slate-500">Xác suất</p>
          <p className="text-2xl font-semibold text-slate-900">{formatPercent(prediction.probability)}</p>
        </div>
        <div className="flex items-center gap-3">
          <p className="text-slate-500">Risk level</p>
          <Badge color={riskColor(prediction.risk_label)}>{prediction.risk_label}</Badge>
        </div>
        <p className="text-slate-500">Lúc: {formatDate(prediction.created_at)}</p>
        {prediction.factors && prediction.factors.length > 0 && (
          <div>
            <p className="text-slate-500">Yếu tố nguy cơ (từ ML):</p>
            <ul className="mt-1 list-disc space-y-1 pl-5">
              {prediction.factors.map((f, idx) => (
                <li key={`${f.field}-${idx}`} className="text-slate-700">
                  {f.message || f.field}
                </li>
              ))}
            </ul>
          </div>
        )}
      </div>
    </Card>
  );
}
