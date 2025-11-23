import { useMemo, useState } from 'react';
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { useNavigate, useParams } from 'react-router-dom';
import { Button } from '../../../components/ui/Button';
import { Card } from '../../../components/ui/Card';
import { PredictForm } from '../../predictions/components/PredictForm';
import { PredictionResultCard } from '../../predictions/components/PredictionResultCard';
import { RecommendationCard } from '../../predictions/components/RecommendationCard';
import { createPrediction } from '../../predictions/api';
import { CreatePredictionRequest, CreatePredictionResponse } from '../../predictions/types';
import { getPatient } from '../api';
import { PatientResponse } from '../types';

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
  const patientId = id ?? '';
  const qc = useQueryClient();
  const { data: patient, isLoading } = useQuery<PatientResponse>({
    queryKey: ['patient', patientId],
    queryFn: () => {
      if (!patientId) throw new Error('Missing patient id');
      return getPatient(patientId);
    },
    enabled: !!patientId,
  });
  const mutation = useMutation({
    mutationFn: (payload: CreatePredictionRequest) => createPrediction(patientId, payload),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ['predictions', patientId] });
    },
  });
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
