import axios from 'axios';
import { useState } from 'react';
import { useQuery } from '@tanstack/react-query';
import { Link, useParams } from 'react-router-dom';
import { Button } from '../../../components/ui/Button';
import { Card } from '../../../components/ui/Card';
import { Input } from '../../../components/ui/Input';
import { RecommendationsList } from '../../exercises/components/RecommendationsList';
import { PredictionResultCard } from '../../predictions/components/PredictionResultCard';
import { listPredictions } from '../../predictions/api';
import { PatientSummaryCard } from '../components/PatientSummaryCard';
import { listRecommendations } from '../../exercises/api';
import { getPatient } from '../api';
import { PatientResponse } from '../types';
import { PredictionResponse } from '../../predictions/types';
import { RecommendationResponse } from '../../exercises/types';
import { createReport, downloadReport, listReports, sendReportEmail } from '../../reports/api';
import { ListReportsResponse, ReportResponse } from '../../reports/types';

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
  const {
    data: reportData,
    isLoading: reportsLoading,
    refetch: refetchReports,
  } = useQuery<ListReportsResponse>({
    queryKey: ['reports', id],
    queryFn: () => {
      if (!id) throw new Error('Missing patient id');
      return listReports(id, { limit: 5 });
    },
    enabled: !!id,
  });

  const [creatingReport, setCreatingReport] = useState(false);
  const [email, setEmail] = useState('');
  const [sendingEmail, setSendingEmail] = useState(false);
  const [emailStatus, setEmailStatus] = useState<{ type: 'success' | 'error'; text: string } | null>(null);
  const [downloadingReportId, setDownloadingReportId] = useState<number | null>(null);

  const handleCreateReport = async () => {
    if (!id) return;
    setEmailStatus(null);
    try {
      setCreatingReport(true);
      await createReport(id);
      await refetchReports();
    } catch (err) {
      const message = axios.isAxiosError(err)
        ? (err.response?.data as { error?: string })?.error ?? err.message
        : err instanceof Error
          ? err.message
          : 'Không thể tạo báo cáo';
      setEmailStatus({ type: 'error', text: message });
    } finally {
      setCreatingReport(false);
    }
  };

  const handleDownloadReport = async (report: ReportResponse) => {
    setDownloadingReportId(report.id);
    try {
      const blob = await downloadReport(report.id);
      const url = window.URL.createObjectURL(blob);
      const link = document.createElement('a');
      link.href = url;
      link.download = report.filename || `report_${report.id}.pdf`;
      link.click();
      window.URL.revokeObjectURL(url);
    } finally {
      setDownloadingReportId(null);
    }
  };

  const handleSendEmail = async (report: ReportResponse) => {
    if (!email.trim()) {
      setEmailStatus({ type: 'error', text: 'Vui lòng nhập email nhận báo cáo.' });
      return;
    }
    setEmailStatus(null);
    try {
      setSendingEmail(true);
      await sendReportEmail(report.id, { email: email.trim() });
      setEmailStatus({ type: 'success', text: `Đã gửi báo cáo tới ${email.trim()}.` });
      await refetchReports();
    } catch (err) {
      let message = 'Không thể gửi email';
      if (axios.isAxiosError(err)) {
        message = (err.response?.data as { error?: string })?.error ?? err.message;
      } else if (err instanceof Error) {
        message = err.message;
      }
      setEmailStatus({ type: 'error', text: message });
    } finally {
      setSendingEmail(false);
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
        </div>
      </div>

      <PatientSummaryCard patient={patient} />

      <Card title="Báo cáo PDF đã lưu">
        <div className="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
          <div className="flex flex-col gap-2 sm:flex-row sm:items-center">
            <Input
              type="email"
              placeholder="Email nhận báo cáo"
              className="sm:max-w-sm"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
            />
            <Button onClick={handleCreateReport} disabled={creatingReport}>
              {creatingReport ? 'Đang tạo...' : 'Tạo báo cáo mới'}
            </Button>
          </div>
        </div>

        <div className="mt-3 space-y-2">
          {reportsLoading ? (
            <p className="text-sm text-slate-600">Đang tải danh sách báo cáo...</p>
          ) : reportData && reportData.reports.length > 0 ? (
            reportData.reports.map((report) => (
              <div
                key={report.id}
                className="flex flex-col gap-2 rounded-lg border border-slate-200 p-3 sm:flex-row sm:items-center sm:justify-between"
              >
                <div>
                  <p className="text-sm font-medium text-slate-900">{report.filename}</p>
                  <p className="text-xs text-slate-500">
                    Tạo lúc: {new Date(report.created_at).toLocaleString('vi-VN')}
                  </p>
                </div>
                <div className="flex gap-2">
                  <Button
                    variant="secondary"
                    onClick={() => handleDownloadReport(report)}
                    disabled={downloadingReportId === report.id}
                  >
                    {downloadingReportId === report.id ? 'Đang tải...' : 'Tải'}
                  </Button>
                  <Button onClick={() => handleSendEmail(report)} disabled={sendingEmail}>
                    {sendingEmail ? 'Đang gửi...' : 'Gửi email'}
                  </Button>
                </div>
              </div>
            ))
          ) : (
            <p className="text-sm text-slate-600">Chưa có báo cáo nào. Nhấn "Tạo báo cáo mới".</p>
          )}
        </div>

        {emailStatus && (
          <p
            className={`mt-3 text-sm ${emailStatus.type === 'success' ? 'text-green-600' : 'text-red-600'}`}
          >
            {emailStatus.text}
          </p>
        )}
      </Card>

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
