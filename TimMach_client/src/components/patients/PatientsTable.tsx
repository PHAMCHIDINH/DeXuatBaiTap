import { Link } from 'react-router-dom';
import { PatientResponse } from '../../types/api';
import { formatDate } from '../../utils/format';
import { Badge } from '../ui/Badge';

interface Props {
  patients: PatientResponse[];
}

function genderLabel(value: number) {
  if (value === 1) return 'Male';
  if (value === 2) return 'Female';
  return 'Other';
}

function riskLabel(label?: string) {
  if (!label) return 'None';
  return label.charAt(0).toUpperCase() + label.slice(1);
}

export function PatientsTable({ patients }: Props) {
  return (
    <div className="overflow-hidden rounded-xl border border-slate-200 bg-white">
      <table className="min-w-full divide-y divide-slate-100 text-sm">
        <thead className="bg-slate-50 text-slate-600">
          <tr>
            <th className="px-4 py-3 text-left font-semibold">Tên</th>
            <th className="px-4 py-3 text-left font-semibold">Giới tính</th>
            <th className="px-4 py-3 text-left font-semibold">DOB</th>
            <th className="px-4 py-3 text-left font-semibold">Nguy cơ</th>
            <th className="px-4 py-3 text-left font-semibold">Tạo lúc</th>
            <th className="px-4 py-3 text-left font-semibold">Actions</th>
          </tr>
        </thead>
        <tbody className="divide-y divide-slate-100">
          {patients.map((p) => (
            <tr key={p.id} className="hover:bg-slate-50">
              <td className="px-4 py-3 font-medium text-slate-900">{p.name}</td>
              <td className="px-4 py-3 text-slate-700">
                <Badge color="gray">{genderLabel(p.gender)}</Badge>
              </td>
              <td className="px-4 py-3 text-slate-700">{p.dob}</td>
              <td className="px-4 py-3 text-slate-700">
                {p.latest_prediction ? (
                  <div className="flex flex-col gap-1">
                    <Badge
                      color={
                        p.latest_prediction.risk_label === 'high'
                          ? 'red'
                          : p.latest_prediction.risk_label === 'medium'
                            ? 'amber'
                            : 'green'
                      }
                    >
                      {riskLabel(p.latest_prediction.risk_label)}
                    </Badge>
                    <p className="text-xs text-slate-500">
                      {Math.round(p.latest_prediction.probability * 100)}% •{' '}
                      {formatDate(p.latest_prediction.created_at)}
                    </p>
                  </div>
                ) : (
                  <Badge color="gray">Chưa có</Badge>
                )}
              </td>
              <td className="px-4 py-3 text-slate-700">{formatDate(p.created_at)}</td>
              <td className="px-4 py-3">
                <div className="flex gap-2 text-blue-600">
                  <Link to={`/patients/${p.id}`} className="hover:underline">
                    Chi tiết
                  </Link>
                  <Link to={`/patients/${p.id}/edit`} className="hover:underline">
                    Sửa
                  </Link>
                </div>
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}
