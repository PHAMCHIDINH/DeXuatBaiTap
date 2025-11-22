import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import * as api from '../api/predictions';
import {
  CreatePredictionRequest,
  ListPredictionsParams,
  ListPredictionsResponse,
  CreatePredictionResponse,
  PredictionResponse,
} from '../types/api';

export function usePredictions(patientId?: string, params: ListPredictionsParams = {}) {
  return useQuery<ListPredictionsResponse>({
    queryKey: ['predictions', patientId, params.limit ?? null, params.offset ?? null],
    queryFn: () => {
      if (!patientId) throw new Error('Missing patient id');
      return api.listPredictions(patientId, params);
    },
    enabled: !!patientId,
  });
}

export function useCreatePrediction(patientId: string) {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (payload: CreatePredictionRequest): Promise<CreatePredictionResponse> =>
      api.createPrediction(patientId, payload),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ['predictions', patientId] });
    },
  });
}

export function useLatestPrediction(patientId?: string) {
  return usePredictions(patientId, { limit: 1 });
}
