import { useState } from 'react';
import { useQuery } from '@tanstack/react-query';
import { Link, useParams } from 'react-router-dom';
import { downloadPatientReport } from '../api';
import { Button } from '../../../components/ui/Button';
import { Card } from '../../../components/ui/Card';
import { RecommendationsList } from '../../exercises/components/RecommendationsList';
import { PredictionResultCard } from '../../predictions/components/PredictionResultCard';
import { listPredictions } from '../../predictions/api';
import { PatientSummaryCard } from '../components/PatientSummaryCard';
import { listRecommendations } from '../../exercises/api';
import { getPatient } from '../api';
import { PatientResponse } from '../types';
import { PredictionResponse } from '../../predictions/types';
import { RecommendationResponse } from '../../exercises/types';

function PatientDetailPage() {
  const { id } = useParams<{ id: string }>();
  const { data: patient, isLoading } = useQuery<PatientResponse>({
    queryKey: ['patient', id],
    queryFn: () => {
      if (!id) throw new Error('Missing patient id');
      return getPatient(id);
    },
    enabled: !!id,
  });
  const { data: latestPredictions } = useQuery<PredictionResponse[]>({
    queryKey: ['predictions', id, 1],
    queryFn: () => {
      if (!id) throw new Error('Missing patient id');
      return listPredictions(id, { limit: 1 });
    },
    enabled: !!id,
  });
  const lastPred = latestPredictions?.[0];
  const { data: recommendations = [], isLoading: recLoading } = useQuery<RecommendationResponse[]>({
    queryKey: ['exercise-recommendations', id, 5],
    queryFn: () => {
      if (!id) throw new Error('Missing patient id');
      return listRecommendations(id, { limit: 5 });
    },
    enabled: !!id,
  });
  const [downloading, setDownloading] = useState(false);

  const handleDownload = async () => {
    if (!id) return;
    try {
      setDownloading(true);
      const blob = await downloadPatientReport(id);
      const url = window.URL.createObjectURL(blob);
      const link = document.createElement('a');
      link.href = url;
      link.download = `patient_${id}_report.pdf`;
      link.click();
      window.URL.revokeObjectURL(url);
    } finally {
      setDownloading(false);
    }
  };

  if (isLoading) return <p className="text-sm text-slate-600">Đang tải...</p>;
  if (!patient) return <p className="text-sm text-red-600">Không tìm thấy bệnh nhân.</p>;

  return (
    <div className="space-y-4">
      <div className="flex items-center justify-between">
        <h2 className="text-2xl font-semibold text-slate-900">{patient.name}</h2>
        <div className="flex gap-2">
          <Link to={`/patients/${patient.id}/predict`}>
            <Button>Predict now</Button>
          </Link>
          <Link to={`/patients/${patient.id}/edit`}>
            <Button variant="secondary">Chỉnh sửa</Button>
          </Link>
          <Link to={`/patients/${patient.id}/history`}>
            <Button variant="secondary">Lịch sử</Button>
          </Link>
          <Button variant="secondary" onClick={handleDownload} disabled={downloading}>
            {downloading ? 'Đang xuất...' : 'Xuất PDF'}
          </Button>
        </div>
      </div>

      <PatientSummaryCard patient={patient} />

      <Card title="Dự đoán gần nhất">
        {lastPred ? (
          <PredictionResultCard prediction={lastPred} />
        ) : (
          <p className="text-sm text-slate-600">Chưa có dự đoán nào.</p>
        )}
      </Card>

      <Card title="Kế hoạch tập luyện gần nhất">
        {recLoading ? (
          <p className="text-sm text-slate-600">Đang tải...</p>
        ) : (
          <RecommendationsList recommendations={recommendations} />
        )}
      </Card>
    </div>
  );
}

export default PatientDetailPage;
