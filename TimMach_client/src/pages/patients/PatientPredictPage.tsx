import { useMemo, useState } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import { PredictForm } from '../../components/predictions/PredictForm';
import { PredictionResultCard } from '../../components/predictions/PredictionResultCard';
import { RecommendationCard } from '../../components/predictions/RecommendationCard';
import { useCreatePrediction } from '../../hooks/usePredictions';
import { usePatient } from '../../hooks/usePatients';
import { Card } from '../../components/ui/Card';
import { Button } from '../../components/ui/Button';
import { CreatePredictionRequest, CreatePredictionResponse } from '../../types/api';

function ageFromDob(dob?: string) {
  if (!dob) return undefined;
  const date = new Date(dob);
  if (Number.isNaN(date.getTime())) return undefined;
  const diff = Date.now() - date.getTime();
  return Math.floor(diff / (1000 * 60 * 60 * 24 * 365.25));
}

function PatientPredictPage() {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const { data: patient, isLoading } = usePatient(id);
  const patientId = id ?? '';
  const mutation = useCreatePrediction(patientId);
  const [result, setResult] = useState<CreatePredictionResponse | null>(null);

  const defaultAge = useMemo(() => ageFromDob(patient?.dob), [patient?.dob]);

  const handleSubmit = async (values: CreatePredictionRequest) => {
    if (!patientId) return;
    const res = await mutation.mutateAsync(values);
    setResult(res);
  };

  if (!id) return <p className="text-sm text-red-600">Thiếu patient id.</p>;
  if (isLoading) return <p className="text-sm text-slate-600">Đang tải...</p>;
  if (!patient) return <p className="text-sm text-red-600">Không tìm thấy bệnh nhân.</p>;

  return (
    <div className="space-y-4">
      <div className="flex items-center justify-between">
        <div>
          <p className="text-xs uppercase tracking-wide text-slate-500">Dự đoán</p>
          <h2 className="text-2xl font-semibold text-slate-900">{patient.name}</h2>
        </div>
        <Button variant="secondary" onClick={() => navigate(-1)}>
          Quay lại
        </Button>
      </div>

      <Card title="Thông tin dự đoán">
        <PredictForm
          defaultValues={{ age_years: defaultAge }}
          onSubmit={handleSubmit}
        />
      </Card>

      {result && (
        <div className="grid gap-3 md:grid-cols-2">
          <PredictionResultCard prediction={result.prediction} />
          <RecommendationCard recommendation={result.recommendation} />
        </div>
      )}
    </div>
  );
}

export default PatientPredictPage;
