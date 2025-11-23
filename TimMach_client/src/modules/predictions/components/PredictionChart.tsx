import {
  LineChart,
  Line,
  CartesianGrid,
  XAxis,
  YAxis,
  Tooltip,
  ResponsiveContainer,
} from 'recharts';
import { PredictionResponse } from '../../../types/api';
import { formatPercent } from '../../../utils/format';

interface Props {
  predictions: PredictionResponse[];
}

export function PredictionChart({ predictions }: Props) {
  const data = predictions.map((p) => ({
    time: new Date(p.created_at).toLocaleDateString(),
    probability: p.probability * 100,
  }));

  return (
    <div className="h-80 w-full">
      <ResponsiveContainer width="100%" height="100%">
        <LineChart data={data} margin={{ left: 12, right: 12, top: 12, bottom: 12 }}>
          <CartesianGrid strokeDasharray="3 3" stroke="#e2e8f0" />
          <XAxis dataKey="time" tick={{ fontSize: 12 }} />
          <YAxis tickFormatter={(v) => `${v.toFixed(0)}%`} tick={{ fontSize: 12 }} domain={[0, 100]} />
          <Tooltip formatter={(value: number) => formatPercent(value / 100)} />
          <Line type="monotone" dataKey="probability" stroke="#2563eb" strokeWidth={2} dot />
        </LineChart>
      </ResponsiveContainer>
    </div>
  );
}
