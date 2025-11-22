import { useParams } from 'react-router-dom';
import { usePredictions } from '../../hooks/usePredictions';
import { usePatient } from '../../hooks/usePatients';
import { useRecommendations } from '../../hooks/useExercises';
import { Card } from '../../components/ui/Card';
import { PredictionHistoryTable } from '../../components/predictions/PredictionHistoryTable';
import { PredictionChart } from '../../components/predictions/PredictionChart';
import { RecommendationsList } from '../../components/exercises/RecommendationsList';

function PatientHistoryPage() {
  const { id } = useParams<{ id: string }>();
  const { data: patient } = usePatient(id);
  const { data, isLoading } = usePredictions(id, { limit: 50 });
  const predictions = data?.predictions ?? [];
  const { data: recs, isLoading: recLoading } = useRecommendations(id, { limit: 20 });
  const recommendations = recs?.recommendations ?? [];

  if (!id) return <p className="text-sm text-red-600">Thiếu patient id.</p>;

  return (
    <div className="space-y-4">
      <div>
        <p className="text-xs uppercase tracking-wide text-slate-500">Lịch sử dự đoán</p>
        <h2 className="text-2xl font-semibold text-slate-900">{patient?.name ?? 'Patient'}</h2>
      </div>

      {isLoading ? (
        <p className="text-sm text-slate-600">Đang tải...</p>
      ) : predictions.length === 0 ? (
        <Card>
          <p className="text-sm text-slate-600">Chưa có dữ liệu dự đoán.</p>
        </Card>
      ) : (
        <div className="space-y-3">
          <Card title="Biểu đồ xu hướng">
            <PredictionChart predictions={predictions} />
          </Card>
          <PredictionHistoryTable predictions={predictions} />
          <Card title="Kế hoạch tập luyện">
            {recLoading ? (
              <p className="text-sm text-slate-600">Đang tải...</p>
            ) : (
              <RecommendationsList recommendations={recommendations} />
            )}
          </Card>
        </div>
      )}
    </div>
  );
}

export default PatientHistoryPage;
