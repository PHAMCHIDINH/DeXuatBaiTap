import { PredictionResponse } from '../../../types/api';
import { formatDate, formatPercent } from '../../../utils/format';
import { Badge } from '../../../components/ui/Badge';

interface Props {
  predictions: PredictionResponse[];
}

function riskColor(label: string): 'green' | 'amber' | 'red' | 'gray' {
  const value = label.toLowerCase();
  if (value.includes('low')) return 'green';
  if (value.includes('medium')) return 'amber';
  if (value.includes('high')) return 'red';
  return 'gray';
}

export function PredictionHistoryTable({ predictions }: Props) {
  return (
    <div className="overflow-hidden rounded-xl border border-slate-200 bg-white">
      <table className="min-w-full divide-y divide-slate-100 text-sm">
        <thead className="bg-slate-50 text-slate-600">
          <tr>
            <th className="px-4 py-3 text-left font-semibold">Thời gian</th>
            <th className="px-4 py-3 text-left font-semibold">Xác suất</th>
            <th className="px-4 py-3 text-left font-semibold">Risk</th>
            <th className="px-4 py-3 text-left font-semibold">Model</th>
          </tr>
        </thead>
        <tbody className="divide-y divide-slate-100">
          {predictions.map((p) => (
            <tr key={p.id} className="hover:bg-slate-50">
              <td className="px-4 py-3 text-slate-800">{formatDate(p.created_at)}</td>
              <td className="px-4 py-3 text-slate-800">{formatPercent(p.probability)}</td>
              <td className="px-4 py-3 text-slate-800">
                <Badge color={riskColor(p.risk_label)}>{p.risk_label}</Badge>
              </td>
              <td className="px-4 py-3 text-slate-800">{p.model_version || 'v1'}</td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}
