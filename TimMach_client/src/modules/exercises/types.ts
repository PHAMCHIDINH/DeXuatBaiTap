export interface TemplateResponse {
  id: number;
  name: string;
  intensity: string;
  description: string;
  duration_min: number;
  freq_per_week: number;
  target_risk_level: string;
  tags: string[];
}

export interface CreateTemplateRequest {
  name: string;
  intensity: string;
  description: string;
  duration_min: number;
  freq_per_week: number;
  target_risk_level: string;
  tags?: string[];
}

export interface RecommendationItem {
  name: string;
  intensity: string;
  duration_min: number;
  freq_per_week: number;
  notes?: string;
}

export interface RecommendationPlan {
  summary: string;
  items: RecommendationItem[];
}

export interface RecommendationResponse {
  id: number;
  patient_id: string;
  prediction_id: string;
  plan?: RecommendationPlan | null;
  created_at: string;
}
