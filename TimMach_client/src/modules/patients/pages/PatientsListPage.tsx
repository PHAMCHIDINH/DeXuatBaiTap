import { Link, useSearchParams } from 'react-router-dom';
import { Card } from '../../../components/ui/Card';
import { Button } from '../../../components/ui/Button';
import { PatientsTable } from '../components/PatientsTable';
import { usePatientsList } from '../hooks/usePatients';

function parseNumber(value: string | null, fallback: number) {
  const n = value ? Number(value) : NaN;
  return Number.isFinite(n) && n >= 0 ? n : fallback;
}

const riskOptions = ['high', 'medium', 'low', 'none'] as const;

function parseRisk(value: string | null): string | undefined {
  if (!value) return undefined;
  const normalized = value.toLowerCase();
  return riskOptions.includes(normalized as (typeof riskOptions)[number]) ? normalized : undefined;
}

function PatientsListPage() {
  const [params, setParams] = useSearchParams();
  const limit = parseNumber(params.get('limit'), 10);
  const offset = parseNumber(params.get('offset'), 0);
  const risk = parseRisk(params.get('risk'));

  const { data, isLoading } = usePatientsList({ limit, offset, risk });
  const patients = data?.patients ?? [];

  const nextPage = () => setParams({ limit: String(limit), offset: String(offset + limit), risk: risk ?? '' });
  const prevPage = () =>
    setParams({ limit: String(limit), offset: String(Math.max(0, offset - limit)), risk: risk ?? '' });

  const onChangeRisk = (value: string) => {
    setParams({
      limit: String(limit),
      offset: '0',
      risk: value,
    });
  };

  return (
    <div className="space-y-4">
      <div className="flex items-center justify-between gap-3">
        <div>
          <p className="text-xs uppercase tracking-wide text-slate-500">Danh sách</p>
          <h2 className="text-2xl font-semibold text-slate-900">Patients</h2>
        </div>
        <Link to="/patients/new">
          <Button>Tạo bệnh nhân</Button>
        </Link>
      </div>

      <Card>
        <div className="mb-3 flex flex-wrap items-center gap-3">
          <label className="text-sm text-slate-600">
            Lọc theo nguy cơ:
            <select
              className="ml-2 rounded-lg border border-slate-200 bg-white px-3 py-2 text-sm shadow-sm focus:border-blue-500 focus:outline-none"
              value={risk ?? ''}
              onChange={(e) => onChangeRisk(e.target.value)}
            >
              <option value="">Tất cả</option>
              {riskOptions.map((r) => (
                <option key={r} value={r}>
                  {r}
                </option>
              ))}
            </select>
          </label>
        </div>
        {isLoading ? (
          <p className="text-sm text-slate-600">Đang tải...</p>
        ) : patients.length === 0 ? (
          <p className="text-sm text-slate-600">Chưa có bệnh nhân nào.</p>
        ) : (
          <PatientsTable patients={patients} />
        )}
        <div className="mt-3 flex items-center justify-between text-sm text-slate-700">
          <div>
            Trang: {Math.floor(offset / limit) + 1}
          </div>
          <div className="flex gap-2">
            <Button variant="secondary" size="sm" onClick={prevPage} disabled={offset === 0}>
              Trước
            </Button>
            <Button variant="secondary" size="sm" onClick={nextPage} disabled={patients.length < limit}>
              Sau
            </Button>
          </div>
        </div>
      </Card>
    </div>
  );
}

export default PatientsListPage;
