import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import * as api from '../api';
import {
  ListTemplatesResponse,
  CreateTemplateRequest,
  TemplateResponse,
  ListRecommendationsResponse,
} from '../../../types/api';

export function useTemplates() {
  return useQuery<ListTemplatesResponse>({
    queryKey: ['exercise-templates'],
    queryFn: api.listTemplates,
  });
}

export function useCreateTemplate() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (payload: CreateTemplateRequest): Promise<TemplateResponse> => api.createTemplate(payload),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ['exercise-templates'] });
    },
  });
}

export function useRecommendations(patientId?: string, params?: { limit?: number; offset?: number }) {
  return useQuery<ListRecommendationsResponse>({
    queryKey: ['exercise-recommendations', patientId, params?.limit ?? null, params?.offset ?? null],
    queryFn: () => {
      if (!patientId) throw new Error('Missing patient id');
      return api.listRecommendations(patientId, params);
    },
    enabled: !!patientId,
  });
}
