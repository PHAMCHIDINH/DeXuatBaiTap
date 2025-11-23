import { useState } from 'react';
import { Link, useParams } from 'react-router-dom';
import { downloadPatientReport } from '../api';
import { Button } from '../../../components/ui/Button';
import { Card } from '../../../components/ui/Card';
import { RecommendationsList } from '../../exercises/components/RecommendationsList';
import { useRecommendations } from '../../exercises/hooks/useExercises';
import { PredictionResultCard } from '../../predictions/components/PredictionResultCard';
import { useLatestPrediction } from '../../predictions/hooks/usePredictions';
import { PatientSummaryCard } from '../components/PatientSummaryCard';
import { usePatient } from '../hooks/usePatients';

function PatientDetailPage() {
  const { id } = useParams<{ id: string }>();
  const { data: patient, isLoading } = usePatient(id);
  const { data: latest } = useLatestPrediction(id);
  const lastPred = latest?.predictions?.[0];
  const { data: recs, isLoading: recLoading } = useRecommendations(id, { limit: 5 });
  const recommendations = recs?.recommendations ?? [];
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
