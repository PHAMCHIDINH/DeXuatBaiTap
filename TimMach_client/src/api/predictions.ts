import client from './client';
import {
  CreatePredictionRequest,
  ListPredictionsParams,
  ListPredictionsResponse,
  CreatePredictionResponse,
  PredictionResponse,
} from '../types/api';

export async function createPrediction(
  patientId: string,
  payload: CreatePredictionRequest,
): Promise<CreatePredictionResponse> {
  const { data } = await client.post<CreatePredictionResponse>(`/patients/${patientId}/predict`, payload);
  return data;
}

export async function listPredictions(
  patientId: string,
  params: ListPredictionsParams,
): Promise<ListPredictionsResponse> {
  const { data } = await client.get<ListPredictionsResponse>(`/patients/${patientId}/predictions`, { params });
  console.log("data", data);
  return data;
}
