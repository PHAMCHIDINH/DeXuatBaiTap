import { RecommendationPlan } from '../exercises/types';

export interface CreatePredictionRequest {
  age_years: number;
  gender: number;
  height: number;
  weight: number;
  ap_hi: number;
  ap_lo: number;
  cholesterol: number;
  gluc: number;
  smoke: number;
  alco: number;
  active: number;
}

export interface PredictionResponse {
  id: string;
  patient_id: string;
  probability: number;
  risk_label: string;
  model_version?: string;
  raw_features?: Record<string, unknown> | null;
  created_at: string;
}

export interface ListPredictionsParams {
  limit?: number;
  offset?: number;
}

export interface CreatePredictionResponse {
  prediction: PredictionResponse;
  recommendation: RecommendationPlan;
}
