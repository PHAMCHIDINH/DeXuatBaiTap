import client from '../../api/client';
import { CreateTemplateRequest, TemplateResponse, RecommendationResponse } from './types';

export async function listTemplates(): Promise<TemplateResponse[]> {
  const { data } = await client.get<{ templates: TemplateResponse[] }>('/exercise-templates');
  return data.templates;
}

export async function createTemplate(payload: CreateTemplateRequest): Promise<TemplateResponse> {
  const { data } = await client.post<TemplateResponse>('/exercise-templates', payload);
  return data;
}

export async function listRecommendations(
  patientId: string,
  params?: { limit?: number; offset?: number },
): Promise<RecommendationResponse[]> {
  const { data } = await client.get<{ recommendations: RecommendationResponse[] }>(
    `/patients/${patientId}/recommendations`,
    { params },
  );
  return data.recommendations;
}
