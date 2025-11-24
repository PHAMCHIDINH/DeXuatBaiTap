import client from '../../api/client';
import {
  CreatePatientRequest,
  ListPatientsParams,
  PatientResponse,
  UpdatePatientRequest,
} from './types';

export async function listPatients(params: ListPatientsParams): Promise<PatientResponse[]> {
  const { data } = await client.get<{ patients: PatientResponse[] }>('/patients', { params });
  return data.patients;
}

export async function getPatient(id: string): Promise<PatientResponse> {
  const { data } = await client.get<PatientResponse>(`/patients/${id}`);
  return data;
}

export async function createPatient(payload: CreatePatientRequest): Promise<PatientResponse> {
  const { data } = await client.post<PatientResponse>('/patients', payload);
  return data;
}

export async function updatePatient(id: string, payload: UpdatePatientRequest): Promise<PatientResponse> {
  const { data } = await client.patch<PatientResponse>(`/patients/${id}`, payload);
  return data;
}

export async function deletePatient(id: string): Promise<void> {
  await client.delete(`/patients/${id}`);
}

export async function downloadPatientReport(id: string): Promise<Blob> {
  const { data } = await client.get(`/patients/${id}/report.pdf`, { responseType: 'blob' });
  return data;
}
