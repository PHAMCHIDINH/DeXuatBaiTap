export interface ReportRecipient {
  email: string;
  sent_at: string;
  status: string;
}

export interface ReportResponse {
  id: number;
  patient_id: number;
  filename: string;
  file_url: string;
  recipients: ReportRecipient[];
  created_at: string;
}

export interface ListReportsResponse {
  reports: ReportResponse[];
  total: number;
}

export interface ListReportsParams {
  limit?: number;
  offset?: number;
}

export interface SendReportEmailRequest {
  email: string;
  subject?: string;
  message?: string;
}
