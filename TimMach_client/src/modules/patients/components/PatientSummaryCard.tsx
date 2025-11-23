import { PatientResponse } from '../types';
import { Card } from '../../../components/ui/Card';
import { formatDate } from '../../../utils/format';

interface Props {
  patient: PatientResponse;
}

export function PatientSummaryCard({ patient }: Props) {
  return (
    <Card title="Thông tin bệnh nhân">
      <div className="grid grid-cols-1 gap-3 text-sm sm:grid-cols-2">
        <div>
          <p className="text-slate-500">Tên</p>
          <p className="font-semibold text-slate-900">{patient.name}</p>
        </div>
        <div>
          <p className="text-slate-500">Giới tính</p>
          <p className="font-semibold text-slate-900">
            {patient.gender === 1 ? 'Male' : patient.gender === 2 ? 'Female' : 'Other'}
          </p>
        </div>
        <div>
          <p className="text-slate-500">Ngày sinh</p>
          <p className="font-semibold text-slate-900">{patient.dob}</p>
        </div>
        <div>
          <p className="text-slate-500">Tạo lúc</p>
          <p className="font-semibold text-slate-900">{formatDate(patient.created_at)}</p>
        </div>
      </div>
    </Card>
  );
}
