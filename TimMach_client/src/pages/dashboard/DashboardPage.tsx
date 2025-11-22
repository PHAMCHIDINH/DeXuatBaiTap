import { Link } from 'react-router-dom';
import { usePatientsList } from '../../hooks/usePatients';
import { useStats } from '../../hooks/useStats';
import { Card } from '../../components/ui/Card';
import { Button } from '../../components/ui/Button';
import { formatDate } from '../../utils/format';

function DashboardPage() {
  const { data, isLoading } = usePatientsList({ limit: 5 });
  const { data: stats, isLoading: statsLoading } = useStats();
  const recentPatients = data?.patients ?? [];
  const riskCounts = stats?.risk_counts ?? [];

  return (
    <div className="space-y-4">
      <div className="grid gap-3 md:grid-cols-3">
        <Card title="Bệnh nhân" className="bg-white">
          <p className="text-3xl font-semibold text-blue-700">
            {statsLoading ? '—' : stats?.total_patients ?? 0}
          </p>
          <p className="text-sm text-slate-600">Tổng số bệnh nhân</p>
        </Card>
        <Card title="Dự đoán" className="bg-white">
          <p className="text-3xl font-semibold text-slate-800">—</p>
          <p className="text-sm text-slate-600">Thêm API /stats để hiển thị</p>
        </Card>
        <Card title="Hành động nhanh" className="bg-white">
          <div className="flex flex-col gap-2">
            <Button variant="primary" onClick={() => {}} disabled>
              Dự đoán nhanh (chọn bệnh nhân)
            </Button>
            <Button variant="secondary" onClick={() => {}} disabled>
              Tạo báo cáo
            </Button>
          </div>
        </Card>
      </div>

      <Card title="Nhóm nguy cơ" className="bg-white">
        {statsLoading ? (
          <p className="text-sm text-slate-600">Đang tải...</p>
        ) : (
          <div className="grid gap-3 sm:grid-cols-2 lg:grid-cols-4">
            {['high', 'medium', 'low', 'none'].map((label) => {
              const found = riskCounts.find((r) => r.risk_label === label);
              const value = found?.count ?? 0;
              const colors: Record<string, string> = {
                high: 'text-red-600',
                medium: 'text-amber-600',
                low: 'text-green-600',
                none: 'text-slate-500',
              };
              const names: Record<string, string> = {
                high: 'High',
                medium: 'Medium',
                low: 'Low',
                none: 'Chưa có dữ liệu',
              };
              return (
                <div key={label} className="rounded-lg border border-slate-100 p-3 shadow-sm">
                  <p className="text-xs uppercase tracking-wide text-slate-500">{names[label]}</p>
                  <p className={`text-2xl font-semibold ${colors[label]}`}>{value}</p>
                  {label !== 'none' && (
                    <Link to={`/patients?risk=${label}`}>
                      <Button variant="secondary" size="sm" className="mt-2 w-full">
                        Xem danh sách
                      </Button>
                    </Link>
                  )}
                </div>
              );
            })}
          </div>
        )}
      </Card>

      <Card
        title="Bệnh nhân mới"
        action={
          <Link to="/patients" className="text-sm font-medium text-blue-600 hover:underline">
            Xem tất cả
          </Link>
        }
      >
        {isLoading && <p className="text-sm text-slate-600">Đang tải...</p>}
        {!isLoading && recentPatients.length === 0 && (
          <p className="text-sm text-slate-600">Chưa có bệnh nhân nào.</p>
        )}
        <div className="divide-y divide-slate-100">
          {recentPatients.map((p) => (
            <div key={p.id} className="flex items-center justify-between py-3">
              <div>
                <p className="font-medium text-slate-900">{p.name}</p>
                <p className="text-xs text-slate-500">DOB: {p.dob} • Tạo: {formatDate(p.created_at)}</p>
              </div>
              <Link to={`/patients/${p.id}`} className="text-sm text-blue-600 hover:underline">
                Chi tiết
              </Link>
            </div>
          ))}
        </div>
      </Card>
    </div>
  );
}

export default DashboardPage;
