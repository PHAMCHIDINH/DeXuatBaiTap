import client from '../../api/client';
import {
  CreatePredictionRequest,
  ListPredictionsParams,
  CreatePredictionResponse,
  PredictionResponse,
} from './types';

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
): Promise<PredictionResponse[]> {
  const { data } = await client.get<{ predictions: PredictionResponse[] }>(
    `/patients/${patientId}/predictions`,
    { params },
  );
  return data.predictions;
}
