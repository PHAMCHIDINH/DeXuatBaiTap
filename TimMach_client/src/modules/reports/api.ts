import client from '../../api/client';
import {
  ListReportsParams,
  ListReportsResponse,
  ReportResponse,
  SendReportEmailRequest,
} from './types';

export async function createReport(patientId: string): Promise<ReportResponse> {
  const { data } = await client.post<ReportResponse>(`/patients/${patientId}/reports`);
  return data;
}

export async function listReports(patientId: string, params?: ListReportsParams): Promise<ListReportsResponse> {
  const { data } = await client.get<ListReportsResponse>(`/patients/${patientId}/reports`, { params });
  return data;
}

export async function downloadReport(reportId: number): Promise<Blob> {
  const { data } = await client.get(`/reports/${reportId}/download`, { responseType: 'blob' });
  return data;
}

export async function sendReportEmail(reportId: number, payload: SendReportEmailRequest): Promise<void> {
  await client.post(`/reports/${reportId}/email`, payload);
}

export async function deleteReport(reportId: number): Promise<void> {
  await client.delete(`/reports/${reportId}`);
}
