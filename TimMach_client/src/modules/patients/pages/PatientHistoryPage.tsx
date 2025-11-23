import { useQuery } from '@tanstack/react-query';
import { useParams } from 'react-router-dom';
import { Card } from '../../../components/ui/Card';
import { RecommendationsList } from '../../exercises/components/RecommendationsList';
import { listRecommendations } from '../../exercises/api';
import { PredictionChart } from '../../predictions/components/PredictionChart';
import { PredictionHistoryTable } from '../../predictions/components/PredictionHistoryTable';
import { listPredictions } from '../../predictions/api';
import { getPatient } from '../api';
import { PatientResponse } from '../types';
import { RecommendationResponse } from '../../exercises/types';
import { PredictionResponse } from '../../predictions/types';

function PatientHistoryPage() {
  const { id } = useParams<{ id: string }>();
  const { data: patient } = useQuery<PatientResponse>({
    queryKey: ['patient', id],
    queryFn: () => {
      if (!id) throw new Error('Missing patient id');
      return getPatient(id);
    },
    enabled: !!id,
  });
  const { data: predictions = [], isLoading } = useQuery<PredictionResponse[]>({
    queryKey: ['predictions', id, 50],
    queryFn: () => {
      if (!id) throw new Error('Missing patient id');
      return listPredictions(id, { limit: 50 });
    },
    enabled: !!id,
  });
  const { data: recommendations = [], isLoading: recLoading } = useQuery<RecommendationResponse[]>({
    queryKey: ['exercise-recommendations', id, 20],
    queryFn: () => {
      if (!id) throw new Error('Missing patient id');
      return listRecommendations(id, { limit: 20 });
    },
    enabled: !!id,
  });

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
