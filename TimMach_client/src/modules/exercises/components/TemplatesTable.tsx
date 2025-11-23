import { TemplateResponse } from '../../../types/api';
import { Card } from '../../../components/ui/Card';

interface Props {
  templates: TemplateResponse[];
}

export function TemplatesTable({ templates }: Props) {
  if (templates.length === 0) {
    return <p className="text-sm text-slate-600">Chưa có template nào.</p>;
  }

  return (
    <Card>
      <div className="overflow-x-auto">
        <table className="min-w-full divide-y divide-slate-100 text-sm">
          <thead className="bg-slate-50 text-slate-600">
            <tr>
              <th className="px-4 py-3 text-left font-semibold">Tên</th>
              <th className="px-4 py-3 text-left font-semibold">Cường độ</th>
              <th className="px-4 py-3 text-left font-semibold">Thời lượng</th>
              <th className="px-4 py-3 text-left font-semibold">Buổi/tuần</th>
              <th className="px-4 py-3 text-left font-semibold">Risk</th>
              <th className="px-4 py-3 text-left font-semibold">Mô tả</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-slate-100">
            {templates.map((t) => (
              <tr key={t.id} className="hover:bg-slate-50">
                <td className="px-4 py-3 font-medium text-slate-900">{t.name}</td>
                <td className="px-4 py-3 text-slate-700">{t.intensity}</td>
                <td className="px-4 py-3 text-slate-700">{t.duration_min} phút</td>
                <td className="px-4 py-3 text-slate-700">{t.freq_per_week} buổi</td>
                <td className="px-4 py-3 text-slate-700">{t.target_risk_level}</td>
                <td className="px-4 py-3 text-slate-700">{t.description}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </Card>
  );
}
