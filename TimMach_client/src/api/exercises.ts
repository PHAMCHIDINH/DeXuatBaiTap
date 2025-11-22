import client from './client';
import {
  ListTemplatesResponse,
  CreateTemplateRequest,
  TemplateResponse,
  ListRecommendationsResponse,
} from '../types/api';

export async function listTemplates(): Promise<ListTemplatesResponse> {
  const { data } = await client.get<ListTemplatesResponse>('/exercise-templates');
  return data;
}

export async function createTemplate(payload: CreateTemplateRequest): Promise<TemplateResponse> {
  const { data } = await client.post<TemplateResponse>('/exercise-templates', payload);
  return data;
}

export async function listRecommendations(
  patientId: string,
  params?: { limit?: number; offset?: number },
): Promise<ListRecommendationsResponse> {
  const { data } = await client.get<ListRecommendationsResponse>(`/patients/${patientId}/recommendations`, {
    params,
  });
  return data;
}
